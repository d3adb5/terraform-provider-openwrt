//go:build acceptance.test

package settings_test

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
		"account_email": lucirpc.String("admin@example.com"),
	}
	ok, err := client.CreateSection(ctx, "acme", "acme", "acme", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_acme_acme" "testing" {
	id = "acme"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_acme_acme.testing", "id", "acme"),
			resource.TestCheckResourceAttr("data.openwrt_acme_acme.testing", "account_email", "admin@example.com"),
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

resource "openwrt_acme_acme" "testing" {
	id            = "acme"
	account_email = "admin@example.com"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_acme_acme.testing", "id", "acme"),
			resource.TestCheckResourceAttr("openwrt_acme_acme.testing", "account_email", "admin@example.com"),
		),
	}

	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_acme_acme.testing",
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_acme_acme" "testing" {
	id            = "acme"
	account_email = "updated@example.com"
	debug         = true
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_acme_acme.testing", "id", "acme"),
			resource.TestCheckResourceAttr("openwrt_acme_acme.testing", "account_email", "updated@example.com"),
			resource.TestCheckResourceAttr("openwrt_acme_acme.testing", "debug", "true"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
