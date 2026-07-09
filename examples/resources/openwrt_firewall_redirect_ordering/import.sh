# The import id is always "firewall.redirect".
# Importing adopts every redirect section in its current order;
# the next plan then compares that order against the configuration.

terraform import openwrt_firewall_redirect_ordering.this firewall.redirect
