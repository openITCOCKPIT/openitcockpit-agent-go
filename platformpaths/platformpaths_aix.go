package platformpaths

type aixPlatformPath struct {
}

func (p *aixPlatformPath) Init() error {
	return nil
}

func (p *aixPlatformPath) LogPath() string {
	return "/var/log/openitcockpit-agent/agent.log"
}

func (p *aixPlatformPath) ConfigPath() string {
	return "/etc/openitcockpit-agent/"
}

func (p *aixPlatformPath) AdditionalData() map[string]string {
	return map[string]string{}
}

func getPlatformPath() PlatformPath {
	return &aixPlatformPath{}
}
