package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/troublete/go-chariot/analyze"
	"github.com/troublete/go-chariot/auxiliary"
	"github.com/troublete/go-chariot/cmd/chariot/internal"
	"github.com/troublete/go-chariot/inspect"
	"log/slog"
	"os"
	"strings"
	"sync"
	"text/template"
)

var (
	routeTpl       = `http.Handle("{{.Path}}", chariot.HTTPHandler({{.Package}}.{{.Name}}))`
	routingFileTpl = `
package main

func routing() {
	{{- range .Routes -}}
		{{.}}
	{{end -}}
}
`
)

func main() {
	path := os.Args[1]
	flag.Parse()

	funcs, err := inspect.FindAllFunctions(path)
	if err != nil {
		panic(err)
	}

	var routingFileMu sync.Mutex
	var routes []string

	pc := internal.NewProcessorRegister()
	pc.Register("chariot.route", func(args map[string]string, f inspect.Function) error {
		ac := auxiliary.AttributeContract(args, map[string]func(v string) bool{
			"path": func(_ string) bool { return true },
		})

		if !ac.Valid {
			return fmt.Errorf("failed to execute chariot.route, attributes invalid please check docs")
		}

		sb := &strings.Builder{}
		tpl := template.Must(template.New(".").Parse(routeTpl))
		err := tpl.Execute(sb, map[string]string{
			"Path":    ac.Attributes["path"],
			"Package": f.Package,
			"Name":    f.Name,
		})
		if err != nil {
			return err
		}

		routingFileMu.Lock()
		defer routingFileMu.Unlock()

		routes = append(routes, sb.String())
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for _, f := range funcs {
		go func() {
			defer wg.Done()
			defs, warns := analyze.ExtractDefinitionsForFunction(f, analyze.FilterCommentNotInForm())
			if len(warns) > 0 {
				for _, w := range warns {
					slog.Warn(w)
				}
			} else {
				for _, d := range defs {
					if proc, ok := pc[d.Identifier]; !ok {
						slog.Error("failed to lookup processor", "proc", proc)
						os.Exit(1)
					} else {
						err := proc(d.Arguments, f)
						if err != nil {
							slog.Error("failed to process", "proc", proc, "args", d.Arguments, "func", f.Name)
						}
					}
				}
			}
		}()
	}
	wg.Wait()

	var buf bytes.Buffer
	err = template.Must(template.New(".").Parse(routingFileTpl)).Execute(&buf, map[string]any{
		"Routes": routes,
	})
	if err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fn := fmt.Sprintf("%v/routes.go", cwd)
	slog.Info("writing routing file", "file", fn)
	err = auxiliary.WriteGoFile(fn, buf.Bytes())
	if err != nil {
		panic(err)
	}
}
