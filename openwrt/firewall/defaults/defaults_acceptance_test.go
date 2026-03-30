//go:build acceptance.test

package defaults_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ORFops/terraform-provider-openwrt/internal/acceptancetest"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_firewall_defaults" "this" {
	id = "cfg01e63d"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_firewall_defaults.this", "id", "cfg01e63d"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_defaults.this", "input", "REJECT"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_defaults.this", "output", "ACCEPT"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_defaults.this", "forward", "REJECT"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_defaults.this", "synflood_protect", "true"),
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

	importValidation := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_defaults" "this" {
	id = "cfg01e63d"
}
`,
			providerBlock,
		),
		ImportState:        true,
		ImportStateId:      "cfg01e63d",
		ImportStatePersist: true,
		ResourceName:       "openwrt_firewall_defaults.this",
	}

	readResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_defaults" "this" {
	id = "cfg01e63d"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "id", "cfg01e63d"),
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "input", "REJECT"),
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "output", "ACCEPT"),
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "forward", "REJECT"),
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "synflood_protect", "true"),
		),
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_defaults" "this" {
	id           = "cfg01e63d"
	drop_invalid = true
	forward      = "REJECT"
	input        = "REJECT"
	output       = "ACCEPT"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "id", "cfg01e63d"),
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "drop_invalid", "true"),
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "forward", "REJECT"),
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "input", "REJECT"),
			resource.TestCheckResourceAttr("openwrt_firewall_defaults.this", "output", "ACCEPT"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		importValidation,
		readResource,
		updateAndReadResource,
	)
}
