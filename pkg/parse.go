package pkg

import (
	"log"
	"net"
	"strings"
)

// ip address format: 192.168.10.1 or 192.168.10.1/32
func ParseIP(ip string) string {
	if strings.Contains(ip, "/") == true {
		if strings.Contains(ip, "/32") == true {
			nip := strings.Replace(ip, "/32", "", -1)
			address := net.ParseIP(nip)
			if address == nil {
				log.Fatal("illegal ip address")
			}
			return address.String()
		}
	}
	address := net.ParseIP(ip)
	if address == nil {
		log.Fatal("illegal ip address")
	}
	return address.String()
}
