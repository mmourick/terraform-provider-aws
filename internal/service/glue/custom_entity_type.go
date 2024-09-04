// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package glue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	awstypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	fwflex "github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource("aws_glue_custom_entity_type",name="Custom Entity Type")
func newResourcCustomEntityType(context.Context) (resource.ResourceWithConfigure, error) {
	r := &resourceCustomEntityType{}

	return r, nil
}

const (
	ResNameCustomEntityType = "aws_glue_custom_entity_type"
)

type resourceCustomEntityType struct {
	framework.ResourceWithConfigure
}

func (r *resourceCustomEntityType) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "aws_glue_custom_entity_type"
}

func (r *resourceCustomEntityType) Schema(ctx context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"context_words": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"regex_string": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			names.AttrTags:    tftags.TagsAttribute(),
			names.AttrTagsAll: tftags.TagsAttributeComputedOnly(),
		},
	}

	response.Schema = s
}

func (r *resourceCustomEntityType) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	conn := r.Meta().GlueClient(ctx)
	var plan resourceCustomEntityTypeData

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)

	if response.Diagnostics.HasError() {
		return
	}

	input := glue.CreateCustomEntityTypeInput{}
	response.Diagnostics.Append(fwflex.Expand(ctx, plan, &input, fwflex.WithFieldNamePrefix("CustomEntityType"))...)

	if response.Diagnostics.HasError() {
		return
	}

	err := retry.RetryContext(ctx, propagationTimeout, func() *retry.RetryError {
		_, err := conn.CreateCustomEntityType(ctx, &input)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.CreateCustomEntityType(ctx, &input)
	}

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Glue, create.ErrActionCreating, ResNameCatalogTableOptimizer, *input.Name, err),
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *resourceCustomEntityType) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	conn := r.Meta().GlueClient(ctx)
	var data resourceCustomEntityTypeData

	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	output, err := FindCustomEntityType(ctx, conn, data.Name.ValueString())

	if tfresource.NotFound(err) {
		response.Diagnostics.Append(fwdiag.NewResourceNotFoundWarningDiagnostic(err))
		response.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Glue, create.ErrActionReading, ResNameCustomEntityType, data.Name.String(), err),
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(fwflex.Flatten(ctx, output, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *resourceCustomEntityType) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (r *resourceCustomEntityType) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	conn := r.Meta().GlueClient(ctx)
	var data resourceCustomEntityTypeData

	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "deleting Glue Custom Entity Type", map[string]interface{}{
		"context_words": data.ContextWords.String(),
		"name":          data.Name.ValueString(),
		"regex_string":  data.RegexString.ValueString(),
		"tags":          data.Tags.String(),
	})

	_, err := conn.DeleteCustomEntityType(ctx, &glue.DeleteCustomEntityTypeInput{
		Name: data.Name.ValueStringPointer(),
	})

	if errs.IsA[*awstypes.EntityNotFoundException](err) {
		return
	}

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Glue, create.ErrActionDeleting, ResNameCustomEntityType, data.Name.ValueString(), err),
			err.Error(),
		)
		return
	}
}

func (r *resourceCustomEntityType) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	parts, err := flex.ExpandResourceId(request.ID, idParts, false)

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Glue, create.ErrActionImporting, ResNameCustomEntityType, request.ID, err),
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root(names.AttrCatalogID), parts[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root(names.AttrDatabaseName), parts[1])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root(names.AttrTableName), parts[2])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root(names.AttrType), parts[3])...)
}

type resourceCustomEntityTypeData struct {
	ContextWords fwtypes.ListValueOf[types.String] `tfsdk:"context_words"`
	Name         types.String                      `tfsdk:"name"`
	RegexString  types.String                      `tfsdk:"regex_string"`
	Tags         types.Map                         `tfsdk:"tags"`
}

type customEntityTypeConfigurationData struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	RoleARN fwtypes.ARN `tfsdk:"role_arn"`
}

func FindCustomEntityType(ctx context.Context, conn *glue.Client, name string) (*glue.GetCustomEntityTypeOutput, error) {
	input := &glue.GetCustomEntityTypeInput{
		Name: aws.String(name),
	}

	output, err := conn.GetCustomEntityType(ctx, input)

	if errs.IsA[*awstypes.EntityNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
