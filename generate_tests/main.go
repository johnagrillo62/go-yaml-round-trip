package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

type StructInfo struct {
	Name   string
	Fields []string
}

func main() {
        fmt.Println("Found structs:")
	fmt.Println("Starting generator")
        fmt.Println("Args:", os.Args)

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input.go> <output_test.go>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	srcBytes, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	src := string(srcBytes)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var structs []StructInfo
	fmt.Println("Found structs:")


	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			s := StructInfo{Name: typeSpec.Name.Name}
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				s.Fields = append(s.Fields, field.Names[0].Name)
			}
			structs = append(structs, s)
		}
	}


const tmplText = `package main

import (
	"os"
	"reflect"
	"testing"
	"gopkg.in/yaml.v3"
)

func validateYAMLMatchesStruct(t *testing.T, yamlFile string, structType interface{}) {
	data, err := os.ReadFile(yamlFile)
	if err != nil { t.Fatal(err) }

	var yamlMap map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlMap); err != nil { t.Fatal(err) }

	val := reflect.TypeOf(structType)
	if val.Kind() == reflect.Ptr { val = val.Elem() }

	tags := map[string]bool{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" { tag = field.Name }
		tags[tag] = true
	}

	for key := range yamlMap {
		if !tags[key] {
			t.Errorf("YAML key '%s' does not exist in struct %s", key, val.Name())
		}
	}

	for tag := range tags {
		if _, ok := yamlMap[tag]; !ok {
			t.Errorf("Struct field '%s' is missing in YAML", tag)
		}
	}
}

{{range .}}
func TestYAML_{{.Name}}_Valid(t *testing.T) {
	validateYAMLMatchesStruct(t, "valid_{{.Name | ToLower}}.yaml", {{.Name}}{})
}

func TestYAML_{{.Name}}_Invalid(t *testing.T) {
	validateYAMLMatchesStruct(t, "invalid_{{.Name | ToLower}}.yaml", {{.Name}}{})
}
{{end}}
`

	tmpl := template.Must(template.New("yamlTests").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).Parse(tmplText))

	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, structs); err != nil {
		panic(err)
	}

	fmt.Println("Generated YAML validation tests:", outputFile)
}
