#!/bin/sh
# FreeBSD: stop the agent if the package is removed. On an update with "pkg upgrade" pkg runs this script
# of the old package with PKG_UPGRADE=true, the agent keeps running until it is restarted (see postinst.sh).
# The configuration files remain in /etc/openitcockpit-agent.

if [ "${PKG_UPGRADE:-}" != "true" ]; then
    if service openitcockpit_agent onestatus >/dev/null 2>&1; then
        service openitcockpit_agent onestop
    fi
fi
exit 0
