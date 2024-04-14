package glue_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfglue "github.com/hashicorp/terraform-provider-aws/internal/service/glue"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccCustomEntityType_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var customEntityType glue.GetCustomEntityTypeOutput

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resourceName := "aws_glue_custom_entity_type.test"

	resource.ParallelTest(t, resource.TestCase{
		// PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckRegistry(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.GlueServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		// CheckDestroy:             testAccCheckRegistryDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEntityTypeConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckCustomEntityTypeExists(ctx, resourceName, &customEntityType),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "regex_string", "^[0-9A-Za-z_$#-]+$"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCustomEntityTypeExists(ctx context.Context, resourceName string, customEntityType *glue.GetCustomEntityTypeOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Glue CustomEntityType ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).GlueConn(ctx)
		output, err := tfglue.FindCustomEntityTypeByName(ctx, conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		if output == nil {
			return fmt.Errorf("Glue CustomEntityType (%s) not found", rs.Primary.ID)
		}

		if aws.StringValue(output.Name) == rs.Primary.ID {
			*customEntityType = *output
			return nil
		}

		return fmt.Errorf("Glue CustomEntityType (%s) not found", rs.Primary.ID)
	}
}

func testAccCustomEntityTypeConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_glue_custom_entity_type" "test" {
  name = %[1]q
  regex_string = "^[0-9A-Za-z_$#-]+$"
}
`, rName)
}
