package pkg

import (
	"log"
	"net"
	"strings"
)

func ParseIP(ip string) string {
	if strings.Contains(ip, "/") == true {
		if strings.Contains(ip, "/32") == true {
			nip := strings.Replace(ip, "/32", "", -1)
			address := net.ParseIP(nip)
			if address == nil {
				log.Fatal("illegal ip address")
			}
			return nip
		}
	}
	address := net.ParseIP(ip)
	if address == nil {
		log.Fatal("illegal ip address")
	}
	return ip
}
