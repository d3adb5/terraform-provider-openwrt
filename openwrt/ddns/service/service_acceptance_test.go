//go:build acceptance.test

package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
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
		"service_name": lucirpc.String("cloudflare.com"),
		"domain":       lucirpc.String("example.com"),
		"username":     lucirpc.String("user@example.com"),
		"password":     lucirpc.String("secret"),
		"interface":    lucirpc.String("wan"),
	}
	ok, err := client.CreateSection(ctx, "ddns", "service", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_ddns_service" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_ddns_service.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_ddns_service.testing", "service_name", "cloudflare.com"),
			resource.TestCheckResourceAttr("data.openwrt_ddns_service.testing", "domain", "example.com"),
			resource.TestCheckResourceAttr("data.openwrt_ddns_service.testing", "username", "user@example.com"),
			resource.TestCheckResourceAttr("data.openwrt_ddns_service.testing", "interface", "wan"),
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

resource "openwrt_ddns_service" "testing" {
	id           = "testing"
	service_name = "cloudflare.com"
	domain       = "example.com"
	username     = "user@example.com"
	password     = "secret"
	interface    = "wan"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "service_name", "cloudflare.com"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "domain", "example.com"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "username", "user@example.com"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "interface", "wan"),
		),
	}

	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_ddns_service.testing",
		// password is Sensitive — the provider reads it back from UCI but
		// Terraform's import verify compares plan-time (null/unknown) vs state.
		// Since password is Optional+Sensitive and not Computed, after import
		// the state has the real value but the config block has no password,
		// so import verify would fail. Skip it.
		ImportStateVerifyIgnore: []string{"password"},
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_ddns_service" "testing" {
	id             = "testing"
	service_name   = "cloudflare.com"
	domain         = "example.com"
	username       = "user@example.com"
	password       = "secret"
	interface      = "wan"
	enabled        = true
	check_interval = 10
	check_unit     = "minutes"
	force_interval = 72
	force_unit     = "hours"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "service_name", "cloudflare.com"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "domain", "example.com"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "username", "user@example.com"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "interface", "wan"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "enabled", "true"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "check_interval", "10"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "check_unit", "minutes"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "force_interval", "72"),
			resource.TestCheckResourceAttr("openwrt_ddns_service.testing", "force_unit", "hours"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
