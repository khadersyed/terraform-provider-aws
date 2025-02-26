package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccVPCNetworkACLsDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_network_acls.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLsDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanValue(dataSourceName, "ids.#", "1"),
				),
			},
		},
	})
}

func TestAccVPCNetworkACLsDataSource_filter(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_network_acls.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLsDataSourceConfig_Filter(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "1"),
				),
			},
		},
	})
}

func TestAccVPCNetworkACLsDataSource_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_network_acls.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLsDataSourceConfig_Tags(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "2"),
				),
			},
		},
	})
}

func TestAccVPCNetworkACLsDataSource_vpcID(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_network_acls.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLsDataSourceConfig_VPCID(rName),
				Check: resource.ComposeTestCheckFunc(
					// The VPC will have a default network ACL
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "3"),
				),
			},
		},
	})
}

func TestAccVPCNetworkACLsDataSource_empty(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_network_acls.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLsDataSourceConfig_Empty(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccNetworkACLsDataSourceConfig_Base(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_acl" "test" {
  count = 2

  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccNetworkACLsDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccNetworkACLsDataSourceConfig_Base(rName), `
data "aws_network_acls" "test" {
  depends_on = [aws_network_acl.test[0], aws_network_acl.test[1]]
}
`)
}

func testAccNetworkACLsDataSourceConfig_Filter(rName string) string {
	return acctest.ConfigCompose(testAccNetworkACLsDataSourceConfig_Base(rName), `
data "aws_network_acls" "test" {
  filter {
    name   = "network-acl-id"
    values = [aws_network_acl.test[0].id]
  }

  depends_on = [aws_network_acl.test[0], aws_network_acl.test[1]]
}
`)
}

func testAccNetworkACLsDataSourceConfig_Tags(rName string) string {
	return acctest.ConfigCompose(testAccNetworkACLsDataSourceConfig_Base(rName), `
data "aws_network_acls" "test" {
  tags = {
    Name = aws_network_acl.test[0].tags.Name
  }

  depends_on = [aws_network_acl.test[0], aws_network_acl.test[1]]
}
`)
}

func testAccNetworkACLsDataSourceConfig_VPCID(rName string) string {
	return acctest.ConfigCompose(testAccNetworkACLsDataSourceConfig_Base(rName), `
data "aws_network_acls" "test" {
  vpc_id = aws_network_acl.test[0].vpc_id

  depends_on = [aws_network_acl.test[0], aws_network_acl.test[1]]
}
`)
}

func testAccNetworkACLsDataSourceConfig_Empty(rName string) string {
	return fmt.Sprintf(`
data "aws_network_acls" "test" {
  tags = {
    Name = %[1]q
  }
}
`, rName)
}
