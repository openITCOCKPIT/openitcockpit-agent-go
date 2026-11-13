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
	}
}
