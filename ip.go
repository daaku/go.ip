package main

/*
#include <arpa/inet.h>
#include <net/if.h>
#include <net/route.h>
#include <string.h>
#include <stropts.h>
#include <sys/ioctl.h>

void activate_network_device() {
  struct ifreq ifr;
  const char * name = "eth0";
  int fd = socket(PF_INET, SOCK_DGRAM, IPPROTO_IP);

  strncpy(ifr.ifr_name, name, IFNAMSIZ);

  ifr.ifr_addr.sa_family = AF_INET;
  struct sockaddr_in* addr = (struct sockaddr_in*) &ifr.ifr_addr;
  inet_pton(AF_INET, "192.168.100.4", &addr->sin_addr);
  ioctl(fd, SIOCSIFADDR, &ifr);

  inet_pton(AF_INET, "255.255.0.0", &addr->sin_addr);
  ioctl(fd, SIOCSIFNETMASK, &ifr);

  ioctl(fd, SIOCGIFFLAGS, &ifr);
  strncpy(ifr.ifr_name, name, IFNAMSIZ);
  ifr.ifr_flags |= (IFF_UP | IFF_RUNNING);

  ioctl(fd, SIOCSIFFLAGS, &ifr);

  struct rtentry route;
  addr = (struct sockaddr_in*) &route.rt_gateway;
  addr->sin_family = AF_INET;
  addr->sin_addr.s_addr = inet_addr("192.168.100.1");
  addr = (struct sockaddr_in*) &route.rt_dst;
  addr->sin_family = AF_INET;
  addr->sin_addr.s_addr = inet_addr("0.0.0.0");
  addr = (struct sockaddr_in*) &route.rt_genmask;
  addr->sin_family = AF_INET;
  addr->sin_addr.s_addr = inet_addr("0.0.0.0");
  route.rt_flags = RTF_UP | RTF_GATEWAY;
  route.rt_metric = 0;
  ioctl(fd, SIOCADDRT, &route);
}
*/
import "C"

import (
	"fmt"
	"os"
)

func main() {
	if err := do(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func do() error {
	_, err := C.activate_network_device()
	if err != nil {
		return nil
	}

	/*
		s, err := syscall.Socket(
			syscall.AF_INET,
			syscall.SOCK_DGRAM|syscall.SOCK_CLOEXEC,
			syscall.IPPROTO_IP,
		)
		if err != nil {
			return err
		}

		var ifreq *C.ifreq
		ifreq.Name = C.CString("eth0")
		ifreq.Flags = syscall.IFF_UP
		_, _, err = syscall.Syscall(
			syscall.SYS_IOCTL,
			uintptr(s),
			syscall.SIOCGIFFLAGS,
			uintptr(unsafe.Pointer(&ifreq)),
		)
		if err != nil {
			return err
		}
	*/

	/*
		err = syscall.Bind(s, syscall.SockaddrNetlink{Family: syscall.AF_NETLINK})
		if err != nil {
			return err
		}

		sa, err := syscall.Getsockname(s)
		if err != nil {
			return err
		}
	*/

	return nil
}
