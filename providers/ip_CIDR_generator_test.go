package providers

import (
	"testing"

	assert "github.com/lexandro/go-assert"
)

func Test_generating_IP_Pool_with_no_range(t *testing.T) {
	generator := IPCIDRGenerator{}
	_, err := generator.Generate("192.168.60.50/32")
	assert.IsNotNil(t, err)
}

func Test_generating_IP_Pool_with_unsupported_range(t *testing.T) {
	generator := IPCIDRGenerator{}
	_, err := generator.Generate("192.168.60.50/1")
	assert.IsNotNil(t, err)
}

func Test_generating_IP_Pool_with_24_range(t *testing.T) {
	generator := IPCIDRGenerator{}
	ips, err := generator.Generate("192.168.60.50/24")
	assert.IsNil(t, err)
	assert.IsTrue(t, len(ips) == 1)
	//fmt.Printf("%s\n", ips[0])
	assert.IsTrue(t, ips[0] == "192.168.60.1/24")
}

func Test_generating_IP_Pool_with_16_range(t *testing.T) {
	generator := IPCIDRGenerator{}
	ips, err := generator.Generate("192.168.60.50/16")
	assert.IsNil(t, err)
	//fmt.Printf("%s\n", ips)
	assert.IsTrue(t, len(ips) == 256)
}

func Test_generating_IP_Pool_with_8_range(t *testing.T) {
	generator := IPCIDRGenerator{}
	ips, err := generator.Generate("192.168.60.50/8")
	assert.IsNil(t, err)
	assert.IsTrue(t, len(ips) == 65536)
}

// func Test_generating_IP_Pool_with_max_range(t *testing.T) {
// 	generator := IPCIDRGenerator{}
// 	ips, err := generator.Generate("192.168.60.50/0")
// 	assert.IsNil(t, err)
// 	assert.IsTrue(t, len(ips) == 16777216)
// }
