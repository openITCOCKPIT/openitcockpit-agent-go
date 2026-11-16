#!/bin/sh
# AIX: register the agent as subsystem of the SRC: runs as root, stopped with SIGTERM (SIGKILL after 30 seconds),
# restarted by the SRC if it terminates abnormally. The agent writes its own log file, panics go to stderr.log.
#   lssrc -s openitcockpit-agent
#   stopsrc -s openitcockpit-agent
#   startsrc -s openitcockpit-agent

if ! lssrc -s openitcockpit-agent >/dev/null 2>&1; then
    mkssys -s openitcockpit-agent -p /usr/bin/openitcockpit-agent -u 0 -S -n 15 -f 9 -w 30 -R -Q \
        -i /dev/null -o /dev/null -e /var/log/openitcockpit-agent/stderr.log
fi

# Start the agent at boot
if ! lsitab oitcagent >/dev/null 2>&1; then
    mkitab "oitcagent:2:once:/usr/bin/startsrc -s openitcockpit-agent >/dev/null 2>&1"
fi

startsrc -s openitcockpit-agent >/dev/null
exit 0
