package cmd

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	listTemplate *template.Template
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	var err error
	a.listTemplate, err = template.New("entries").Parse(`<select name="entries", id="secret", size="100">
		{{range .}}
			<option value="{{.}}">{{.}}</option>
		{{end}}
		</select>`)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Template parsing failed: %v", err)
	}
}

// Return a list of found files.
func (a *App) getList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	entries := []string{}
	expanded := os.ExpandEnv(appConfig.ScanDirectory)
	filepath.WalkDir(expanded,
		func(path string, de fs.DirEntry, err error) error {
			if de.IsDir() {
				runtime.LogDebugf(a.ctx, ">>> IsDir %s", path)
				if strings.HasSuffix(de.Name(), ".git") {
					runtime.LogDebugf(a.ctx, ">>> .git %s", path)
					return filepath.SkipDir
				}
				return nil
			}
			if filepath.Ext(de.Name()) == ".gpg" {
				var relative string
				if relative, err = filepath.Rel(expanded, path); err != nil {
					return err
				}

				entries = append(entries, strings.TrimSuffix(relative, filepath.Ext(relative)))
			} else {
				runtime.LogDebugf(a.ctx, ">>> IsFile %s", path)
			}
			return nil
		})
	if err := a.listTemplate.Execute(w, entries); err != nil {
		runtime.LogErrorf(a.ctx, "Template execution failed: %v", err)
		return
	}
}
