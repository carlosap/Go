package reflection

import "runtime"

func GetCpuParallelCapabilities() (int, int) {
	var parallel int
	cpus := runtime.NumCPU()
	if cpus < 2 {
		parallel = 1
	} else {
		parallel = cpus - 1
	}

	if runtime.GOOS == "solaris" {
		parallel = 3
	}
	return parallel, cpus
}
