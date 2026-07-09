resource "openwrt_firewall_redirect" "specific" {
  id        = "specific"
  name      = "Forward HTTPS from the office"
  src       = "wan"
  src_ip    = "203.0.113.0/24"
  src_dport = "443"
  dest_ip   = "192.168.1.10"
  proto     = ["tcp"]
  target    = "DNAT"
}

resource "openwrt_firewall_redirect" "catch_all" {
  id        = "catch_all"
  name      = "Forward HTTPS from anywhere"
  src       = "wan"
  src_dport = "443"
  dest_ip   = "192.168.1.20"
  proto     = ["tcp"]
  target    = "DNAT"
}

# The first matching redirect wins,
# so the specific redirect must come before the catch-all.
resource "openwrt_firewall_redirect_ordering" "this" {
  ids = [
    openwrt_firewall_redirect.specific.id,
    openwrt_firewall_redirect.catch_all.id,
  ]
}
