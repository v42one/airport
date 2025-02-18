package singbox

import "strings"

type Addr struct {
	Scheme string
	Host   string
	Path   string
}

func (a Addr) String() string {
	b := strings.Builder{}

	if a.Scheme != "" {
		b.WriteString(a.Scheme)
		b.WriteString("://")
	}

	b.WriteString(a.Host)

	if a.Path != "" {
		b.WriteString("/")
		b.WriteString(a.Path)
	}

	return b.String()
}
