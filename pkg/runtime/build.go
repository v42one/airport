package runtime

type ApplyFunc[T any] func(target *T)

func Build[T any](applyFuncs ...ApplyFunc[T]) *T {
	opts := new(T)
	Apply(opts, applyFuncs...)
	return opts
}

func Apply[T any](target *T, applyFns ...ApplyFunc[T]) {
	for _, apply := range applyFns {
		if apply == nil {
			continue
		}
		apply(target)
	}
}

type Applier[T any] interface {
	ApplyTo(x *T)
}

func With[T any](a Applier[T]) ApplyFunc[T] {
	return func(target *T) {
		if a != nil {
			a.ApplyTo(target)
		}
	}
}
