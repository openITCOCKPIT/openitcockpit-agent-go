#!/bin/sh
# AIX: stop the agent before an update (the agent runs as subsystem of the SRC)

if lssrc -s openitcockpit-agent >/dev/null 2>&1; then
    if lssrc -s openitcockpit-agent | grep -w active >/dev/null; then
        stopsrc -s openitcockpit-agent >/dev/null
        waited=0
        while lssrc -s openitcockpit-agent | grep -w active >/dev/null && [ $waited -lt 30 ]; do
            sleep 1
            waited=$((waited + 1))
        done
    fi
fi
exit 0
