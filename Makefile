.PHONY=template
template: cmd/list_templ.go
	templ generate cmd/list.templ

cmd/list_templ.go: cmd/list.templ
	templ generate cmd/list.templ
