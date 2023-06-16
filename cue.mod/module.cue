module: "github.com/v42one/clash-proxy"

require: {
	"dagger.io":                      "v0.3.0"
	"github.com/innoai-tech/runtime": "v0.0.0-20221114082425-7a5e0cdc3035"
	"github.com/octohelm/kubepkg":    "v0.4.1"
	"universe.dagger.io":             "v0.3.0"
	"wagon.octohelm.tech":            "v0.0.0-20200202235959-e64a70c55ed2"
}

require: {
	"k8s.io/api":          "v0.25.4" @indirect()
	"k8s.io/apimachinery": "v0.25.4" @indirect()
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
