package configs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	ConfigSrcFileName = "config.go"
	ConfigDstFileName = "config_gen.go"
)

func Generate(rootDir string) error {
	pkg, err := findPackage(rootDir)
	if err != nil {
		return err
	}

	err = joinPackageErrors(pkg)
	if err != nil {
		return err
	}

	srcFile, err := getSrcFile(pkg)
	if err != nil {
		return err
	}

	gen := newFileGen(pkg, srcFile)

	err = inspectSrc(gen)
	if err != nil {
		return err
	}

	gen.updateDeclarations()

	err = gen.generateFile()
	if err != nil {
		return err
	}

	err = writeFormattedFile(rootDir, gen)
	if err != nil {
		return err
	}

	return nil
}

func joinPackageErrors(pkg *packages.Package) error {
	if len(pkg.Errors) > 0 {
		errs := make([]error, len(pkg.Errors))

		for i := range pkg.Errors {
			errs[i] = pkg.Errors[i]
		}

		return errors.Join(errs...)
	}

	return nil
}

func getSrcFile(pkg *packages.Package) (*ast.File, error) {
	if len(pkg.Syntax) < 1 {
		return nil, errors.New("srcFile not found")
	}

	return pkg.Syntax[0], nil
}

func findPackage(rootDir string) (*packages.Package, error) {
	fileset := token.NewFileSet()

	parsedPackages, err := packages.Load(
		&packages.Config{
			Mode:       packages.NeedImports | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedCompiledGoFiles | packages.NeedDeps | packages.NeedModule,
			Context:    context.TODO(),
			Dir:        rootDir,
			BuildFlags: []string{"-tags=vanya"},
			Fset:       fileset,
			Tests:      true,
		}, ConfigSrcFileName,
	)
	if err != nil {
		return nil, err
	}

	if len(parsedPackages) == 0 {
		return nil, errors.New("unable to find package")
	}

	if len(parsedPackages) > 1 {
		return nil, errors.New("multiple packages found with config.go srcFile")
	}

	return parsedPackages[0], nil
}

func inspectSrc(gen *fileGen) error {
	errs := make([]error, 0)

	ast.Inspect(
		gen.srcFile, func(node ast.Node) bool {
			var err error

			switch node.(type) {
			case *ast.FuncDecl:
				err = inspectMain(gen, node.(*ast.FuncDecl))
			case *ast.ImportSpec:
				err = inspectImport(gen, node.(*ast.ImportSpec))
			}

			if err != nil {
				errs = append(errs, err)
			}

			return true
		},
	)

	err := errors.Join(errs...)
	if err != nil {
		return err
	}

	for _, arg := range gen.buildArgs {
		inspectTypes(gen, arg)
	}

	return nil
}

func inspectImport(gen *fileGen, node *ast.ImportSpec) error {
	pkg, ok := gen.pkg.Imports[strings.Trim(node.Path.Value, `"`)]
	if !ok {
		return fmt.Errorf("imported package %s not found", node.Path.Value)
	}

	pkgName := pkg.Types.Name()

	if node.Name != nil {
		pkgName = node.Name.Name
	}

	gen.imports[pkgName] = pkg

	return nil
}

func inspectMain(gen *fileGen, node *ast.FuncDecl) error {
	if node.Name.Name != "main" {
		return nil
	}

	if node.Body == nil || node.Body.List == nil {
		return errors.New("invalid main function body")
	}

	for _, statement := range node.Body.List {
		exprStmt, ok := statement.(*ast.ExprStmt)
		if !ok {
			continue
		}

		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		callIdent, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			continue
		}

		if callIdent.Name != "vanya" && selectorExpr.Sel.Name != "BuildConfigs" {
			continue
		}

		gen.buildArgs = callExpr.Args
	}

	return nil
}

func inspectTypes(gen *fileGen, arg ast.Expr) {
	switch arg.(type) {
	case *ast.UnaryExpr:
		inspectTypes(gen, arg.(*ast.UnaryExpr).X)
	case *ast.CompositeLit:
		compositeLit, ok := arg.(*ast.CompositeLit)
		if !ok {
			log.Println("only composite literals are expected as argument in BuildConfig")
			return
		}

		switch compositeLit.Type.(type) {
		case *ast.SelectorExpr:
			selectorExpr := compositeLit.Type.(*ast.SelectorExpr)
			xIdent, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				return
			}

			pkg, ok := gen.imports[xIdent.Name]
			if !ok {
				log.Printf("could not find package providing %s.%s\n", xIdent.Name, selectorExpr.Sel.Name)
				return
			}

			decl, ok := findDeclaration(pkg, selectorExpr.Sel.Name)
			if !ok {
				log.Printf("could not find declaration %s.%s in imported package", xIdent.Name, selectorExpr.Sel.Name)
				return
			}

			gen.objects = append(gen.objects, decl)
		case *ast.Ident:
			name := compositeLit.Type.(*ast.Ident).Name

			decl, ok := findDeclaration(gen.pkg, name)
			if !ok {
				log.Printf("could not find declaration for %s", name)
				return
			}

			gen.objects = append(gen.objects, decl)
		}

		return
	}
}

func findDeclaration(pkg *packages.Package, name string) (ast.Decl, bool) {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			if genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if typeSpec.Name.Name == name {
					return decl, true
				}
			}
		}
	}

	return nil, false
}

type alias = string

type fileGen struct {
	pkg       *packages.Package
	srcFile   *ast.File
	imports   map[alias]*packages.Package
	buildArgs []ast.Expr
	objects   []ast.Decl
	buf       *bytes.Buffer
}

func newFileGen(pkg *packages.Package, file *ast.File) *fileGen {
	return &fileGen{
		pkg:       pkg,
		srcFile:   file,
		imports:   make(map[alias]*packages.Package),
		buildArgs: make([]ast.Expr, 0),
		objects:   make([]ast.Decl, 0),
		buf:       &bytes.Buffer{},
	}
}

func (f *fileGen) updateDeclarations() {
	for _, decl := range f.objects {
		typeSpec := decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
		structType := typeSpec.Type.(*ast.StructType)
		for _, field := range structType.Fields.List {
			field.Tag = &ast.BasicLit{
				Value: fmt.Sprintf("`mapstructure:\"%s\"`", toSnakeCase(field.Names[0].Name)),
			}
		}
	}
}

func (f *fileGen) generateFile() error {
	err := f.generateFrame()
	if err != nil {
		return err
	}

	err = f.generateConfig()
	if err != nil {
		return err
	}

	err = f.generateConfigConstructor()
	if err != nil {
		return err
	}

	err = f.generateDefaultConstructor()
	if err != nil {
		return err
	}

	return nil
}

const frameFormat = `// Code generated by Vanya: DO NOT EDIT.
// versions:
// 	vanya v0.0.0
// source: github.com/ivanmashin/my-service/configs/config.go

//go:generate go run github.com/ivanmashin/vanya/cmd/config
//go:build !vanya
// +build !vanya

package %s

import "github.com/ivanmashin/vanya/pkg/configs"
`

func (f *fileGen) generateFrame() error {
	_, err := f.buf.WriteString(fmt.Sprintf(frameFormat, f.pkg.Types.Name()))
	if err != nil {
		return err
	}

	return nil
}

func (f *fileGen) generateConfig() error {
	fields := 1 + len(f.objects)*2 // = 1 embeddingField + (1 field + 1 space)*len(objects)
	fieldList := make([]*ast.Field, fields)
	fieldList[0] = &ast.Field{
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("configs"),
			Sel: ast.NewIdent("Embedding"),
		},
	}

	cfgObj := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Config"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fieldList,
					},
				},
			},
		},
	}

	for i, obj := range f.objects {
		genDecl := obj.(*ast.GenDecl)
		spec := genDecl.Specs[0].(*ast.TypeSpec)

		fieldList[i*2+1] = &ast.Field{
			Type: ast.NewIdent(""),
		}

		fieldList[i*2+2] = &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(spec.Name.Name),
			},
			Type: spec.Type,
			Tag: &ast.BasicLit{
				Value: fmt.Sprintf("`mapstructure:\"%s\"`", toSnakeCase(spec.Name.Name)),
			},
		}
	}

	err := printer.Fprint(f.buf, f.pkg.Fset, cfgObj)
	if err != nil {
		return err
	}

	return nil
}

const constructor = `
func NewConfig(opts ...configs.Option) (Config, error) {
	c := NewDefaultConfig()

	err := c.Init(&c, opts...)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
`

func (f *fileGen) generateConfigConstructor() error {
	_, err := f.buf.WriteString("\n" + constructor + "\n")
	if err != nil {
		return err
	}

	return nil
}

const defaultConfigTemplate = `
func NewDefaultConfig() Config {
	return Config{
		Embedding: configs.Embedding{},
		{{ range . }}{{ .Key }}: {{ .Type }}{
			{{ range .Defaults }}{{ . }},
			{{ end }}
		},
		{{ end }}
	}
}
`

type templateData struct {
	Key      string
	Type     string
	Defaults []string
}

func (f *fileGen) generateDefaultConstructor() error {
	data := make([]templateData, 0)

	for i, obj := range f.objects {
		genDecl := obj.(*ast.GenDecl)
		spec := genDecl.Specs[0].(*ast.TypeSpec)

		structType, ok := spec.Type.(*ast.StructType)
		if ok {
			removeAllCommentsFromStruct(structType)
		}

		mainFuncArg := f.buildArgs[i].(*ast.CompositeLit)

		bt := &bytes.Buffer{}
		err := printer.Fprint(bt, f.pkg.Fset, spec.Type)
		if err != nil {
			return err
		}

		defaults := make([]string, 0)
		b := &bytes.Buffer{}
		for _, expr := range mainFuncArg.Elts {
			err := printer.Fprint(b, f.pkg.Fset, expr)
			if err != nil {
				return err
			}

			defaults = append(defaults, b.String())
			b.Reset()
		}

		data = append(
			data, templateData{
				Key: spec.Name.Name, Type: bt.String(), Defaults: defaults,
			},
		)
	}

	err := template.Must(template.New("value").Parse(defaultConfigTemplate)).Execute(
		f.buf, data,
	)
	if err != nil {
		return err
	}

	return nil
}

func removeAllCommentsFromStruct(node *ast.StructType) {
	for _, field := range node.Fields.List {
		field.Doc = nil
		field.Comment = nil
	}
}

func writeFormattedFile(rootDir string, gen *fileGen) error {
	f, err := os.OpenFile(filepath.Join(rootDir, ConfigDstFileName), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	src, err := format.Source(gen.buf.Bytes())
	if err != nil {
		return err
	}

	_, err = io.Copy(f, bytes.NewReader(src))
	if err != nil {
		return err
	}

	return nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(s string) string {
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
