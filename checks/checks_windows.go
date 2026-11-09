package checks

func getPlatformChecks() []Check {
	return []Check{
		&CheckMem{},
		&CheckProcess{},
		&CheckAgent{},
		&CheckSwap{},
		&CheckUser{},
		&CheckDisk{},
		&CheckDiskIo{},
		&CheckNet{},
		&CheckNetIo{},
		&CheckCpu{},
		&CheckDocker{},
		&CheckWinService{},
		&CheckWindowsEventLog{},
		&CheckNtp{},
		&CheckLogfile{},
		&CheckCertificate{},
		&CheckFile{},
		&CheckHTTP{},
		&CheckRedfish{},
		&CheckSmart{},
		&CheckGpu{},
		&CheckSecurity{},
		&CheckWindowsScheduledTasks{},
		&CheckWindowsPerfCounter{},
		&CheckNetworkPorts{},
		&CheckIpmi{},
	}
}
