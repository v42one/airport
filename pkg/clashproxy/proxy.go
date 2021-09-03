package clashproxy

import (
	"hash/crc32"
	"strconv"
)

type ProxyYAML struct {
	Proxies []Proxy `yaml:"proxies"`
}

type ProxyDeployment struct {
	Name     string
	Type     string
	Cipher   string
	Password string
	Ports    []int
}

type Proxy struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Cipher   string `yaml:"cipher"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
}

func HashID(v string) string {
	return strconv.FormatUint(uint64(crc32.Checksum([]byte(v), crc32.MakeTable(crc32.IEEE))), 16)
}
