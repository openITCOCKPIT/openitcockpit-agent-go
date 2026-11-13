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
		&CheckLoad{},
		&CheckCpu{},
		&CheckNet{},
		&CheckNetIo{},
		&CheckNtp{},
		&CheckLpar{},
		&CheckSrcServices{},
		&CheckErrpt{},
		&CheckVolumeGroups{},
		&CheckMultipath{},
		&CheckNetworkAdapters{},
		&CheckPagingSpaces{},
		&CheckLogfile{},
		&CheckCertificate{},
		&CheckFile{},
		&CheckHTTP{},
		&CheckRedfish{},
		&CheckPostfix{},
		&CheckMounts{},
		&CheckNetworkPorts{},
	}
}
