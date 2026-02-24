package networkinterface

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	bringUpOnBootAttribute            = "auto"
	bringUpOnBootAttributeDescription = "Specifies whether to bring up this interface on boot."
	bringUpOnBootUCIOption            = "auto"

	deviceAttribute            = "device"
	deviceAttributeDescription = "Name of the (physical or virtual) device. This name is what the device is known as in LuCI or the `name` field in Terraform. This is not the UCI config name."
	deviceUCIOption            = "device"

	disabledAttribute            = "disabled"
	disabledAttributeDescription = "Disables this interface."
	disabledUCIOption            = "disabled"

	dnsAttribute            = "dns"
	dnsAttributeDescription = "DNS servers"
	dnsUCIOption            = "dns"

	gatewayAttribute            = "gateway"
	gatewayAttributeDescription = "Gateway of the interface"
	gatewayUCIOption            = "gateway"

	ip6AddressAttribute            = "ip6addr"
	ip6AddressAttributeDescription = `Static IPv6 address of the interface in CIDR notation (e.g. "2001:db8::1/64").`
	ip6AddressUCIOption            = "ip6addr"

	ip6AssignAttribute            = "ip6assign"
	ip6AssignAttributeDescription = "Delegate a prefix of given length to this interface"
	ip6AssignUCIOption            = "ip6assign"

	ip6ClassAttribute            = "ip6class"
	ip6ClassAttributeDescription = `Accept only the given class of IPv6 prefixes from upstream (e.g. "wan6").`
	ip6ClassUCIOption            = "ip6class"

	ip6GatewayAttribute            = "ip6gw"
	ip6GatewayAttributeDescription = `IPv6 default gateway of the interface (e.g. "2001:db8::1").`
	ip6GatewayUCIOption            = "ip6gw"

	ip6HintAttribute            = "ip6hint"
	ip6HintAttributeDescription = "Subprefix ID hint for prefix delegation, in hexadecimal (e.g. \"1\")."
	ip6HintUCIOption            = "ip6hint"

	ip6IfaceIDAttribute            = "ip6ifaceid"
	ip6IfaceIDAttributeDescription = `IPv6 interface identifier suffix. Accepted values: "eui64", "random", or a fixed suffix like "::1".`
	ip6IfaceIDUCIOption            = "ip6ifaceid"

	ip6PrefixAttribute            = "ip6prefix"
	ip6PrefixAttributeDescription = `IPv6 prefix for distribution to downstream interfaces, in CIDR notation (e.g. "2001:db8:1::/48").`
	ip6PrefixUCIOption            = "ip6prefix"

	ipAddressAttribute            = "ipaddr"
	ipAddressAttributeDescription = "IP address of the interface"
	ipAddressUCIOption            = "ipaddr"

	ipv6Attribute            = "ipv6"
	ipv6AttributeDescription = "Enable or disable IPv6 on this interface."
	ipv6UCIOption            = "ipv6"

	macAddressAttribute            = "macaddr"
	macAddressAttributeDescription = "Override the MAC Address of this interface."
	macAddressUCIOption            = "macaddr"

	metricAttribute            = "metric"
	metricAttributeDescription = "Route metric for this interface. Lower values are preferred."
	metricUCIOption            = "metric"

	mtuAttribute            = "mtu"
	mtuAttributeDescription = "Override the default MTU on this interface."
	mtuUCIOption            = "mtu"

	netmaskAttribute            = "netmask"
	netmaskAttributeDescription = "Netmask of the interface"
	netmaskUCIOption            = "netmask"

	peerDNSAttribute            = "peerdns"
	peerDNSAttributeDescription = "Use DHCP-provided DNS servers."
	peerDNSUCIOption            = "peerdns"

	pppoeACAttribute            = "ac"
	pppoeACAttributeDescription = "PPPoE Access Concentrator name to connect to (PPPoE only)."
	pppoeACUCIOption            = "ac"

	pppoeKeepAliveAttribute            = "keepalive"
	pppoeKeepAliveAttributeDescription = `LCP echo keepalive configuration in the format "failures interval" (e.g. "5 15" means 5 failures with 15 second interval). PPPoE only.`
	pppoeKeepAliveUCIOption            = "keepalive"

	pppoePasswordAttribute            = "password"
	pppoePasswordAttributeDescription = "PPPoE password for authentication. PPPoE only."
	pppoePasswordUCIOption            = "password"

	pppoeServiceAttribute            = "service"
	pppoeServiceAttributeDescription = "PPPoE service name (ISP-specific). PPPoE only."
	pppoeServiceUCIOption            = "service"

	pppoeUsernameAttribute            = "username"
	pppoeUsernameAttributeDescription = "PPPoE username for authentication. PPPoE only."
	pppoeUsernameUCIOption            = "username"

	protocolAttribute            = "proto"
	protocolAttributeDescription = `The protocol type of the interface (e.g. "static", "dhcp", "dhcpv6", "pppoe").`
	protocolDHCP                 = "dhcp"
	protocolDHCPV6               = "dhcpv6"
	protocolPPPoE                = "pppoe"
	protocolStatic               = "static"
	protocolUCIOption            = "proto"

	requestingAddressAttribute            = "reqaddress"
	requestingAddressAttributeDescription = `Behavior for requesting address. Can only be one of "force", "try", or "none".`
	requestingAddressForce                = "force"
	requestingAddressNone                 = "none"
	requestingAddressTry                  = "try"
	requestingAddressUCIOption            = "reqaddress"

	// The fact we can only support `"auto"` is because we haven't figured out how to represent unions.
	// Once we do,
	// we can support `"auto"`, `no`, or 0-64.
	requestingPrefixAttribute            = "reqprefix"
	requestingPrefixAttributeDescription = `Behavior for requesting prefixes. Currently, only "auto" is supported.`
	requestingPrefixAuto                 = "auto"
	requestingPrefixUCIOption            = "reqprefix"

	schemaDescription = "A logic network."

	uciConfig = "network"
	uciType   = "interface"
)

var (
	bringUpOnBootSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       bringUpOnBootAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetBringUpOnBoot, bringUpOnBootAttribute, bringUpOnBootUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetBringUpOnBoot, bringUpOnBootAttribute, bringUpOnBootUCIOption),
	}

	deviceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       deviceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDevice, deviceAttribute, deviceUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDevice, deviceAttribute, deviceUCIOption),
	}

	disabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       disabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDisabled, disabledAttribute, disabledUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDisabled, disabledAttribute, disabledUCIOption),
	}

	dnsSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dnsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetDNS, dnsAttribute, dnsUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetDNS, dnsAttribute, dnsUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.AnyWithAllWarnings(
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCP,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCPV6,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolPPPoE,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolStatic,
				),
			),
		},
	}

	gatewaySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       gatewayAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetGateway, gatewayAttribute, gatewayUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetGateway, gatewayAttribute, gatewayUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^([[:digit:]]{1,3}\.){3}[[:digit:]]{1,3}$`),
				`must be a valid gateway (e.g. "192.168.1.1")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	ip6AddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6AddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIP6Address, ip6AddressAttribute, ip6AddressUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIP6Address, ip6AddressAttribute, ip6AddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^[0-9a-fA-F:]+/\d{1,3}$`),
				`must be a valid IPv6 CIDR address (e.g. "2001:db8::1/64")`,
			),
		},
	}

	ip6AssignSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6AssignAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetIP6Assign, ip6AssignAttribute, ip6AssignUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetIP6Assign, ip6AssignAttribute, ip6AssignUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(0, 64),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	ip6ClassSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6ClassAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIP6Class, ip6ClassAttribute, ip6ClassUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIP6Class, ip6ClassAttribute, ip6ClassUCIOption),
	}

	ip6GatewaySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6GatewayAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIP6Gateway, ip6GatewayAttribute, ip6GatewayUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIP6Gateway, ip6GatewayAttribute, ip6GatewayUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^[0-9a-fA-F:]+$`),
				`must be a valid IPv6 address (e.g. "2001:db8::1")`,
			),
		},
	}

	ip6HintSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6HintAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIP6Hint, ip6HintAttribute, ip6HintUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIP6Hint, ip6HintAttribute, ip6HintUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^[0-9a-fA-F]+$`),
				`must be a hexadecimal value (e.g. "1" or "a3")`,
			),
		},
	}

	ip6IfaceIDSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6IfaceIDAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIP6IfaceID, ip6IfaceIDAttribute, ip6IfaceIDUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIP6IfaceID, ip6IfaceIDAttribute, ip6IfaceIDUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^(eui64|random|::[0-9a-fA-F:]+)$`),
				`must be "eui64", "random", or a fixed IPv6 suffix (e.g. "::1")`,
			),
		},
	}

	ip6PrefixSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6PrefixAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIP6Prefix, ip6PrefixAttribute, ip6PrefixUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIP6Prefix, ip6PrefixAttribute, ip6PrefixUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^[0-9a-fA-F:]+/\d{1,3}$`),
				`must be a valid IPv6 CIDR prefix (e.g. "2001:db8:1::/48")`,
			),
		},
	}

	ipAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^([[:digit:]]{1,3}\.){3}[[:digit:]]{1,3}$`),
				`must be a valid IP address (e.g. "192.168.3.1")`,
			),
		},
	}

	ipv6SchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipv6AttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetIPv6, ipv6Attribute, ipv6UCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetIPv6, ipv6Attribute, ipv6UCIOption),
	}

	macAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       macAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetMacAddress, macAddressAttribute, macAddressUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetMacAddress, macAddressAttribute, macAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:xdigit:]][[:xdigit:]]:){5}[[:xdigit:]][[:xdigit:]]$"),
				`must be a valid MAC address (e.g. "12:34:56:78:90:ab")`,
			),
		},
	}

	metricSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       metricAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMetric, metricAttribute, metricUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMetric, metricAttribute, metricUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(0, 4294967295),
		},
	}

	mtuSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mtuAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMTU, mtuAttribute, mtuUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMTU, mtuAttribute, mtuUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(576, 9200),
		},
	}

	netmaskSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       netmaskAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetNetmask, netmaskAttribute, netmaskUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetNetmask, netmaskAttribute, netmaskUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^([[:digit:]]{1,3}\.){3}[[:digit:]]{1,3}$`),
				`must be a valid netmask (e.g. "255.255.255.0")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	peerDNSSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       peerDNSAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetPeerDNS, peerDNSAttribute, peerDNSUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetPeerDNS, peerDNSAttribute, peerDNSUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.AnyBool(
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCP,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCPV6,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolPPPoE,
				),
			),
		},
	}

	pppoeACSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       pppoeACAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetPPPoEAC, pppoeACAttribute, pppoeACUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetPPPoEAC, pppoeACAttribute, pppoeACUCIOption),
		Validators: []validator.String{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolPPPoE,
			),
		},
	}

	pppoeKeepAliveSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       pppoeKeepAliveAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetPPPoEKeepAlive, pppoeKeepAliveAttribute, pppoeKeepAliveUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetPPPoEKeepAlive, pppoeKeepAliveAttribute, pppoeKeepAliveUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^\d+ \d+$`),
				`must be in the format "failures interval" (e.g. "5 15")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolPPPoE,
			),
		},
	}

	pppoePasswordSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       pppoePasswordAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetPPPoEPassword, pppoePasswordAttribute, pppoePasswordUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		Sensitive:         true,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetPPPoEPassword, pppoePasswordAttribute, pppoePasswordUCIOption),
		Validators: []validator.String{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolPPPoE,
			),
		},
	}

	pppoeServiceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       pppoeServiceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetPPPoEService, pppoeServiceAttribute, pppoeServiceUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetPPPoEService, pppoeServiceAttribute, pppoeServiceUCIOption),
		Validators: []validator.String{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolPPPoE,
			),
		},
	}

	pppoeUsernameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       pppoeUsernameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetPPPoEUsername, pppoeUsernameAttribute, pppoeUsernameUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		Sensitive:         true,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetPPPoEUsername, pppoeUsernameAttribute, pppoeUsernameUCIOption),
		Validators: []validator.String{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolPPPoE,
			),
		},
	}

	protocolSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       protocolAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetProtocol, protocolAttribute, protocolUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetProtocol, protocolAttribute, protocolUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				protocolDHCP,
				protocolDHCPV6,
				protocolPPPoE,
				protocolStatic,
			),
		},
	}

	requestingAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       requestingAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetRequestingAddress, requestingAddressAttribute, requestingAddressUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetRequestingAddress, requestingAddressAttribute, requestingAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				requestingAddressForce,
				requestingAddressNone,
				requestingAddressTry,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolDHCPV6,
			),
		},
	}

	requestingPrefixSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       requestingPrefixAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetRequestingPrefix, requestingPrefixAttribute, requestingPrefixUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetRequestingPrefix, requestingPrefixAttribute, requestingPrefixUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				requestingPrefixAuto,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolDHCPV6,
			),
		},
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		bringUpOnBootAttribute:     bringUpOnBootSchemaAttribute,
		deviceAttribute:            deviceSchemaAttribute,
		disabledAttribute:          disabledSchemaAttribute,
		dnsAttribute:               dnsSchemaAttribute,
		gatewayAttribute:           gatewaySchemaAttribute,
		ip6AddressAttribute:        ip6AddressSchemaAttribute,
		ip6AssignAttribute:         ip6AssignSchemaAttribute,
		ip6ClassAttribute:          ip6ClassSchemaAttribute,
		ip6GatewayAttribute:        ip6GatewaySchemaAttribute,
		ip6HintAttribute:           ip6HintSchemaAttribute,
		ip6IfaceIDAttribute:        ip6IfaceIDSchemaAttribute,
		ip6PrefixAttribute:         ip6PrefixSchemaAttribute,
		ipAddressAttribute:         ipAddressSchemaAttribute,
		ipv6Attribute:              ipv6SchemaAttribute,
		macAddressAttribute:        macAddressSchemaAttribute,
		mtuAttribute:               mtuSchemaAttribute,
		metricAttribute:            metricSchemaAttribute,
		netmaskAttribute:           netmaskSchemaAttribute,
		peerDNSAttribute:            peerDNSSchemaAttribute,
		pppoeACAttribute:            pppoeACSchemaAttribute,
		pppoeKeepAliveAttribute:     pppoeKeepAliveSchemaAttribute,
		pppoePasswordAttribute:      pppoePasswordSchemaAttribute,
		pppoeServiceAttribute:       pppoeServiceSchemaAttribute,
		pppoeUsernameAttribute:      pppoeUsernameSchemaAttribute,
		protocolAttribute:           protocolSchemaAttribute,
		requestingAddressAttribute: requestingAddressSchemaAttribute,
		requestingPrefixAttribute:  requestingPrefixSchemaAttribute,
		lucirpcglue.IdAttribute:    lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
	}
)

func NewDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

func NewResource() resource.Resource {
	return lucirpcglue.NewResource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

type model struct {
	BringUpOnBoot     types.Bool   `tfsdk:"auto"`
	Device            types.String `tfsdk:"device"`
	Disabled          types.Bool   `tfsdk:"disabled"`
	DNS               types.List   `tfsdk:"dns"`
	Gateway           types.String `tfsdk:"gateway"`
	Id                types.String `tfsdk:"id"`
	IP6Address        types.String `tfsdk:"ip6addr"`
	IP6Assign         types.Int64  `tfsdk:"ip6assign"`
	IP6Class          types.String `tfsdk:"ip6class"`
	IP6Gateway        types.String `tfsdk:"ip6gw"`
	IP6Hint           types.String `tfsdk:"ip6hint"`
	IP6IfaceID        types.String `tfsdk:"ip6ifaceid"`
	IP6Prefix         types.String `tfsdk:"ip6prefix"`
	IPAddress         types.String `tfsdk:"ipaddr"`
	IPv6              types.Bool   `tfsdk:"ipv6"`
	MacAddress        types.String `tfsdk:"macaddr"`
	Metric            types.Int64  `tfsdk:"metric"`
	MTU               types.Int64  `tfsdk:"mtu"`
	Netmask           types.String `tfsdk:"netmask"`
	PeerDNS           types.Bool   `tfsdk:"peerdns"`
	PPPoEAC           types.String `tfsdk:"ac"`
	PPPoEKeepAlive    types.String `tfsdk:"keepalive"`
	PPPoEPassword     types.String `tfsdk:"password"`
	PPPoEService      types.String `tfsdk:"service"`
	PPPoEUsername     types.String `tfsdk:"username"`
	Protocol          types.String `tfsdk:"proto"`
	RequestingAddress types.String `tfsdk:"reqaddress"`
	RequestingPrefix  types.String `tfsdk:"reqprefix"`
}

func modelGetBringUpOnBoot(m model) types.Bool       { return m.BringUpOnBoot }
func modelGetDevice(m model) types.String            { return m.Device }
func modelGetDisabled(m model) types.Bool            { return m.Disabled }
func modelGetDNS(m model) types.List                 { return m.DNS }
func modelGetGateway(m model) types.String           { return m.Gateway }
func modelGetId(m model) types.String                { return m.Id }
func modelGetIP6Address(m model) types.String        { return m.IP6Address }
func modelGetIP6Assign(m model) types.Int64          { return m.IP6Assign }
func modelGetIP6Class(m model) types.String          { return m.IP6Class }
func modelGetIP6Gateway(m model) types.String        { return m.IP6Gateway }
func modelGetIP6Hint(m model) types.String           { return m.IP6Hint }
func modelGetIP6IfaceID(m model) types.String        { return m.IP6IfaceID }
func modelGetIP6Prefix(m model) types.String         { return m.IP6Prefix }
func modelGetIPAddress(m model) types.String         { return m.IPAddress }
func modelGetIPv6(m model) types.Bool                { return m.IPv6 }
func modelGetMacAddress(m model) types.String        { return m.MacAddress }
func modelGetMetric(m model) types.Int64             { return m.Metric }
func modelGetMTU(m model) types.Int64                { return m.MTU }
func modelGetNetmask(m model) types.String           { return m.Netmask }
func modelGetPeerDNS(m model) types.Bool             { return m.PeerDNS }
func modelGetPPPoEAC(m model) types.String           { return m.PPPoEAC }
func modelGetPPPoEKeepAlive(m model) types.String    { return m.PPPoEKeepAlive }
func modelGetPPPoEPassword(m model) types.String     { return m.PPPoEPassword }
func modelGetPPPoEService(m model) types.String      { return m.PPPoEService }
func modelGetPPPoEUsername(m model) types.String     { return m.PPPoEUsername }
func modelGetProtocol(m model) types.String          { return m.Protocol }
func modelGetRequestingAddress(m model) types.String { return m.RequestingAddress }
func modelGetRequestingPrefix(m model) types.String  { return m.RequestingPrefix }

func modelSetBringUpOnBoot(m *model, value types.Bool)       { m.BringUpOnBoot = value }
func modelSetDevice(m *model, value types.String)            { m.Device = value }
func modelSetDisabled(m *model, value types.Bool)            { m.Disabled = value }
func modelSetDNS(m *model, value types.List)                 { m.DNS = value }
func modelSetGateway(m *model, value types.String)           { m.Gateway = value }
func modelSetId(m *model, value types.String)                { m.Id = value }
func modelSetIP6Address(m *model, value types.String)        { m.IP6Address = value }
func modelSetIP6Assign(m *model, value types.Int64)          { m.IP6Assign = value }
func modelSetIP6Class(m *model, value types.String)          { m.IP6Class = value }
func modelSetIP6Gateway(m *model, value types.String)        { m.IP6Gateway = value }
func modelSetIP6Hint(m *model, value types.String)           { m.IP6Hint = value }
func modelSetIP6IfaceID(m *model, value types.String)        { m.IP6IfaceID = value }
func modelSetIP6Prefix(m *model, value types.String)         { m.IP6Prefix = value }
func modelSetIPAddress(m *model, value types.String)         { m.IPAddress = value }
func modelSetIPv6(m *model, value types.Bool)                { m.IPv6 = value }
func modelSetMacAddress(m *model, value types.String)        { m.MacAddress = value }
func modelSetMetric(m *model, value types.Int64)             { m.Metric = value }
func modelSetMTU(m *model, value types.Int64)                { m.MTU = value }
func modelSetNetmask(m *model, value types.String)           { m.Netmask = value }
func modelSetPeerDNS(m *model, value types.Bool)             { m.PeerDNS = value }
func modelSetPPPoEAC(m *model, value types.String)           { m.PPPoEAC = value }
func modelSetPPPoEKeepAlive(m *model, value types.String)    { m.PPPoEKeepAlive = value }
func modelSetPPPoEPassword(m *model, value types.String)     { m.PPPoEPassword = value }
func modelSetPPPoEService(m *model, value types.String)      { m.PPPoEService = value }
func modelSetPPPoEUsername(m *model, value types.String)     { m.PPPoEUsername = value }
func modelSetProtocol(m *model, value types.String)          { m.Protocol = value }
func modelSetRequestingAddress(m *model, value types.String) { m.RequestingAddress = value }
func modelSetRequestingPrefix(m *model, value types.String)  { m.RequestingPrefix = value }
