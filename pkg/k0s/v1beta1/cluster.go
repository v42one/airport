package v1beta1

import (
	k0sctlv1beta1 "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1"
	k0sctlclusterv1beta1 "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1/cluster"
)

const APIVersion = k0sctlv1beta1.APIVersion

type (
	Spec = k0sctlclusterv1beta1.Spec
	K0s  = k0sctlclusterv1beta1.K0s
	Host = k0sctlclusterv1beta1.Host
)

type Cluster struct {
	Kind       string          `yaml:"kind"`
	APIVersion string          `yaml:"apiVersion"`
	Metadata   ClusterMetadata `yaml:"metadata"`
	Spec       *Spec           `yaml:"spec"`
}

type ClusterMetadata struct {
	Name string `yaml:"name"`
}
