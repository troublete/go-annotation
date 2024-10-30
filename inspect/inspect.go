package inspect

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

type Receiver struct {
	ReceiverType string
	Pointer      bool
}

type Function struct {
	Comments []string
	FilePath string
	Name     string
	Package  string
	Receiver *Receiver
}

type FunctionList []Function

func (f FunctionList) Find(name string) *Function {
	for _, f := range f {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

// FindAllFunctions will walk a root dir and search for every function declarations; it will return file,
// function and package information + the doc block and in case of receiver function the corresponding receiver type
// This function is used in a context where a set of file is analysed based on comments to be used for code generation
func FindAllFunctions(root string) (FunctionList, error) {
	var pkgs []map[string]*ast.Package
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			pkg, err := parser.ParseDir(fset, path, func(info fs.FileInfo) bool {
				return true
			}, parser.ParseComments)
			if err != nil {
				return err
			}
			pkgs = append(pkgs, pkg)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var lock sync.Mutex
	var results []Function
	var wg sync.WaitGroup
	wg.Add(len(pkgs))
	for _, p := range pkgs {
		go func(pkglist map[string]*ast.Package) {
			defer wg.Done()
			for pkgname, pkg := range pkglist {
				for fpath, file := range pkg.Files {
					for _, decl := range file.Decls {
						v, ok := decl.(*ast.FuncDecl)
						if ok {
							var recv *Receiver
							if v.Recv != nil {
								recv = &Receiver{}
								def := v.Recv.List[0]
								if t, isPointer := def.Type.(*ast.StarExpr); isPointer {
									recv.Pointer = true
									recv.ReceiverType = t.X.(*ast.Ident).Name
								} else {
									recv.ReceiverType = def.Type.(*ast.Ident).Name
								}
							}

							fname := v.Name.String()
							docs := strings.Split(v.Doc.Text(), "\n")
							var lines []string
							for _, dl := range docs {
								if t := strings.TrimSpace(dl); t != "" {
									lines = append(lines, t)
								}
							}

							doc := Function{
								Comments: lines,
								FilePath: fpath,
								Name:     fname,
								Package:  pkgname,
								Receiver: recv,
							}

							lock.Lock()
							results = append(results, doc)
							lock.Unlock()
						}
					}
				}
			}
		}(p)
	}
	wg.Wait()
	return results, nil
}
