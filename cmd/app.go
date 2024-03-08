package cmd

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

const FILTER_PREFIX = "filter="

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Return a list of found files.
func (a *App) getList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a.filterList(ctx, w, r)
}

func (a *App) filterList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if b, err := io.ReadAll(r.Body); err != nil {
		runtime.LogErrorf(a.ctx, "Body reading failed: %v", err)
	} else {
		rawFilter := strings.TrimSpace(string(b))
		filterString := strings.ReplaceAll(strings.Replace(rawFilter, FILTER_PREFIX, "", 1), " ", ".*")
		var pat *regexp.Regexp
		if rePat, err := regexp.Compile(filterString); err != nil {
			runtime.LogErrorf(a.ctx, "regexp.Compile failed: %v", err)
			pat = nil
		} else {
			pat = rePat
		}
		var entries []string
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

					if pat != nil {
						if pat.Match([]byte(path)) {
							runtime.LogDebugf(a.ctx, ">>> Pattern matched, adding %s", path)
							entries = append(entries, strings.TrimSuffix(relative, filepath.Ext(relative)))
						} else {
							runtime.LogDebugf(a.ctx, ">>> Pattern NOT matched, bypass %s", path)
						}
					} else {
						runtime.LogDebugf(a.ctx, ">>> No filter, adding %s", path)
						entries = append(entries, strings.TrimSuffix(relative, filepath.Ext(relative)))
					}
				} else {
					runtime.LogDebugf(a.ctx, ">>> IsFile %s", path)
				}
				return nil
			})
		listComponent := generateList(entries)
		listComponent.Render(a.ctx, w)
		return
	}
}
