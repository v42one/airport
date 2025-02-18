package kubepkg

import (
	"fmt"
	"path/filepath"

	kubepkgv1alpha1 "github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/x/ptr"
	"github.com/v42one/airport/pkg/runtime"
	"github.com/v42one/airport/pkg/singbox"
)

type SingBox struct {
	Version    string
	ServerName string
	ServerIP   string
	VMess      *singbox.InboundVMess
	Wireguard  *singbox.WireguardEndpoint
}

func (s SingBox) ApplyTo(k *kubepkgv1alpha1.KubePkg) {
	k.SetGroupVersionKind(kubepkgv1alpha1.SchemeGroupVersion.WithKind("KubePkg"))
	k.Name = "sing-box"

	runtime.Apply(&k.Spec, func(spec *kubepkgv1alpha1.Spec) {
		spec.Version = s.Version

		spec.Deploy.SetUnderlying(runtime.Build(func(deploy *kubepkgv1alpha1.DeployDeployment) {
			deploy.Kind = deploy.GetKind()
			deploy.Spec.Replicas = ptr.Ptr(int32(1))
		}))

		spec.Services = map[string]kubepkgv1alpha1.Service{
			"#": *runtime.Build(func(svc *kubepkgv1alpha1.Service) {
				svc.Expose = &kubepkgv1alpha1.Expose{
					Underlying: runtime.Build(func(e *kubepkgv1alpha1.ExposeNodePort) {
						e.Type = e.GetType()
					}),
				}

				svc.Ports = map[string]int32{
					"http":       30100,
					"tcp-main":   int32(s.VMess.ListenPort),
					"tcp-backup": int32(s.VMess.ListenBackupPort),
				}
			}),
		}

		serverConfigMountPath := "/etc/sing-box/config.json"
		clientConfigMountPath := filepath.Join("/usr/share/nginx/html/ui", s.VMess.Secret, "sing-box.json")

		spec.Containers = map[string]kubepkgv1alpha1.Container{
			"sing-box": *runtime.Build(func(c *kubepkgv1alpha1.Container) {
				c.Image.Name = "ghcr.io/sagernet/sing-box"
				c.Image.Tag = fmt.Sprintf("v%s", s.Version)
				c.Args = []string{
					"run", "-c", serverConfigMountPath,
				}
				c.Ports = map[string]int32{
					"tcp-main":   int32(s.VMess.ListenPort),
					"tcp-backup": int32(s.VMess.ListenBackupPort),
				}
			}),
			"config": *runtime.Build(func(c *kubepkgv1alpha1.Container) {
				c.Image.Name = "docker.io/library/nginx"
				c.Image.Tag = "alpine"
				c.Ports = map[string]int32{
					"http": 80,
				}
			}),
		}

		spec.Volumes = map[string]kubepkgv1alpha1.Volume{
			"provider": {
				Underlying: runtime.Build(func(v *kubepkgv1alpha1.VolumeConfigMap) {
					v.Type = "ConfigMap"
					v.MountPath = clientConfigMountPath
					v.SubPath = filepath.Base(clientConfigMountPath)

					data, err := singbox.MarshalJSON(s.ClientConfig())
					if err != nil {
						panic(err)
					}

					v.Spec = runtime.Build(func(spec *kubepkgv1alpha1.ConfigMapSpec) {
						spec.Data = map[string]string{
							v.SubPath: string(data),
						}
					})
				}),
			},
			"config": {
				Underlying: runtime.Build(func(v *kubepkgv1alpha1.VolumeConfigMap) {
					v.Type = "ConfigMap"
					v.MountPath = serverConfigMountPath
					v.SubPath = filepath.Base(serverConfigMountPath)

					data, err := singbox.MarshalJSON(s.ServerConfig())
					if err != nil {
						panic(err)
					}

					v.Spec = runtime.Build(func(spec *kubepkgv1alpha1.ConfigMapSpec) {
						spec.Data = map[string]string{
							v.SubPath: string(data),
						}
					})
				}),
			},
		}
	})
}
