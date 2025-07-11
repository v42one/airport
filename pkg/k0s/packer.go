package k0s

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	k0sctlclusterv1beta1 "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1/cluster"
	"github.com/k0sproject/rig"
	k0sprojectversion "github.com/k0sproject/version"
	kubepkgv1alpha1 "github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/manifest"
	"github.com/octohelm/unifs/pkg/filesystem"
	"github.com/octohelm/x/ptr"
	"go.yaml.in/yaml/v3"
	"sigs.k8s.io/yaml/kyaml"

	"github.com/v42one/airport/pkg/k0s/v1beta1"
	"github.com/v42one/airport/pkg/runtime"
)

const (
	manifestsHostBase  = "/var/lib/k0s/manifests"
	manifestsLocalBase = "manifests"
)

type Cluster struct {
	Name         string
	K0sVersion   string
	RemoteServer string
	Components   []*kubepkgv1alpha1.KubePkg
}

func toKYAML(raw []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(nil)
	if err := (&kyaml.Encoder{}).FromYAML(bytes.NewBuffer(raw), b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (clt *Cluster) ExportTo(ctx context.Context, fsys filesystem.FileSystem) error {
	clusterYamlRaw, err := toKYAML(yaml.Marshal(clt.toCluster()))
	if err != nil {
		return fmt.Errorf("marshal cluster.kyaml failed: %w", err)
	}

	if err := saveTo(ctx, fsys, filepath.Join(clt.Name, "cluster.yaml"), clusterYamlRaw); err != nil {
		return err
	}

	for _, kpkg := range clt.Components {
		manifests, err := manifest.SortedExtract(kpkg)
		if err != nil {
			return err
		}

		b := bytes.NewBuffer(nil)
		enc := &kyaml.Encoder{}

		for _, m := range manifests {
			if err := enc.FromObject(m, b); err != nil {
				return fmt.Errorf("marshal to kyaml failed: %w", err)
			}
		}

		if err := saveTo(ctx,
			fsys,
			filepath.Join(clt.Name, manifestsLocalBase, kpkg.Name, fmt.Sprintf("%s.yaml", kpkg.Name)),
			b.Bytes(),
		); err != nil {
			return err
		}

	}

	return nil
}

func (clt *Cluster) toCluster() *v1beta1.Cluster {
	return runtime.Build(func(c *v1beta1.Cluster) {
		c.Kind = "Cluster"
		c.APIVersion = v1beta1.APIVersion
		c.Spec = clt.toClusterSpec()
		c.Metadata.Name = clt.Name
	})
}

func (clt *Cluster) toClusterSpec() *v1beta1.Spec {
	return runtime.Build(func(spec *v1beta1.Spec) {
		spec.K0s = runtime.Build(func(k0s *v1beta1.K0s) {
			k0s.Version, _ = k0sprojectversion.NewVersion(clt.K0sVersion)
		})

		spec.Hosts = append(spec.Hosts, runtime.Build(func(host *v1beta1.Host) {
			host.Metadata.Hostname = clt.Name
			host.Role = "single"
			host.InstallFlags = []string{
				"--disable-components", "helm",
			}

			host.SSH = runtime.Build(func(ssh *rig.SSH) {
				ssh.Address = clt.RemoteServer
				ssh.Port = 22
				ssh.KeyPath = ptr.Ptr("~/.ssh/id_rsa")
				ssh.User = "root"
			})

			for _, x := range clt.Components {
				distYAML := filepath.Join(manifestsHostBase, x.Name, fmt.Sprintf("%s.yaml", x.Name))

				host.Files = append(host.Files, runtime.Build(func(uploadFile *k0sctlclusterv1beta1.UploadFile) {
					uploadFile.Name = "manifests-" + x.Name
					uploadFile.Source = filepath.Join(manifestsLocalBase, x.Name, fmt.Sprintf("%s.yaml", x.Name))
					uploadFile.DestinationFile = distYAML
				}))
			}
		}))
	})
}

func saveTo(ctx context.Context, fsys filesystem.FileSystem, filename string, data []byte) error {
	if err := filesystem.MkdirAll(ctx, fsys, filepath.Dir(filename)); err != nil {
		return err
	}
	f, err := fsys.OpenFile(ctx, filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}
