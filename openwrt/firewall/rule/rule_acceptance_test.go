//go:build acceptance.test

package rule_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ORFops/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/ORFops/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	client := openWrtServer.LuCIRPCClient(
		ctx,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()
	options := lucirpc.Options{
		"name":   lucirpc.String("testing"),
		"target": lucirpc.String("ACCEPT"),
	}
	ok, err := client.CreateSection(ctx, "firewall", "rule", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_firewall_rule" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_firewall_rule.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_rule.testing", "name", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_rule.testing", "target", "ACCEPT"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestResourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_rule" "testing" {
	id     = "testing"
	name   = "testing"
	target = "ACCEPT"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "name", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "target", "ACCEPT"),
		),
	}

	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_firewall_rule.testing",
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_rule" "testing" {
	id     = "testing"
	name   = "testing"
	target = "ACCEPT"
	src    = "lan"
	dest   = "wan"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "name", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "target", "ACCEPT"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "src", "lan"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "dest", "wan"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
