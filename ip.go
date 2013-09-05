// Package ip provides a Linux specific network configuration library.
//
// This is a very limited implementation and currently only handles ipv4 and
// uses the older ioctl APIs rather than the more modern netlink APIs. The only
// scenario this has been tested under is bridged interfaces in a lxc
// container.
package ip

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

type ifreqAddr struct {
	Name [16]byte
	Addr syscall.RawSockaddrInet4
	Pad  [8]byte
}

type ifreqFlags struct {
	Name  [16]byte
	Flags uint16
	Pad   [22]byte
}

type rtentry struct {
	Pad1    [8]byte
	Dst     syscall.RawSockaddrInet4
	Gateway syscall.RawSockaddrInet4
	Genmask syscall.RawSockaddrInet4
	Flags   uint16
	Pad2    [14]byte
	Tos     uint8
	Class   uint8
	Pad3    [3]int16
	Metric  int16
	Pad4    [6]byte
	Dev     uintptr
	Mtu     uint64
	Window  uint64
	Irtt    uint16
	Pad5    [6]byte
}

func ipv4(addr string) (o [4]byte, err error) {
	ip := net.ParseIP(addr).To4()
	if ip == nil {
		err = fmt.Errorf("failed to parse addr %s", addr)
		return
	}
	for i, b := range ip {
		o[i] = b
	}
	return
}

func ioctl(a1, a2, a3 uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, a1, a2, a3)
	if errno != 0 {
		return errno
	}
	return nil
}

// Network device configuration.
type Config struct {
	Device  string
	Addr    string
	Mask    string
	Gateway string
}

func (c Config) deviceBytes() (o [16]byte) {
	for i, c := range c.Device {
		o[i] = byte(c)
	}
	return
}

// Activate the device as configured.
func (c Config) Activate() error {
	devName := c.deviceBytes()

	addr, err := ipv4(c.Addr)
	if err != nil {
		return err
	}

	mask, err := ipv4(c.Mask)
	if err != nil {
		return err
	}

	gateway, err := ipv4(c.Gateway)
	if err != nil {
		return err
	}

	s, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_DGRAM,
		syscall.IPPROTO_IP,
	)
	if err != nil {
		return err
	}
	fd := uintptr(s)

	ifra := ifreqAddr{
		Name: devName,
		Addr: syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
			Addr:   addr,
		},
	}
	if err = ioctl(fd, syscall.SIOCSIFADDR, uintptr(unsafe.Pointer(&ifra))); err != nil {
		return err
	}

	ifra.Addr.Addr = mask
	if err = ioctl(fd, syscall.SIOCSIFNETMASK, uintptr(unsafe.Pointer(&ifra))); err != nil {
		return err
	}

	ifrf := ifreqFlags{Name: devName}
	if err = ioctl(fd, syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifrf))); err != nil {
		return err
	}

	ifrf.Flags = ifrf.Flags | syscall.IFF_UP | syscall.IFF_RUNNING
	if err = ioctl(fd, syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifrf))); err != nil {
		return err
	}

	route := rtentry{
		Dev: uintptr(unsafe.Pointer(&devName)),
		Gateway: syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
			Addr:   gateway,
		},
		Dst: syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
		},
		Genmask: syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
		},
		Flags: syscall.RTF_UP | syscall.RTF_GATEWAY,
	}
	if err = ioctl(fd, syscall.SIOCADDRT, uintptr(unsafe.Pointer(&route))); err != nil {
		return err
	}

	return nil
}
