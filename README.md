# terraform-provider-openwrt

A [Terraform][] provider for [OpenWrt][] devices, published at
[registry.terraform.io/ORFops/openwrt](https://registry.terraform.io/providers/ORFops/openwrt/latest).

## Usage

```hcl
terraform {
  required_providers {
    openwrt = {
      source  = "ORFops/openwrt"
      version = "~> 0.1"
    }
  }
}

provider "openwrt" {
  host     = "192.168.1.1"
  port     = 80
  username = "root"
  password = "password"
}
```

## Supported Resources

| Resource | Data Source |
|---|---|
| `openwrt_acme_acme` | `openwrt_acme_acme` |
| `openwrt_acme_cert` | `openwrt_acme_cert` |
| `openwrt_ddns_service` | `openwrt_ddns_service` |
| `openwrt_dhcp_dhcp` | `openwrt_dhcp_dhcp` |
| `openwrt_dhcp_dnsmasq` | `openwrt_dhcp_dnsmasq` |
| `openwrt_dhcp_domain` | `openwrt_dhcp_domain` |
| `openwrt_dhcp_host` | `openwrt_dhcp_host` |
| `openwrt_dhcp_odhcpd` | `openwrt_dhcp_odhcpd` |
| `openwrt_firewall_defaults` | `openwrt_firewall_defaults` |
| `openwrt_firewall_forwarding` | `openwrt_firewall_forwarding` |
| `openwrt_firewall_redirect` | `openwrt_firewall_redirect` |
| `openwrt_firewall_rule` | `openwrt_firewall_rule` |
| `openwrt_firewall_zone` | `openwrt_firewall_zone` |
| `openwrt_network_bridge_vlan` | `openwrt_network_bridge_vlan` |
| `openwrt_network_device` | `openwrt_network_device` |
| `openwrt_network_globals` | `openwrt_network_globals` |
| `openwrt_network_interface` | `openwrt_network_interface` |
| `openwrt_network_route` | `openwrt_network_route` |
| `openwrt_network_route6` | `openwrt_network_route6` |
| `openwrt_network_rule` | `openwrt_network_rule` |
| `openwrt_network_rule6` | `openwrt_network_rule6` |
| `openwrt_network_switch` | `openwrt_network_switch` |
| `openwrt_network_switch_vlan` | `openwrt_network_switch_vlan` |
| `openwrt_system_system` | `openwrt_system_system` |
| `openwrt_system_timeserver` | `openwrt_system_timeserver` |
| `openwrt_wireless_wifi_device` | `openwrt_wireless_wifi_device` |
| `openwrt_wireless_wifi_iface` | `openwrt_wireless_wifi_iface` |

Full documentation for each resource is available on the
[Terraform Registry](https://registry.terraform.io/providers/ORFops/openwrt/latest/docs).

## Development

[`make`][] is used to build and test this repo.

```sh
make build   # compile the provider
make test    # build + docs check + unit + acceptance tests
make docs    # regenerate docs/
```

### Running acceptance tests

Acceptance tests require Docker to spin up a real OpenWrt instance:

```sh
make start-acceptance-test-server
TF_ACC=1 go test -tags=acceptance.test ./...
make clean
```

### Releasing

Releases are created automatically by GitHub Actions when a tag is pushed:

```sh
git tag v0.x.y
git push origin v0.x.y
```

The workflow builds binaries for all supported platforms, signs the checksums
with GPG, and publishes the release to GitHub. HCP Terraform Registry syncs
automatically from there.

## Credits

This provider is a fork of
[northfuse/terraform-provider-openwrt](https://github.com/northfuse/terraform-provider-openwrt),
which is itself a fork of
[joneshf/terraform-provider-openwrt](https://github.com/joneshf/terraform-provider-openwrt)
by Hardy Jones. These upstream projects laid the foundation for the LuCI
JSON-RPC client, the generic Terraform Plugin Framework glue layer, and the
majority of the supported resources.

Additional development on this fork was assisted by
[Claude Code](https://claude.ai/code) (Anthropic), which contributed bug fixes,
new resources (`network_route`, `network_route6`, `network_rule`, `network_rule6`),
and correctness fixes (firewall UCI option names, schema attribute types).

## License

See [LICENSE](LICENSE).

[`make`]: https://www.gnu.org/software/make/
[openwrt]: https://openwrt.org/
[terraform]: https://www.terraform.io/
