//go:build acceptance.test

package redirect_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ORFops/terraform-provider-openwrt/internal/acceptancetest"
)

func TestOrderingResourceAcceptance(t *testing.T) {
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

	redirects := `
resource "openwrt_firewall_redirect" "first" {
	id        = "first"
	name      = "first"
	src       = "wan"
	src_ip    = "10.0.0.0/8"
	src_dport = "8080"
	dest_ip   = "192.168.1.100"
	proto     = ["tcp"]
	target    = "DNAT"
}

resource "openwrt_firewall_redirect" "second" {
	id        = "second"
	name      = "second"
	src       = "wan"
	src_dport = "8080"
	dest_ip   = "192.168.1.200"
	proto     = ["tcp"]
	target    = "DNAT"
}
`

	checkDeviceOrder := func(want ...string) resource.TestCheckFunc {
		return func(*terraform.State) error {
			sections, err := client.GetSections(ctx, "firewall", "redirect")
			if err != nil {
				return err
			}

			got := []string{}
			for _, section := range sections {
				name, err := section.GetString(".name")
				if err != nil {
					return err
				}

				got = append(got, name)
			}

			if strings.Join(got, ",") != strings.Join(want, ",") {
				return fmt.Errorf("expected device order %v, got %v", want, got)
			}

			return nil
		}
	}

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

%s

resource "openwrt_firewall_redirect_ordering" "this" {
	ids = [
		openwrt_firewall_redirect.second.id,
		openwrt_firewall_redirect.first.id,
	]
}
`,
			providerBlock,
			redirects,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_redirect_ordering.this", "ids.0", "second"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect_ordering.this", "ids.1", "first"),
			checkDeviceOrder("second", "first"),
		),
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

%s

resource "openwrt_firewall_redirect_ordering" "this" {
	ids = [
		openwrt_firewall_redirect.first.id,
		openwrt_firewall_redirect.second.id,
	]
}
`,
			providerBlock,
			redirects,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_redirect_ordering.this", "ids.0", "first"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect_ordering.this", "ids.1", "second"),
			checkDeviceOrder("first", "second"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		updateAndReadResource,
	)
}
