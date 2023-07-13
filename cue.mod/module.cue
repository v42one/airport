module: "github.com/v42one/clash-proxy"

require: {
	"github.com/innoai-tech/runtime": "v0.0.0-20221114082425-7a5e0cdc3035"
	"github.com/octohelm/kubepkg":    "v0.5.2"
	"wagon.octohelm.tech":            "v0.0.0-20200202235959-e64a70c55ed2"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
