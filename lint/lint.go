package lint

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

type Trifle struct {
	Position token.Position
	Text     string
}

func Run(conf Config, target string, includeTest bool) ([]*Trifle, error) {
	matters := conf.matterByTarget()

	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.LoadSyntax | packages.LoadTypes,
		Tests: includeTest,
	}, target)
	if err != nil {
		return nil, errors.Wrap(err, "load package failed")
	}

	var trifles []*Trifle
	for _, pkg := range pkgs {
		matter, ok := getMatter(pkg.PkgPath, matters)
		if !ok {
			continue
		}

		trifles = append(trifles, checkImports(pkg.Fset, matter.Import, pkg)...)
		trifles = append(trifles, checkRefs(pkg.Fset, matter.Ref, pkg)...)
		trifles = append(trifles, checkBuiltin(pkg.Fset, matter.Write, pkg)...)
	}

	sort.Slice(trifles, func(i, j int) bool {
		pi, pj := trifles[i].Position, trifles[j].Position

		if pi.Filename != pj.Filename {
			return pi.Filename < pj.Filename
		}
		if pi.Line != pj.Line {
			return pi.Line < pj.Line
		}
		return pi.Column < pj.Column
	})

	return trifles, nil
}

func getMatter(pkgPath string, matters map[string]*Matter) (*Matter, bool) {
	matter, ok := matters[pkgPath]
	if ok {
		return matter, true
	}
	matter, ok = matters[globalTarget]
	if ok {
		return matter, true
	}
	return nil, false
}

func checkImports(fs *token.FileSet, list []string, pkg *packages.Package) []*Trifle {
	var trifles []*Trifle
	for _, file := range pkg.Syntax {
		for _, spec := range file.Imports {
			for _, t := range list {
				if val := strings.Trim(spec.Path.Value, "\""); t != val {
					continue
				}

				trifles = append(trifles, &Trifle{
					Position: fs.Position(spec.Pos()),
					Text:     fmt.Sprintf("forbidden to import %s", t),
				})
			}
		}
	}

	return trifles
}

func checkRefs(fs *token.FileSet, list map[string][]string, pkg *packages.Package) []*Trifle {
	var trifles []*Trifle

	mp := make(map[string]map[string]struct{})
	for p, l := range list {
		if _, ok := mp[p]; !ok {
			mp[p] = make(map[string]struct{})
		}
		for _, t := range l {
			mp[p][t] = struct {
			}{}
		}
	}

	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			sx, ok := node.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pathname := pkg.TypesInfo.ObjectOf(sx.Sel).Pkg().Path()
			list, ok := mp[pathname]
			if !ok {
				return true
			}

			if _, ok := list[sx.Sel.Name]; ok {
				trifles = append(trifles, &Trifle{
					Position: fs.Position(sx.Sel.Pos()),
					Text:     fmt.Sprintf("forbidden to refer '%s'.%s", pathname, sx.Sel.Name),
				})
			}
			return true
		})
	}

	return trifles
}

func checkBuiltin(fs *token.FileSet, list []string, pkg *packages.Package) []*Trifle {
	var trifles []*Trifle

	mp := toMap(list)

	for _, syntax := range pkg.Syntax {
		ast.Inspect(syntax, func(node ast.Node) bool {
			var pos token.Position
			var text string

			switch x := node.(type) {
			case *ast.IfStmt:
				if _, ok := mp["if"]; !ok {
					return true
				}
				pos, text = fs.Position(x.Pos()), "forbidden to write if"
			case *ast.SwitchStmt:
				if _, ok := mp["switch"]; !ok {
					return true
				}
				pos, text = fs.Position(x.Pos()), "forbidden to write switch"
			case *ast.Ident:
				if _, ok := mp["panic"]; !ok {
					return true
				}
				if x.Name != "panic" || pkg.TypesInfo.ObjectOf(x).Pkg() != nil {
					return true
				}
				pos, text = fs.Position(x.Pos()), "forbidden to write panic"
			case *ast.ForStmt, *ast.RangeStmt:
				if _, ok := mp["for"]; !ok {
					return true
				}
				pos, text = fs.Position(x.Pos()), "forbidden to write for"
			default:
				return true
			}

			trifles = append(trifles, &Trifle{Position: pos, Text: text})

			return true
		})
	}

	return trifles
}

func toMap(sl []string) map[string]struct{} {
	sm := make(map[string]struct{})
	for _, s := range sl {
		sm[s] = struct{}{}
	}
	return sm
}
