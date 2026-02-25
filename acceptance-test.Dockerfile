FROM openwrt/rootfs:x86_64-SNAPSHOT

RUN mkdir -p /var/lock
RUN apk update && apk add \
    # Install curl so we can make a healthcheck
    # wget is installed, but it's hard to use for a health check.
    curl \
    # Install LuCI JSON-RPC packages.
    # See https://github.com/openwrt/luci/wiki/JsonRpcHowTo#basics
    luci-compat \
    luci-lib-ipkg \
    luci-mod-rpc \
    # Install LuCI (and HTTPS support)
    # This is entirely for debugging/diagnosis purposes.
    luci \
    luci-ssl
RUN /etc/init.d/dropbear enable
RUN /etc/init.d/uhttpd enable
# Create empty config files for packages not installed in the test image
# so UCI can create sections in them (the packages create these on real routers)
RUN touch /etc/config/acme /etc/config/ddns

EXPOSE 22 80 443

CMD ["/sbin/init"]

HEALTHCHECK --interval=5s CMD curl \
    --data '{"id": 1, "method": "login", "params": ["root", ""]}' \
    --fail \
    --no-progress-meter \
    http://localhost/cgi-bin/luci/rpc/auth
