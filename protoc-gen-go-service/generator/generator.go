package generator

import (
	"bytes"
	"html/template"

	"github.com/easeq/go-service/protoc-gen-go-service/options"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

// Generator holds the config used to generate the files
type Generator struct {
	FilenamePrefix string
	GoPackageName  protogen.GoPackageName
	GoImportPath   protogen.GoImportPath
	Gen            *protogen.Plugin
	Streams        map[string]int
	RegistryTags   map[string]*options.RegistryTag
	Services       []*protogen.Service
	Imports        map[string]bool
}

// GenerateFile generates a _ascii.pb.go file containing gRPC service definitions.
func (g *Generator) GenerateFile() *protogen.GeneratedFile {
	filename := g.FilenamePrefix + ".pb.gs.go"
	gf := g.Gen.NewGeneratedFile(filename, g.GoImportPath)

	t, err := template.New("gs.tmpl").Funcs(template.FuncMap{
		"camelCase": strcase.ToLowerCamel,
		"add": func(x, y int) int {
			x = x + y
			return x
		},
	}).Parse(tmpl)
	if err != nil {
		panic(err)
	}

	var result bytes.Buffer
	if err := t.Execute(&result, g); err != nil {
		panic(err)
	}

	gf.P(result.String())

	return gf
}
