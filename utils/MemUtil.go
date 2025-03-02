// Neo

package utils

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func IsLinuxSystem() bool {
	return "linux" == runtime.GOOS
}

func MemoryUsageInKB() (kb int64, err error) {
	if IsLinuxSystem() {
		path := fmt.Sprintf("/proc/%d/status", os.Getpid())
		var f *os.File
		f, err = os.OpenFile(path, os.O_RDONLY, 0444)
		if nil != err {
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if 0 == strings.Index(line, "VmRSS") {
				_, err = fmt.Sscanf(line, "VmRSS:\t%d kB", &kb)
				return
			}
		}
	}
	return
}
