package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

const CMD_PREFIX = "/wails-go"

type commandHandler func(context.Context, http.ResponseWriter, *http.Request)

func getList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	payload := `
	<ol>
		<li>Eka</li>
		<li>Toka</li>
		<li>Kola</li>
	</ol>
	`
	w.Write([]byte(payload))
	runtime.LogInfo(ctx, ">>> List returned")
}

func main() {
	// Create an instance of the app structure
	app := NewApp()
	handlers := map[string]commandHandler{
		CMD_PREFIX + "/get-list": getList,
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "waitwhat",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Middleware: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasPrefix(r.URL.Path, CMD_PREFIX) {
						if handler, ok := handlers[r.URL.Path]; ok {
							handler(app.ctx, w, r)
						} else {
							runtime.LogWarning(app.ctx, fmt.Sprintf("Unknown path: %v", r.URL.Path))
						}
					} else {
						next.ServeHTTP(w, r)
					}
				})
			},
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
