package spx

import (
	"fmt"
	"github.com/goplus/igop"
	"github.com/goplus/igop/gopbuild"
	"go/constant"
	"go/types"
	"sort"
	"strings"
	"unsafe"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/format"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/typesutil"
	"github.com/goplus/mod/gopmod"
	"github.com/goplus/mod/modfile"
	"github.com/goplus/mod/modload"
)

var spxProject = &modfile.Project{
	Ext: ".gmx", Class: "*Game",
	Works:    []*modfile.Class{{Ext: ".spx", Class: "Sprite"}},
	PkgPaths: []string{"github.com/goplus/spx", "math"}}

func StartSPXTypesAnalyser(fileName string, fileCode string) interface{} {
	fset := token.NewFileSet()
	info, err := spxInfo(initSPXMod(), fset, fileName, fileCode, initSPXParserConf())
	if err != nil {
		fmt.Println(err)
	}
	// convert type info to some valid value
	defs := ""
	for k, v := range info.Defs {
		defs += fmt.Sprintf("k: %v, v: %v\n", k, v)
	}
	types := ""
	for k, v := range info.Types {
		types += fmt.Sprintf("k: %v, v: %v\n", k, v)
	}
	instances := ""
	for k, v := range info.Instances {
		instances += fmt.Sprintf("k: %v, v: %v\n", k, v)
	}
	result := map[string]interface{}{
		"Defs":      defs,
		"Types":     types,
		"Instances": instances,
	}
	s := typesList(fset, info.Types)
	for _, v := range s {
		fmt.Println(v)
	}
	return result
}

// init function
func initSPXMod() *gopmod.Module {
	//init spxMod
	var spxMod *gopmod.Module
	spxMod = gopmod.New(modload.Default)
	spxMod.Opt.Projects = append(spxMod.Opt.Projects, spxProject)
	err := spxMod.ImportClasses()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return spxMod
}

// init function
func initSPXParserConf() parser.Config {
	return parser.Config{
		ClassKind: func(fname string) (isProj bool, ok bool) {
			ext := modfile.ClassExt(fname)
			c, ok := lookupClass(ext)
			if ok {
				isProj = c.IsProj(ext, fname)
			}
			return
		},
	}
}

// check function
func lookupClass(ext string) (c *modfile.Project, ok bool) {
	switch ext {
	case ".gmx", ".spx":
		return spxProject, true
	}
	return
}

func spxInfo(mod *gopmod.Module, fileSet *token.FileSet, fileName string, fileCode string, parseConf parser.Config) (*typesutil.Info, error) {
	// new parser
	file, err := parser.ParseEntry(fileSet, fileName, fileCode, parseConf)
	if err != nil {
		return nil, err
	}
	// init types conf
	ctx := igop.NewContext(0)
	c := gopbuild.NewContext(ctx)
	//TODO: igop
	conf := &types.Config{}
	// replace it!
	conf.Importer = c
	chkOpts := &typesutil.Config{
		Types:                 types.NewPackage("main", file.Name.Name),
		Fset:                  fileSet,
		Mod:                   mod,
		UpdateGoTypesOverload: false,
	}

	// init info
	info := &typesutil.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
		Overloads:  make(map[*ast.Ident][]types.Object),
	}
	check := typesutil.NewChecker(conf, chkOpts, nil, info)
	err = check.Files(nil, []*ast.File{file})
	return info, err
}

func typesList(fset *token.FileSet, types map[ast.Expr]types.TypeAndValue) []string {
	var items []string
	for expr, tv := range types {
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		tvstr := tv.Type.String()
		if tv.Value != nil {
			tvstr += " = " + tv.Value.String()
		}
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%3d:%3d | %-19s %-40T | %-14s : %s | %v",
			posn.Line, posn.Column, exprString(fset, expr), expr,
			mode(tv), tvstr, (*TypeAndValue)(unsafe.Pointer(&tv)).mode)
		items = append(items, buf.String())
	}
	sort.Strings(items)
	return items
}

type operandMode byte

type TypeAndValue struct {
	mode  operandMode
	Type  types.Type
	Value constant.Value
}

func mode(tv types.TypeAndValue) string {
	switch {
	case tv.IsVoid():
		return "void"
	case tv.IsType():
		return "type"
	case tv.IsBuiltin():
		return "builtin"
	case tv.IsNil():
		return "nil"
	case tv.Assignable():
		if tv.Addressable() {
			return "var"
		}
		return "mapindex"
	case tv.IsValue():
		return "value"
	default:
		return "unknown"
	}
}

func exprString(fset *token.FileSet, expr ast.Expr) string {
	var buf strings.Builder
	format.Node(&buf, fset, expr)
	return buf.String()
}
