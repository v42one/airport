package clashproxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"gopkg.in/yaml.v2"
)

func ParsePublicIPs(s string) *Provider {
	p := &Provider{}

	if strings.Contains(s, ",") {
		p.PublicIPs = strings.SplitN(s, ",", -1)
	} else if strings.Contains(s, " ") {
		p.PublicIPs = strings.SplitN(s, " ", -1)
	} else {
		p.PublicIPs = []string{s}
	}

	return p
}

// Provider public ip list
type Provider struct {
	PublicIPs        []string
	proxyDeployments sync.Map
}

func (p *Provider) LoadPublicIPsFromKube(ctx context.Context, c client.Reader) error {
	nl := &corev1.NodeList{}

	if err := c.List(ctx, nl); err != nil {
		return err
	}

	publicIPs := make([]string, 0, len(nl.Items))

	for _, i := range nl.Items {
		for _, a := range i.Status.Addresses {
			if a.Type == corev1.NodeExternalIP {
				publicIPs = append(publicIPs, a.Address)
			}
		}
	}

	p.PublicIPs = publicIPs

	return nil
}

func (p *Provider) Add(name string, pd *ProxyDeployment) {
	pd.Name = name
	p.proxyDeployments.Store(name, pd)
}

func (p *Provider) Remove(name string) {
	p.proxyDeployments.Delete(name)
}

func (p *Provider) ToProxyYAML() *ProxyYAML {
	py := &ProxyYAML{}

	p.proxyDeployments.Range(func(key, value interface{}) bool {
		pd := value.(*ProxyDeployment)

		for _, publicIP := range p.PublicIPs {
			for _, port := range pd.Ports {
				py.Proxies = append(py.Proxies, Proxy{
					Server: publicIP,
					Port:   port,

					Name: fmt.Sprintf("%s-%d", HashID(publicIP), port),
					Type: pd.Type,

					Cipher:   pd.Cipher,
					Password: pd.Password,
				})
			}
		}

		return true
	})

	return py
}

func (p *Provider) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	if request.RequestURI != "/proxy.yaml" {
		rw.WriteHeader(http.StatusNotFound)
		_, _ = rw.Write(nil)
		return
	}

	rw.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	_ = yaml.NewEncoder(rw).Encode(p.ToProxyYAML())
}

func (p *Provider) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    ":80",
		Handler: p,
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println(err)
			} else {
				log.Fatalln(err)
			}
		}
	}()

	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
