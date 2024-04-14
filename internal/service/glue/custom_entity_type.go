package glue

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_glue_custom_entity_type", name="CustomEntityType")
// @Tags(identifierAttribute="arn")
func ResourceCustomEntityType() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCustomEntityTypeCreate,
		ReadWithoutTimeout:   resourceCustomEntityTypeRead,
		DeleteWithoutTimeout: resourceCustomEntityTypeDelete,
		// DeleteWithoutTimeout: resourceRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"regex_string": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"context_words": {
				Type:     schema.TypeList,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceCustomEntityTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	input := &glue.CreateCustomEntityTypeInput{
		Name:        aws.String(d.Get("name").(string)),
		RegexString: aws.String(d.Get("regex_string").(string)),
		Tags:        getTagsIn(ctx),
	}

	if v, ok := d.GetOk("context_words"); ok {
		input.ContextWords = aws.StringSlice(v.([]string))
	}

	log.Printf("[DEBUG] Creating Glue Custom Entity Type: %s", input)
	output, err := conn.CreateCustomEntityTypeWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Glue Custom Entity Type: %s", err)
	}

	d.SetId(aws.StringValue(output.Name))

	return append(diags, resourceCustomEntityTypeRead(ctx, d, meta)...)
}

func resourceCustomEntityTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	output, err := conn.GetCustomEntityTypeWithContext(ctx, &glue.GetCustomEntityTypeInput{
		Name: aws.String(d.Id()),
	})

	if err != nil {
		if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
			log.Printf("[WARN] Glue Custom Entity Type (%s) not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "reading Glue Custom Entity Type (%s): %s", d.Id(), err)
	}

	if output == nil {
		log.Printf("[WARN] Glue Custom Entity Type (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	d.Set("name", output.Name)
	d.Set("regex_string", output.RegexString)
	d.Set("context_words", output.ContextWords)

	return diags
}

func resourceCustomEntityTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	log.Printf("[DEBUG] Deleting Glue Custom Entity Type: %s", d.Id())

	input := &glue.DeleteCustomEntityTypeInput{
		Name: aws.String(d.Id()),
	}

	_, err := conn.DeleteCustomEntityTypeWithContext(ctx, input)

	if err != nil {
		if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "deleting Glue Custom Entity Type (%s): %s", d.Id(), err)
	}

	return diags
}
