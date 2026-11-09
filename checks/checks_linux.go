package checks

func getPlatformChecks() []Check {
	checks := []Check{
		&CheckMem{},
		&CheckProcess{},
		&CheckAgent{},
		&CheckSwap{},
		&CheckUser{},
		&CheckDisk{},
		&CheckDiskIo{},
		&CheckLoad{},
		&CheckCpu{},
		&CheckNet{},
		&CheckNetIo{},
		&CheckSensor{},
		&CheckDocker{},
		&CheckSystemd{},
		&CheckNtp{},
		&CheckLogfile{},
		&CheckCertificate{},
		&CheckFile{},
		&CheckHTTP{},
		&CheckRedfish{},
		&CheckSmart{},
		&CheckGpu{},
		&CheckSecurity{},
		&CheckZfs{},
		&CheckCeph{},
		&CheckProxmox{},
		&CheckPostfix{},
		&CheckExim{},
		&CheckLvm{},
		&CheckMounts{},
		&CheckNetworkPorts{},
		&CheckIpmi{},
		&CheckSystemHealth{},
	}

	// Only contains the libvirt check if the agent got build with the "libvirt" build tag
	return append(checks, libvirtChecks()...)
}
