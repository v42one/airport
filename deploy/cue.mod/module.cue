module: "github.com/morlay/clash-proxy/deploy"

require: {
	"github.com/octohelm/cuem": "v0.0.0-20210520091405-7e9ddaa903c7"
}

require: {
	"k8s.io/api":          "v0.22.1" @indirect()
	"k8s.io/apimachinery": "v0.22.1" @indirect()
}
