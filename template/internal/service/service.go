package service

var impl *{{ServiceNameInCamelCase}}Impl

type {{ServiceNameInCamelCase}}Impl struct{}

func Setup{{ServiceNameInCamelCase}}Impl() {
	impl = &{{ServiceNameInCamelCase}}Impl{}
}

func Get{{ServiceNameInCamelCase}}Impl() *{{ServiceNameInCamelCase}}Impl {
	return impl
}

func Close{{ServiceNameInCamelCase}}Impl() {}
