module: "github.com/v42one/clash-proxy"

require: {
	"dagger.io":                      "v0.3.0"
	"github.com/innoai-tech/runtime": "v0.0.0-20221114082425-7a5e0cdc3035"
	"universe.dagger.io":             "v0.3.0"
	"wagon.octohelm.tech":            "v0.0.0-20200202235959-a41d305b4507"
}

require: {
	"k8s.io/api":          "v0.25.4" @indirect()
	"k8s.io/apimachinery": "v0.25.4" @indirect()
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
