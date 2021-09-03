package controllers

import (
	"strconv"
	"strings"

	"github.com/morlay/clash-proxy/pkg/clashproxy"

	appsv1 "k8s.io/api/apps/v1"
)

const (
	LabelClashProxyType     = "clash-proxy-type"
	LabelClashProxyCipher   = "clash-proxy-cipher"
	LabelClashProxyPassword = "clash-proxy-password"
	LabelClashProxyPorts    = "clash-proxy-ports"
)

func ProxyOptionFromDeployment(s *appsv1.Deployment) *clashproxy.ProxyDeployment {
	p := &clashproxy.ProxyDeployment{}
	p.Type = s.Labels[LabelClashProxyType]

	podAnnotations := s.Spec.Template.Annotations
	if podAnnotations == nil {
		podAnnotations = map[string]string{}
	}

	p.Cipher = podAnnotations[LabelClashProxyCipher]
	p.Password = podAnnotations[LabelClashProxyPassword]

	for _, portString := range strings.Split(podAnnotations[LabelClashProxyPorts], ",") {
		if portString == "" {
			continue
		}
		port, err := strconv.Atoi(portString)
		if err != nil {
			continue
		}
		p.Ports = append(p.Ports, port)
	}

	return p
}
