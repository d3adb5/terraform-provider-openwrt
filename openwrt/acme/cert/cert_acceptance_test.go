//go:build acceptance.test

package cert_test

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
		"domains":           lucirpc.ListString([]string{"example.com", "www.example.com"}),
		"validation_method": lucirpc.String("dns"),
		"dns":               lucirpc.String("dns_cf"),
		"key_type":          lucirpc.String("ec256"),
		"enabled":           lucirpc.Boolean(true),
	}
	ok, err := client.CreateSection(ctx, "acme", "cert", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_acme_cert" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_acme_cert.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_acme_cert.testing", "domains.0", "example.com"),
			resource.TestCheckResourceAttr("data.openwrt_acme_cert.testing", "domains.1", "www.example.com"),
			resource.TestCheckResourceAttr("data.openwrt_acme_cert.testing", "validation_method", "dns"),
			resource.TestCheckResourceAttr("data.openwrt_acme_cert.testing", "dns", "dns_cf"),
			resource.TestCheckResourceAttr("data.openwrt_acme_cert.testing", "key_type", "ec256"),
			resource.TestCheckResourceAttr("data.openwrt_acme_cert.testing", "enabled", "true"),
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

resource "openwrt_acme_cert" "testing" {
	id                = "testing"
	domains           = ["example.com", "www.example.com"]
	validation_method = "dns"
	dns               = "dns_cf"
	key_type          = "ec256"
	enabled           = true
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "domains.0", "example.com"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "domains.1", "www.example.com"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "validation_method", "dns"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "dns", "dns_cf"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "key_type", "ec256"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "enabled", "true"),
		),
	}

	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_acme_cert.testing",
		// credentials is Sensitive — after import the state has the real value
		// but the empty config has no credentials, so verify would fail.
		ImportStateVerifyIgnore: []string{"credentials"},
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_acme_cert" "testing" {
	id                = "testing"
	domains           = ["example.com", "www.example.com"]
	validation_method = "dns"
	dns               = "dns_cf"
	credentials       = ["CF_Email=user@example.com", "CF_Key=abc123"]
	key_type          = "ec256"
	enabled           = true
	staging           = true
	days              = 30
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "domains.0", "example.com"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "domains.1", "www.example.com"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "validation_method", "dns"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "dns", "dns_cf"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "key_type", "ec256"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "enabled", "true"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "staging", "true"),
			resource.TestCheckResourceAttr("openwrt_acme_cert.testing", "days", "30"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
