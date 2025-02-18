package singbox

import (
	"github.com/sagernet/sing-box/option"
	"github.com/v42one/airport/pkg/runtime"
)

func WithLogOptions(applyFns ...runtime.ApplyFunc[option.LogOptions]) runtime.ApplyFunc[option.Options] {
	return func(o *option.Options) {
		if o.Log == nil {
			o.Log = new(option.LogOptions)
		}
		runtime.Apply(o.Log, applyFns...)
	}
}

func WithExperimentalOptions(applyFns ...runtime.ApplyFunc[option.ExperimentalOptions]) runtime.ApplyFunc[option.Options] {
	return func(o *option.Options) {
		if o.Experimental == nil {
			o.Experimental = new(option.ExperimentalOptions)
		}
		runtime.Apply(o.Experimental, applyFns...)
	}
}

func WithDNSOptions(fns ...runtime.ApplyFunc[option.DNSOptions]) runtime.ApplyFunc[option.Options] {
	return func(target *option.Options) {
		if target.DNS == nil {
			target.DNS = new(option.DNSOptions)
		}
		runtime.Apply(target.DNS, fns...)
	}
}

func WithRouteOptions(applyFns ...runtime.ApplyFunc[option.RouteOptions]) runtime.ApplyFunc[option.Options] {
	return func(o *option.Options) {
		if o.Route == nil {
			o.Route = new(option.RouteOptions)
		}
		runtime.Apply(o.Route, applyFns...)
	}
}

func WithDefaultDNSRule(applyFns ...runtime.ApplyFunc[option.DefaultDNSRule]) runtime.ApplyFunc[option.DNSRule] {
	return func(r *option.DNSRule) {
		r.Type = "default"
		runtime.Apply(&r.DefaultOptions, applyFns...)
	}
}

func WithDefaultDNSRouteActionOptions(applyFns ...runtime.ApplyFunc[option.DNSRouteActionOptions]) runtime.ApplyFunc[option.DefaultDNSRule] {
	return func(r *option.DefaultDNSRule) {
		r.Action = "route"
		runtime.Apply(&r.RouteOptions, applyFns...)
	}
}

func WithLogicalDNSRule(applyFns ...runtime.ApplyFunc[option.LogicalDNSRule]) runtime.ApplyFunc[option.DNSRule] {
	return func(r *option.DNSRule) {
		r.Type = "logical"
		runtime.Apply(&r.LogicalOptions, applyFns...)
	}
}

func WithLogicalDNSRouteActionOptions(applyFns ...runtime.ApplyFunc[option.DNSRouteActionOptions]) runtime.ApplyFunc[option.LogicalDNSRule] {
	return func(r *option.LogicalDNSRule) {
		r.Action = "route"
		runtime.Apply(&r.RouteOptions, applyFns...)
	}
}

func WithDefaultRule(applyFns ...runtime.ApplyFunc[option.DefaultRule]) runtime.ApplyFunc[option.Rule] {
	return func(r *option.Rule) {
		r.Type = "default"
		runtime.Apply(&r.DefaultOptions, applyFns...)
	}
}

func WithDefaultRouteActionOptions(applyFns ...runtime.ApplyFunc[option.RouteActionOptions]) runtime.ApplyFunc[option.DefaultRule] {
	return func(r *option.DefaultRule) {
		r.Action = "route"
		runtime.Apply(&r.RouteOptions, applyFns...)
	}
}

func WithLogicalRule(applyFns ...runtime.ApplyFunc[option.LogicalRule]) runtime.ApplyFunc[option.Rule] {
	return func(r *option.Rule) {
		r.Type = "logical"
		runtime.Apply(&r.LogicalOptions, applyFns...)
	}
}

func WithLogicalRouteActionOptions(applyFns ...runtime.ApplyFunc[option.RouteActionOptions]) runtime.ApplyFunc[option.LogicalRule] {
	return func(r *option.LogicalRule) {
		r.Action = "route"
		runtime.Apply(&r.RouteOptions, applyFns...)
	}
}
