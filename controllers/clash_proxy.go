package controllers

import (
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
)

var CLASH_PROXY_CONFIG_MAP_NAME = strings.TrimSpace(os.Getenv("CLASH_PROXY_CONFIG_MAP_NAME"))
var CLASH_PROXY_PUBLIC_IPS = strings.TrimSpace(os.Getenv("CLASH_PROXY_PUBLIC_IPS"))

var publicIPs []string

func init() {
	if CLASH_PROXY_CONFIG_MAP_NAME == "" {
		CLASH_PROXY_CONFIG_MAP_NAME = "ss-proxy"
	}

	if CLASH_PROXY_PUBLIC_IPS != "" {
		if strings.Contains(CLASH_PROXY_PUBLIC_IPS, ",") {
			publicIPs = strings.SplitN(CLASH_PROXY_PUBLIC_IPS, ",", -1)
		} else if strings.Contains(CLASH_PROXY_PUBLIC_IPS, " ") {
			publicIPs = strings.SplitN(CLASH_PROXY_PUBLIC_IPS, " ", -1)
		} else {
			publicIPs = []string{CLASH_PROXY_PUBLIC_IPS}
		}
	}
}

func ProxiesFromDeployment(s *appsv1.Deployment, publicIPs []string) (proxies []Proxy) {
	p := &Proxy{}
	p.Type = s.Labels["clash-proxy-type"]

	podAnnotations := s.Spec.Template.Annotations
	if podAnnotations == nil {
		podAnnotations = map[string]string{}
	}

	p.Cipher = podAnnotations["clash-proxy-cipher"]
	p.Password = podAnnotations["clash-proxy-password"]

	for _, publicIP := range publicIPs {
		for _, port := range strings.Split(podAnnotations["clash-proxy-ports"], ",") {
			if port == "" {
				continue
			}

			p2 := *p
			p2.Port, _ = strconv.Atoi(port)
			p2.Name = fmt.Sprintf("%s-%d", hashID(publicIP), p2.Port)
			p2.Server = publicIP

			proxies = append(proxies, p2)
		}
	}

	return
}

type Proxy struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Cipher   string `yaml:"cipher"`
}

func hashID(v string) string {
	return strconv.FormatUint(uint64(crc32.Checksum([]byte(v), crc32.MakeTable(crc32.IEEE))), 16)
}
