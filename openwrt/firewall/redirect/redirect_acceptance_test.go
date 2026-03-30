//go:build acceptance.test

package redirect_test

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
		"name":      lucirpc.String("testing"),
		"src":       lucirpc.String("wan"),
		"src_dport": lucirpc.String("8080"),
		"dest":      lucirpc.String("lan"),
		"dest_ip":   lucirpc.String("192.168.1.100"),
		"dest_port": lucirpc.String("80"),
		"proto":     lucirpc.String("tcp"),
		"target":    lucirpc.String("DNAT"),
	}
	ok, err := client.CreateSection(ctx, "firewall", "redirect", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_firewall_redirect" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "name", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "src", "wan"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "src_dport", "8080"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "dest", "lan"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "dest_ip", "192.168.1.100"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "dest_port", "80"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "proto.0", "tcp"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "target", "DNAT"),
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

resource "openwrt_firewall_redirect" "testing" {
	id        = "testing"
	name      = "testing"
	src       = "wan"
	src_dport = "8080"
	dest      = "lan"
	dest_ip   = "192.168.1.100"
	dest_port = "80"
	proto     = ["tcp"]
	target    = "DNAT"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "name", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "src", "wan"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "src_dport", "8080"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest", "lan"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest_ip", "192.168.1.100"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest_port", "80"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "proto.0", "tcp"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "target", "DNAT"),
		),
	}

	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_firewall_redirect.testing",
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_redirect" "testing" {
	id         = "testing"
	name       = "testing"
	src        = "wan"
	src_dport  = "8080"
	dest       = "lan"
	dest_ip    = "192.168.1.100"
	dest_port  = "80"
	proto      = ["tcp"]
	target     = "DNAT"
	reflection = true
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "name", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "src", "wan"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "src_dport", "8080"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest", "lan"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest_ip", "192.168.1.100"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest_port", "80"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "proto.0", "tcp"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "target", "DNAT"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "reflection", "true"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
