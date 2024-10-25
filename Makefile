cmd/list_templ.go: cmd/list.templ
	templ generate cmd/list.templ

.PHONY=build
build: build/bin/traypass.app

build/bin/traypass.app: main.go cmd/*.go
	wails build -ldflags '-w -s'

.PHONY=dev
dev: APPARGS=-appargs '-v'
dev:
	wails dev $(APPARGS)
