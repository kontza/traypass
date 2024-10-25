package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/exp/slog"
)

const CMD_PREFIX = "/wails-go"

type commandHandler func(context.Context, http.ResponseWriter, *http.Request)

func rootRunner(cmd *cobra.Command, args []string) {
	// Create an instance of the app structure
	app := NewApp()
	handlers := map[string]commandHandler{
		CMD_PREFIX + "/get-list": app.getList,
		CMD_PREFIX + "/filter":   app.filterList,
		CMD_PREFIX + "/decrypt":  app.decrypt,
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "traypass",
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
		BackgroundColour: &options.RGBA{R: 29, G: 35, B: 42, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		slog.Error("Error: %v", err)
	}
}
