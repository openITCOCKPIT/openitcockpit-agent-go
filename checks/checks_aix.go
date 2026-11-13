package checks

func getPlatformChecks() []Check {
	return []Check{
		&CheckMem{},
		&CheckAgent{},
		&CheckSwap{},
		&CheckDisk{},
		&CheckLoad{},
		&CheckCpu{},
		&CheckLpar{},
		&CheckSrcServices{},
		&CheckErrpt{},
		&CheckVolumeGroups{},
		&CheckMultipath{},
		&CheckLogfile{},
		&CheckCertificate{},
		&CheckFile{},
		&CheckHTTP{},
		&CheckRedfish{},
		&CheckPostfix{},
	}
}
