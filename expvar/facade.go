package expvar

func Publish(name string, f VarFunc) {
	Default.Publish(name, f)
}

func PublishMap(f VarMapFunc) {
	Default.PublishMap(f)
}

func Each(f func(key string, val interface{}), sortKeys ...int) {
	Default.Each(f, sortKeys...)
}
