#!/bin/sh
# FreeBSD: pkg does not keep changed configuration files on an update. Like the ports, the package contains
# the files as *.sample and they are copied only if the configuration file does not exist yet.

for sample in /etc/openitcockpit-agent/*.ini.sample; do
    config="${sample%.sample}"
    if [ ! -f "$config" ]; then
        cp -p "$sample" "$config"
    fi
done

mkdir -p /var/log/openitcockpit-agent

# Like other packages of FreeBSD, the package does not enable or restart the agent. pkg kills all processes
# which are started by the scripts of a package, so the agent can not be restarted here after an update.
if service openitcockpit_agent onestatus >/dev/null 2>&1; then
    echo "Restart the openITCOCKPIT Monitoring Agent to use the new version:"
    echo "  service openitcockpit_agent restart"
elif ! service openitcockpit_agent enabled; then
    echo "Start the openITCOCKPIT Monitoring Agent at boot:"
    echo "  sysrc openitcockpit_agent_enable=YES"
    echo "  service openitcockpit_agent start"
fi
exit 0
