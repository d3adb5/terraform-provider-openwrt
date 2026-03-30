//go:build acceptance.test

package networkinterface_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ORFops/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/ORFops/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestResourceInvalidIPValidationAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	invalidIPAddr := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_interface" "testing" {
	device  = "eth0"
	id      = "testing"
	ipaddr  = "192x168x1x1"
	netmask = "255.255.255.0"
	proto   = "static"
}
`,
			providerBlock,
		),
		ExpectError: regexp.MustCompile(`must be a valid IP address`),
	}
	invalidNetmask := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_interface" "testing" {
	device  = "eth0"
	id      = "testing"
	ipaddr  = "192.168.1.1"
	netmask = "255x255x255x0"
	proto   = "static"
}
`,
			providerBlock,
		),
		ExpectError: regexp.MustCompile(`must be a valid netmask`),
	}
	invalidGateway := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_interface" "testing" {
	device  = "eth0"
	gateway = "192x168x1x1"
	id      = "testing"
	ipaddr  = "192.168.1.1"
	netmask = "255.255.255.0"
	proto   = "static"
}
`,
			providerBlock,
		),
		ExpectError: regexp.MustCompile(`must be a valid gateway`),
	}

	acceptancetest.TerraformSteps(
		t,
		invalidIPAddr,
		invalidNetmask,
		invalidGateway,
	)
}

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
		"device":  lucirpc.String("br-testing"),
		"ipaddr":  lucirpc.String("192.168.3.1"),
		"netmask": lucirpc.String("255.255.255.0"),
		"proto":   lucirpc.String("static"),
	}
	ok, err := client.CreateSection(ctx, "network", "interface", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_network_interface" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "ipaddr", "192.168.3.1"),
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "netmask", "255.255.255.0"),
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "proto", "static"),
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

resource "openwrt_network_interface" "testing" {
	device = "br-testing"
	id = "testing"
	ipaddr = "192.168.3.1"
	netmask = "255.255.255.0"
	proto = "static"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "ipaddr", "192.168.3.1"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "netmask", "255.255.255.0"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "proto", "static"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_interface.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_interface" "testing" {
	device = "br-testing"
	dns = [
		"9.9.9.9",
		"1.1.1.1",
	]
	id = "testing"
	ipaddr = "192.168.3.1"
	macaddr = "12:34:56:78:90:ab"
	mtu = 1505
	netmask = "255.255.255.0"
	proto = "static"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "dns.1", "1.1.1.1"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "dns.0", "9.9.9.9"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "ipaddr", "192.168.3.1"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "macaddr", "12:34:56:78:90:ab"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "mtu", "1505"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "netmask", "255.255.255.0"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "proto", "static"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}

func TestResourcePeerDNSWithDHCPAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	step := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_interface" "testing" {
	device = "br-testing"
	dns = [
		"9.9.9.9",
		"1.1.1.1",
	]
	id = "testing"
	peerdns = false
	proto = "dhcp"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "dns.0", "9.9.9.9"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "dns.1", "1.1.1.1"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "peerdns", "false"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "proto", "dhcp"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		step,
	)
}

func TestResourcePeerDNSWithDHCPV6Acceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	step := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_interface" "testing" {
	device = "br-testing"
	dns = [
		"9.9.9.9",
		"1.1.1.1",
	]
	id = "testing"
	peerdns = false
	proto = "dhcpv6"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "dns.0", "9.9.9.9"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "dns.1", "1.1.1.1"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "peerdns", "false"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "proto", "dhcpv6"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		step,
	)
}

func TestResourcePppoeAcceptance(t *testing.T) {
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

resource "openwrt_network_interface" "testing" {
	device    = "eth0"
	id        = "testing"
	keepalive = "5 15"
	proto     = "pppoe"
	username  = "user@isp.example"
	password  = "s3cr3t"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "device", "eth0"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "keepalive", "5 15"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "proto", "pppoe"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "username", "user@isp.example"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "password", "s3cr3t"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_interface.testing",
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
	)
}

func TestResourcePppoeUsernameWithoutPppoeProtoAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	invalidStep := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_interface" "testing" {
	device   = "eth0"
	id       = "testing"
	proto    = "dhcp"
	username = "user@isp.example"
}
`,
			providerBlock,
		),
		ExpectError: regexp.MustCompile(`username`),
	}

	acceptancetest.TerraformSteps(
		t,
		invalidStep,
	)
}
