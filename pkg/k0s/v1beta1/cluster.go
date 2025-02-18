package v1beta1

import (
	k0sctlv1beta1 "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1"
	k0sctlclusterv1beta1 "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1/cluster"
)

const APIVersion = k0sctlv1beta1.APIVersion

type Spec = k0sctlclusterv1beta1.Spec
type K0s = k0sctlclusterv1beta1.K0s
type Host = k0sctlclusterv1beta1.Host

type Cluster struct {
	Kind       string          `json:"kind"`
	APIVersion string          `yaml:"apiVersion"`
	Metadata   ClusterMetadata `json:"metadata"`
	Spec       *Spec           `yaml:"spec"`
}

type ClusterMetadata struct {
	Name string `yaml:"name"`
}
