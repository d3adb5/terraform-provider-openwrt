//go:build acceptance.test

package timeserver_test

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
		"enabled":       lucirpc.Boolean(true),
		"enable_server": lucirpc.Boolean(false),
		"server":        lucirpc.ListString([]string{"0.openwrt.pool.ntp.org", "1.openwrt.pool.ntp.org"}),
	}
	ok, err := client.CreateSection(ctx, "system", "timeserver", "ntp", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_system_timeserver" "testing" {
	id = "ntp"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_system_timeserver.testing", "id", "ntp"),
			resource.TestCheckResourceAttr("data.openwrt_system_timeserver.testing", "enabled", "true"),
			resource.TestCheckResourceAttr("data.openwrt_system_timeserver.testing", "enable_server", "false"),
			resource.TestCheckResourceAttr("data.openwrt_system_timeserver.testing", "server.0", "0.openwrt.pool.ntp.org"),
			resource.TestCheckResourceAttr("data.openwrt_system_timeserver.testing", "server.1", "1.openwrt.pool.ntp.org"),
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

resource "openwrt_system_timeserver" "testing" {
	id            = "ntp"
	enabled       = true
	enable_server = false
	server        = ["0.openwrt.pool.ntp.org", "1.openwrt.pool.ntp.org"]
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "id", "ntp"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "enabled", "true"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "enable_server", "false"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "server.0", "0.openwrt.pool.ntp.org"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "server.1", "1.openwrt.pool.ntp.org"),
		),
	}

	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_system_timeserver.testing",
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_system_timeserver" "testing" {
	id            = "ntp"
	enabled       = true
	enable_server = true
	server        = ["0.openwrt.pool.ntp.org", "1.openwrt.pool.ntp.org", "2.openwrt.pool.ntp.org"]
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "id", "ntp"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "enabled", "true"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "enable_server", "true"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "server.0", "0.openwrt.pool.ntp.org"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "server.1", "1.openwrt.pool.ntp.org"),
			resource.TestCheckResourceAttr("openwrt_system_timeserver.testing", "server.2", "2.openwrt.pool.ntp.org"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
