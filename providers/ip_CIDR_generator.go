package providers

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type IPCIDRGenerator struct {
}

func (g *IPCIDRGenerator) Generate(baseIp string) ([]string, error) {
	if baseIp == "" {
		return nil, errors.New("The parameter: \"baseIp\" MUST be set.")
	}
	ip, ipnet, err := net.ParseCIDR(baseIp)
	if err != nil {
		return nil, fmt.Errorf("Error occured while parsing you passed argument: ip-range (%s), error: %s", baseIp, err.Error())
	}
	ipchars := strings.Split(ip.String(), ".")
	ones, _ := ipnet.Mask.Size()
	if ones%8 != 0 {
		return nil, errors.New("We have not supoorted any others network masks which % 8 != 0")
	}
	if ones == 32 {
		return nil, errors.New("No more available sub network can be allocates!")
	}
	if ones == 24 {
		return []string{fmt.Sprintf("%s.%s.%s.1/24", ipchars[0], ipchars[1], ipchars[2])}, nil
	}
	var ips []string
	if ones == 16 {
		for i := 0; i < 256; i++ {
			ips = append(ips, fmt.Sprintf("%s.%s.%d.1/24", ipchars[0], ipchars[1], i))
		}
		return ips, nil
	}
	if ones == 8 {
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				ips = append(ips, fmt.Sprintf("%s.%d.%d.1/24", ipchars[0], i, j))
			}
		}
		return ips, nil
	}
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			for k := 0; k < 256; k++ {
				ips = append(ips, fmt.Sprintf("%d.%d.%d.1/24", i, j, k))
			}
		}
	}
	return ips, nil
}
