#!/bin/sh
# AIX: remove the subsystem of the SRC and the entry of /etc/inittab
# $1 is 0 if the package is removed and 1 on an update (the new package stops and starts the agent)

if [ "$1" = "0" ]; then
    if lssrc -s openitcockpit-agent >/dev/null 2>&1; then
        stopsrc -s openitcockpit-agent >/dev/null 2>&1
        waited=0
        while lssrc -s openitcockpit-agent | grep -w active >/dev/null && [ $waited -lt 30 ]; do
            sleep 1
            waited=$((waited + 1))
        done
        rmssys -s openitcockpit-agent >/dev/null
    fi
    if lsitab oitcagent >/dev/null 2>&1; then
        rmitab oitcagent
    fi
fi
exit 0
