package inspect

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

type Receiver struct {
	ReceiverType string `json:"receiver_type"`
	Pointer      bool   `json:"is_pointer"`
}

type Function struct {
	Comments []string  `json:"-"`
	FilePath string    `json:"file_path"`
	Name     string    `json:"name"`
	Package  string    `json:"package"`
	Receiver *Receiver `json:"receiver,omitempty"`
}

func (f Function) Doc() []string {
	return f.Comments
}

type FunctionList []Function

// Find returns a function by name
// Used to quickly lookup a function
func (f FunctionList) Find(name string) *Function {
	for _, f := range f {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

// FindAllFunctions uses the go parser to traverse (starting on root) all valid go files and extract
// all functions + comments found (all functions declarations including the ones defined on receivers)
// it returns a simplified representation of everything found
// To be used to use go code as metaprogramming input for code generation and similar
// functions
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
						f, fok := decl.(*ast.FuncDecl)
						if fok {
							var recv *Receiver
							if f.Recv != nil {
								recv = &Receiver{}
								def := f.Recv.List[0]
								if t, isPointer := def.Type.(*ast.StarExpr); isPointer {
									recv.Pointer = true
									recv.ReceiverType = t.X.(*ast.Ident).Name
								} else {
									recv.ReceiverType = def.Type.(*ast.Ident).Name
								}
							}

							docs := strings.Split(f.Doc.Text(), "\n")
							var lines []string
							for _, dl := range docs {
								if t := strings.TrimSpace(dl); t != "" {
									lines = append(lines, t)
								}
							}

							doc := Function{
								Comments: lines,
								FilePath: fpath,
								Name:     f.Name.String(),
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

type Type struct {
	Comments []string `json:"-"`
	FilePath string   `json:"file_path"`
	Name     string   `json:"name"`
	Package  string   `json:"package"`
	Fields   []Field  `json:"fields"`
}

func (t Type) Doc() []string {
	return t.Comments
}

type Field struct {
	Comments []string          `json:"-"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Tags     map[string]string `json:"tags"`
}

func (f Field) Doc() []string {
	return f.Comments
}

type TypeList []Type

// Find searches a type with name and returns a pointer or nil
// Used to quickly find a type
func (tl TypeList) Find(name string) *Type {
	for _, t := range tl {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

// FindAllTypes uses the go parser to traverse (starting on root) all valid go files and extract
// all types + comments found (so all type declarations) alongside their fields + comments
// it returns a simplified representation of everything found
// To be used to use go code as metaprogramming input for code generation and similar
// functions
func FindAllTypes(root string) (TypeList, error) {
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
	var results []Type
	var wg sync.WaitGroup
	wg.Add(len(pkgs))
	for _, p := range pkgs {
		go func(pkglist map[string]*ast.Package) {
			defer wg.Done()
			for pkgname, pkg := range pkglist {
				for fpath, file := range pkg.Files {
					for _, decl := range file.Decls {
						g, gok := decl.(*ast.GenDecl)
						if gok {
							docs := strings.Split(g.Doc.Text(), "\n")
							var lines []string
							for _, dl := range docs {
								if t := strings.TrimSpace(dl); t != "" {
									lines = append(lines, t)
								}
							}

							for _, s := range g.Specs {
								t, tok := s.(*ast.TypeSpec)
								if tok {
									var fields []Field
									s, sok := t.Type.(*ast.StructType)
									if sok {
										if s.Fields != nil {
											for _, f := range s.Fields.List {
												docs := strings.Split(f.Doc.Text(), "\n")
												var lines []string
												for _, dl := range docs {
													if t := strings.TrimSpace(dl); t != "" {
														lines = append(lines, t)
													}
												}

												var tags map[string]string
												if f.Tag != nil {
													tags = map[string]string{}
													pairs := strings.Split(f.Tag.Value[1:len(f.Tag.Value)-1], ",")
													for _, p := range pairs {
														kv := strings.Split(p, ":")
														tags[kv[0]] = kv[1]
													}
												}

												st, stok := f.Type.(*ast.Ident)
												if stok {
													fields = append(fields, Field{
														Comments: lines,
														Name:     f.Names[0].String(),
														Type:     st.Name,
														Tags:     tags,
													})
												}

												set, setok := f.Type.(*ast.SelectorExpr)
												if setok {
													fields = append(fields, Field{
														Comments: lines,
														Name:     f.Names[0].String(),
														Type:     fmt.Sprintf("%v.%v", set.X.(*ast.Ident).Name, set.Sel.Name),
														Tags:     tags,
													})
												}

												stet, stetok := f.Type.(*ast.StarExpr)
												if stetok {
													fields = append(fields, Field{
														Comments: lines,
														Name:     f.Names[0].String(),
														Type:     fmt.Sprintf("*%v.%v", stet.X.(*ast.SelectorExpr).X.(*ast.Ident).Name, stet.X.(*ast.SelectorExpr).Sel.Name),
														Tags:     tags,
													})
												}
											}
										}
									}

									doc := Type{
										Comments: lines,
										FilePath: fpath,
										Name:     t.Name.String(),
										Package:  pkgname,
										Fields:   fields,
									}
									lock.Lock()
									results = append(results, doc)
									lock.Unlock()
								}
							}
						}
					}
				}
			}
		}(p)
	}
	wg.Wait()
	return results, nil
}
