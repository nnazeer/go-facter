package host

import (
	"os"
	"os/user"
	"strings"
	"syscall"

	h "github.com/shirou/gopsutil/host"
)

type Facter interface {
	Add(string, interface{})
}

// int8ToString converts [65]int8 in syscall.Utsname to string
func int8ToString(bs [65]int8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}
	return string(b)
}

func GetHostFacts(f Facter) error {
	hostInfo, err := h.Info()
	if err != nil {
		return err
	}
	// TODO - capitalize the first letter of kernel and OS
	f.Add("fqdn", hostInfo.Hostname)
	splitted := strings.SplitN(hostInfo.Hostname, ".", 2)
	var hostname *string
	if len(splitted) > 1 {
		hostname = &splitted[0]
		f.Add("domain", splitted[1])
	} else {
		hostname = &hostInfo.Hostname
	}
	f.Add("hostname", *hostname)

	var is_virtual bool
	if hostInfo.VirtualizationRole == "host" {
		is_virtual = false
	} else {
		is_virtual = true
	}
	f.Add("is_virtual", is_virtual)

	f.Add("kernel", hostInfo.OS)
	f.Add("operatingsystemrelease", hostInfo.PlatformVersion)
	f.Add("operatingsystem", hostInfo.Platform)
	f.Add("osfamily", hostInfo.PlatformFamily)
	f.Add("uptime_seconds", hostInfo.Uptime)
	f.Add("uptime_minutes", hostInfo.Uptime/60)
	f.Add("uptime_hours", hostInfo.Uptime/60/60)
	f.Add("uptime_days", hostInfo.Uptime/60/60/24)
	f.Add("virtual", hostInfo.VirtualizationSystem)

	envPath := os.Getenv("PATH")
	if envPath != "" {
		f.Add("path", envPath)
	}

	user, err := user.Current()
	if err == nil {
		f.Add("id", user.Username)
	} else {
		panic(err)
	}

	var uname syscall.Utsname
	err = syscall.Uname(&uname)
	if err == nil {
		kernelRelease := int8ToString(uname.Release)
		f.Add("kernelrelease", kernelRelease)
		f.Add("kernelversion", strings.Split(kernelRelease, "-")[0])
		f.Add("hardwaremodel", int8ToString(uname.Machine))
	}

	return nil
}