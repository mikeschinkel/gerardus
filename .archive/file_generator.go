package _archive

//import (
//	"fmt"
//	"strings"
//
//	"gerardus/.archive"
//	"gerardus/collector"
//	"gerardus/parser"
//	"gerardus/surveyor"
//)
//
//type FileGenerator struct {
//	file   parser.File
//	survey *surveyor.CodeSurveyor
//}
//
//func NewFileGenerator(file parser.File, survey *surveyor.CodeSurveyor) *FileGenerator {
//	return &FileGenerator{
//		file:   file,
//		survey: survey,
//	}
//}
//
//func (g *FileGenerator) GeneratePackage() string {
//	return fmt.Sprintf("package %s\n\n",
//		parser.AsGoFile(g.file).LocalPackageName(),
//	)
//}
//
//func (g *FileGenerator) ShouldGenerate() (generate bool) {
//	switch t := g.file.(type) {
//	case *parser.ModFile:
//		generate = true
//	case *parser.GoFile:
//		generate = true
//		if t.HasImports() {
//			goto end
//		}
//		if t.HasTypes() {
//			goto end
//		}
//		generate = false
//	}
//end:
//	return generate
//}
//
//func (g *FileGenerator) GenerateContent() (content string) {
//	switch t := g.file.(type) {
//	case *parser.ModFile:
//		content = string(t.Content)
//	case *parser.GoFile:
//		sb := strings.Builder{}
//		sb.WriteString(g.GeneratePackage())
//		sb.WriteString(g.GenerateImports())
//		sb.WriteString(g.GenerateInterfaces())
//		//sb.WriteString(g.GenerateTypestubs())
//		content = sb.String()
//	}
//	return content
//}
//
//func (g *FileGenerator) GenerateImports() string {
//	var sb strings.Builder
//
//	f, ok := g.file.(*parser.GoFile)
//	if !ok {
//		goto end
//	}
//
//	if len(f.Imports) == 0 {
//		goto end
//	}
//	sb = strings.Builder{}
//	sb.WriteString("import (\n")
//	for _, imp := range f.Imports.Sorted() {
//		sb.WriteString(g.GenerateImport(imp))
//	}
//	sb.WriteString(")\n")
//	sb.WriteByte('\n')
//end:
//	return sb.String()
//}
//
//func (g *FileGenerator) GenerateImport(imp *parser.Import) (s string) {
//	if imp.Alias == "" {
//		s = fmt.Sprintf(`"%s"`, imp.Package.Name)
//	} else {
//		s = fmt.Sprintf(`%s "%s"`, imp.Alias, imp.Package.Name)
//	}
//	return fmt.Sprintf("\t%s\n", s)
//}
//
//func (g *FileGenerator) GenerateInterfaces() string {
//	var sb strings.Builder
//
//	f, ok := g.file.(*parser.GoFile)
//	if !ok {
//		goto end
//	}
//
//	sb = strings.Builder{}
//	for _, iface := range f.Types.Sorted() {
//		if iface.symType != collector.InterfaceSymbol {
//			continue
//		}
//		if !iface.HasMethods() {
//			continue
//		}
//		sb.WriteString(g.GenerateInterface(iface))
//	}
//	sb.WriteByte('\n')
//end:
//	return sb.String()
//}
//
//func (g *FileGenerator) GenerateInterface(iface *_archive.Type) string {
//	var sb strings.Builder
//	sb.WriteString(fmt.Sprintf("type %s interface {\n", iface.Name()))
//	for _, m := range iface.Methods {
//		sb.WriteString(fmt.Sprintf("\t%s\n", g.GenerateInterfaceMethod(m)))
//	}
//	sb.WriteString("}\n\n")
//	return sb.String()
//}
//
////func (g *FileGenerator) GenerateTypestubs() string {
////	var sb strings.Builder
////
////	f, ok := g.file.(*GoFile)
////	if !ok {
////		goto end
////	}
////
////	if len(g.survey.Types) == 0 {
////		goto end
////	}
////	sb = strings.Builder{}
////	for _, t := range f.Types.Sorted() {
////		sb.WriteString(g.GenerateTypeStub(t))
////	}
////	sb.WriteByte('\n')
////end:
////	return sb.String()
////}
//
//// GenerateTypeStub generates only enough of a type as to satisfy interfaces
//func (g *FileGenerator) GenerateTypeStub(t *_archive.Type) string {
//	return fmt.Sprintf("type %s struct {}\n", t.Name())
//}
//
//func (g *FileGenerator) GenerateInterfaceMethod(m *parser.Method) string {
//	var sb strings.Builder
//
//	sb.WriteString(m.Name)
//	if !m.IsEmbed() {
//		sb.WriteString("(")
//	}
//	i := 0
//	for _, param := range m.params {
//		if i > 0 {
//			sb.WriteString(", ")
//		}
//		if len(param.Name) > 0 {
//			sb.WriteString(param.Name)
//			sb.WriteByte(' ')
//		}
//		sb.WriteString(param.Type)
//		i++
//	}
//	if !m.IsEmbed() {
//		sb.WriteString(") ")
//	} else {
//		sb.WriteByte(' ')
//	}
//
//	if len(m.results) > 1 {
//		sb.WriteString("(")
//	}
//	i = 0
//	for _, result := range m.results {
//		if i > 0 {
//			sb.WriteString(", ")
//		}
//		if len(result.Name) > 0 {
//			sb.WriteString(result.Name)
//			sb.WriteByte(' ')
//		}
//		sb.WriteString(result.Type)
//		i++
//	}
//	if len(m.results) > 1 {
//		sb.WriteString(")")
//	}
//	return sb.String()
//}
