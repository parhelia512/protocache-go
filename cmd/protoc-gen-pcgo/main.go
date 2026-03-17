package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/peterrk/slices"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func debugPrintln(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}

func debugPrintf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

const NoneKind = protoreflect.Kind(0)

type Alias struct {
	Key       protoreflect.Kind
	Value     protoreflect.Kind
	ValueType string
}

type Type struct {
	Package string
	GoName  string
}

var (
	typeBook  = make(map[string]Type)
	aliasBook = make(map[string]Alias)

	builinTypes = map[string]string{
		"Bytes":        "[]byte",
		"String":       "string",
		"Float64":      "float64",
		"Float32":      "float32",
		"Uint64":       "uint64",
		"Int64":        "int64",
		"Uint32":       "uint32",
		"Int32":        "int32",
		"Bool":         "bool",
		"Float64Array": "[]float64",
		"Float32Array": "[]float32",
		"Uint64Array":  "[]uint64",
		"Int64Array":   "[]int64",
		"Uint32Array":  "[]uint32",
		"Int32Array":   "[]int32",
		"BoolArray":    "[]bool",
	}
)

func (a *Alias) CalcName(imports map[string]string) (string, error) {
	if a.Key == NoneKind {
		switch a.Value {
		case protoreflect.MessageKind:
			if t, got := typeBook[a.ValueType]; !got {
				return "", fmt.Errorf("unknown type: %s", a.ValueType)
			} else if pkg, got := imports[t.Package]; got {
				return "ARRAY_" + pkg + t.GoName, nil
			} else {
				return "ARRAY_" + t.GoName, nil
			}
		case protoreflect.BytesKind:
			return "protocache.BytesArray", nil
		case protoreflect.StringKind:
			return "protocache.StringArray", nil
		case protoreflect.DoubleKind:
			return "protocache.Float64Array", nil
		case protoreflect.FloatKind:
			return "protocache.Float32Array", nil
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return "protocache.Uint64Array", nil
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return "protocache.Int64Array", nil
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return "protocache.Uint32Array", nil
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return "protocache.Int32Array", nil
		case protoreflect.BoolKind:
			return "protocache.BoolArray", nil
		case protoreflect.EnumKind:
			if t, got := typeBook[a.ValueType]; !got {
				return "", fmt.Errorf("unknown type: %s", a.ValueType)
			} else if pkg, got := imports[t.Package]; got {
				return fmt.Sprintf("protocache.EnumArray[%s.%s]", pkg, t.GoName), nil
			} else {
				return fmt.Sprintf("protocache.EnumArray[%s]", t.GoName), nil
			}
		default:
			return "", fmt.Errorf("unsupported type: %v", a.Value)
		}
	}

	prefix := "MAP_"
	switch a.Key {
	case protoreflect.StringKind:
		prefix += "string_"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		prefix += "uint64_"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		prefix += "int64_"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		prefix += "uint32_"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		prefix += "int32_"
	default:
		return "", fmt.Errorf("unsupported key type: %v", a.Key)
	}

	switch a.Value {
	case protoreflect.MessageKind,
		protoreflect.EnumKind:
		if t, got := typeBook[a.ValueType]; !got {
			return "", fmt.Errorf("unknown type: %s", a.ValueType)
		} else if pkg, got := imports[t.Package]; got {
			return prefix + pkg + t.GoName, nil
		} else {
			return prefix + t.GoName, nil
		}
	case protoreflect.BytesKind:
		return prefix + "bytes", nil
	case protoreflect.StringKind:
		return prefix + "string", nil
	case protoreflect.DoubleKind:
		return prefix + "float64", nil
	case protoreflect.FloatKind:
		return prefix + "float32", nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return prefix + "uint64", nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return prefix + "int64", nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return prefix + "uint32", nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return prefix + "int32", nil
	case protoreflect.BoolKind:
		return prefix + "bool", nil
	default:
		return "", fmt.Errorf("unsupported value type: %v", a.Value)
	}
}

func CollectEnums(pkg string, list []*protogen.Enum) {
	for _, one := range list {
		fullname := string(one.Desc.FullName())
		typeBook[fullname] = Type{
			Package: pkg,
			GoName:  one.GoIdent.GoName,
		}
	}
}

func CollectMessages(pkg string, list []*protogen.Message) {
	for _, one := range list {
		if one.Desc.IsMapEntry() {
			continue
		}
		fullname := string(one.Desc.FullName())
		CollectEnums(pkg, one.Enums)
		CollectMessages(pkg, one.Messages)
		typeBook[fullname] = Type{
			Package: pkg,
			GoName:  one.GoIdent.GoName,
		}
		if len(one.Fields) != 1 || one.Fields[0].Desc.Name() != "_" {
			continue
		}

		var alias Alias
		field := one.Fields[0].Desc
		if field.IsMap() {
			alias.Key = field.MapKey().Kind()
			field = field.MapValue()
		} else if field.IsList() {
			alias.Key = NoneKind
		} else {
			continue
		}
		alias.Value = field.Kind()
		if alias.Value == protoreflect.EnumKind {
			alias.ValueType = string(field.Enum().FullName())
		} else if alias.Value == protoreflect.MessageKind {
			alias.ValueType = string(field.Message().FullName())
		}
		aliasBook[fullname] = alias
	}
}

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		opt := parseOptions(gen.Request.GetParameter())
		for _, file := range gen.Files {
			CollectEnums(string(file.GoImportPath), file.Enums)
			CollectMessages(string(file.GoImportPath), file.Messages)
		}

		//debugPrintln(typeBook)
		//debugPrintln(aliasBook)

		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}
			err := GenFile(gen, file)
			if err != nil {
				return err
			}
			if opt.EX {
				err = GenEXFile(gen, file)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

type Options struct {
	EX bool
}

func parseOptions(raw string) Options {
	var out Options
	for _, one := range strings.Split(raw, ",") {
		one = strings.TrimSpace(one)
		switch one {
		case "", ".":
		case "extra", "extra=true":
			out.EX = true
		}
	}
	return out
}

func CollectImports(pkg string, list []*protogen.Message, book map[string]string) {
	for _, one := range list {
		if one.Desc.IsMapEntry() {
			continue
		}
		CollectImports(pkg, one.Messages, book)
		for _, field := range one.Fields {
			desc := field.Desc
			if desc.IsMap() {
				desc = desc.MapValue()
			}
			var valueType string
			if desc.Kind() == protoreflect.EnumKind {
				valueType = string(desc.Enum().FullName())
			} else if desc.Kind() == protoreflect.MessageKind {
				valueType = string(desc.Message().FullName())
			} else {
				continue
			}
			t, got := typeBook[valueType]
			if got && t.Package != pkg {
				book[t.Package] = ""
			}
		}
	}
}

func GenEnums(g *protogen.GeneratedFile, list []*protogen.Enum) error {
	for _, one := range list {
		g.P()
		g.P("type ", one.GoIdent.GoName, " protocache.EnumValue")
		g.P()
		g.P("const (")
		for _, v := range one.Values {
			g.P("	", v.GoIdent.GoName, " ", one.GoIdent.GoName, " = ", v.Desc.Number())
		}
		g.P(")")
	}
	return nil
}

func GenMessages(g *protogen.GeneratedFile, imports map[string]string,
	derived map[string]Alias, list []*protogen.Message) error {
	getPkgPrefix := func(fullname string) string {
		if t, got := typeBook[fullname]; got {
			if pkg, got := imports[t.Package]; got {
				return pkg + "."
			}
		}
		return ""
	}

	for _, one := range list {
		if one.Desc.IsMapEntry() {
			continue
		}
		GenEnums(g, one.Enums)
		if err := GenMessages(g, imports, derived, one.Messages); err != nil {
			return err
		}

		alias, ok := aliasBook[string(one.Desc.FullName())]
		if ok {
			name, err := alias.CalcName(imports)
			if err != nil {
				return err
			}
			derived[name] = alias
			g.P()
			g.P("type ", one.GoIdent.GoName, " = ", name)
			g.P()
			g.P("func AS_", one.GoIdent.GoName, "(data []byte)", one.GoIdent.GoName, " {")
			if strings.HasPrefix(name, "protocache.") {
				g.P("	return protocache.As", name[11:], "(data)")
			} else {
				g.P("	return AS_", name, "(data)")
			}
			g.P("}")
			continue
		}

		fields := make([]*protogen.Field, len(one.Fields))
		copy(fields, one.Fields)
		order := slices.Order[*protogen.Field]{
			Less: func(a, b *protogen.Field) bool {
				return a.Desc.Number() < b.Desc.Number()
			},
		}
		order.Sort(fields)
		maxID := protoreflect.FieldNumber(0)
		for _, field := range fields {
			if maxID < field.Desc.Number() {
				maxID = field.Desc.Number()
			}
		}

		g.P()
		g.P("const (")
		for _, field := range fields {
			g.P("	_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), " uint16 = ", field.Desc.Number()-1)
		}
		g.P("	_FIELD_TOTAL_", one.GoIdent.GoName, " uint16 = ", maxID)
		g.P(")")
		g.P()
		g.P("type ", one.GoIdent.GoName, " struct { core protocache.Message }")
		g.P()
		g.P("func AS_", one.GoIdent.GoName, "(data []byte) ", one.GoIdent.GoName,
			"{ return ", one.GoIdent.GoName, "{core: protocache.AsMessage(data)}}")
		g.P()
		g.P("func (m *", one.GoIdent.GoName, ") IsValid() bool {return m.core.IsValid()}")

		handleSimpleField := func(field *protogen.Field, t string) {
			g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() ", builinTypes[t], " {")
			g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
			g.P("	return field.Get", t, "()")
			g.P("}")
		}

		handleStingrArray := func(field *protogen.Field, t string) {
			g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() protocache.", t, "Array {")
			g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
			g.P("	return protocache.As", t, "Array(field.GetObject())")
			g.P("}")
		}

		handleComlexField := func(field *protogen.Field, t string) {
			g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() ", t, " {")
			g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
			g.P("	return AS_", t, "(field.GetObject())")
			g.P("}")
		}

		for _, field := range fields {
			desc := field.Desc
			g.P()
			if desc.IsMap() {
				vDesc := desc.MapValue()
				alias := Alias{
					Key:   desc.MapKey().Kind(),
					Value: vDesc.Kind(),
				}
				if alias.Value == protoreflect.EnumKind {
					alias.ValueType = string(vDesc.Enum().FullName())
				} else if alias.Value == protoreflect.MessageKind {
					alias.ValueType = string(vDesc.Message().FullName())
				}
				name, err := alias.CalcName(imports)
				if err != nil {
					return err
				}
				derived[name] = alias
				handleComlexField(field, name)
			} else if desc.IsList() {
				switch desc.Kind() {
				case protoreflect.MessageKind:
					alias := Alias{
						Key:       NoneKind,
						Value:     desc.Kind(),
						ValueType: string(desc.Message().FullName()),
					}
					name, err := alias.CalcName(imports)
					if err != nil {
						return err
					}
					derived[name] = alias
					handleComlexField(field, name)
				case protoreflect.BytesKind:
					handleStingrArray(field, "Bytes")
				case protoreflect.StringKind:
					handleStingrArray(field, "String")
				case protoreflect.DoubleKind:
					handleSimpleField(field, "Float64Array")
				case protoreflect.FloatKind:
					handleSimpleField(field, "Float32Array")
				case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					handleSimpleField(field, "Uint64Array")
				case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
					handleSimpleField(field, "Int64Array")
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
					handleSimpleField(field, "Uint32Array")
				case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
					handleSimpleField(field, "Int32Array")
				case protoreflect.BoolKind:
					handleSimpleField(field, "BoolArray")
				case protoreflect.EnumKind:
					pkgPrefix := getPkgPrefix(string(field.Enum.Desc.FullName()))
					typeName := pkgPrefix + field.Enum.GoIdent.GoName
					g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() []", typeName, " {")
					g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
					g.P("	return protocache.CastEnumArray[", typeName, "](field.GetEnumValueArray())")
					g.P("}")
				default:
					return fmt.Errorf("unsupported type: %v", desc.Kind())
				}
			} else {
				switch desc.Kind() {
				case protoreflect.MessageKind:
					pkgPrefix := getPkgPrefix(string(field.Message.Desc.FullName()))
					typeName := field.Message.GoIdent.GoName
					g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() ", pkgPrefix, typeName, " {")
					g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
					g.P("	return ", pkgPrefix, "AS_", typeName, "(field.GetObject())")
					g.P("}")

					if _, hit := aliasBook[string(field.Message.Desc.FullName())]; !hit {
						g.P("func (m *", one.GoIdent.GoName, ") Has", field.GoName, "() bool {")
						g.P("	return m.core.HasField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
						g.P("}")
					}
				case protoreflect.BytesKind:
					handleSimpleField(field, "Bytes")
				case protoreflect.StringKind:
					handleSimpleField(field, "String")
				case protoreflect.DoubleKind:
					handleSimpleField(field, "Float64")
				case protoreflect.FloatKind:
					handleSimpleField(field, "Float32")
				case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					handleSimpleField(field, "Uint64")
				case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
					handleSimpleField(field, "Int64")
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
					handleSimpleField(field, "Uint32")
				case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
					handleSimpleField(field, "Int32")
				case protoreflect.BoolKind:
					handleSimpleField(field, "Bool")
				case protoreflect.EnumKind:
					pkgPrefix := getPkgPrefix(string(field.Enum.Desc.FullName()))
					typeName := pkgPrefix + field.Enum.GoIdent.GoName
					g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() ", typeName, " {")
					g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
					g.P("	return ", typeName, "(field.GetEnumValue())")
					g.P("}")
				default:
					return fmt.Errorf("unsupported type: %v", desc.Kind())
				}
			}
		}
	}
	return nil
}

func GenFile(gen *protogen.Plugin, file *protogen.File) error {
	imports := make(map[string]string)
	CollectImports(string(file.GoImportPath), file.Messages, imports)

	pkgs := make([]string, 0, len(imports))
	for one, _ := range imports {
		pkgs = append(pkgs, one)
	}
	slices.Sort(pkgs)
	for i, one := range pkgs {
		imports[one] = fmt.Sprintf("p%d", i+1)
	}

	filename := file.GeneratedFilenamePrefix + ".pc.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("import (")
	g.P(`	"github.com/peterrk/protocache-go"`)
	for name, mark := range imports {
		g.P("	", mark, ` "`, name, `"`)
	}
	g.P(")")

	if err := GenEnums(g, file.Enums); err != nil {
		return err
	}

	derived := make(map[string]Alias)
	if err := GenMessages(g, imports, derived, file.Messages); err != nil {
		return err
	}

	names := make([]string, 0, len(derived))
	for one, _ := range derived {
		if !strings.HasPrefix(one, "protocache.") {
			names = append(names, one)
		}
	}
	slices.Sort(names)
	//debugPrintln(names)

	handleSimpleKV := func(name, kt, vt string) {
		g.P("func (x *", name, ") Value(i uint32) ", builinTypes[vt], " {")
		g.P("	field := x.core.Value(i)")
		g.P("	return field.Get", vt, "()")
		g.P("}")
		g.P()
		g.P("func (x *", name, ") Find(key ", builinTypes[kt], ") (", builinTypes[vt], ", bool) {")
		g.P("	field := x.core.FindBy", kt, "(key)")
		g.P("	return field.Get", vt, "(), field.IsValid()")
		g.P("}")
	}

	getValueType := func(fullname string) (prefix string, name string) {
		t := typeBook[fullname]
		if pkg, got := imports[t.Package]; got {
			return pkg + ".", t.GoName
		} else {
			return "", t.GoName
		}
	}
	for _, name := range names {
		g.P()
		alias := derived[name]
		if alias.Key == NoneKind {
			if alias.Value != protoreflect.MessageKind {
				panic("message array only")
			}
			g.P("type ", name, " struct { core protocache.Array }")
			g.P()
			g.P("func AS_", name, "(data []byte) ", name, " {")
			g.P("	return ", name, "{core: protocache.AsArray(data)}")
			g.P("}")
			g.P()
			prefix, value := getValueType(alias.ValueType)
			g.P("func (x *", name, ") Get(i uint32) ", prefix, value, " {")
			g.P("	field := x.core.Get(i)")
			g.P("	return ", prefix, "AS_", value, "(field.GetObject())")
			g.P("}")
		} else {
			g.P("type ", name, " struct { core protocache.Map }")
			g.P()
			g.P("func AS_", name, "(data []byte) ", name, " {")
			g.P("	return ", name, "{core: protocache.AsMap(data)}")
			g.P("}")

			var key string
			switch alias.Key {
			case protoreflect.StringKind:
				key = "String"
			case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				key = "Uint64"
			case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				key = "Int64"
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
				key = "Uint32"
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
				key = "Int32"
			default:
				panic("unsupported key type")
			}

			g.P()
			g.P("func (x *", name, ") Key(i uint32) ", builinTypes[key], " {")
			g.P("	field := x.core.Key(i)")
			g.P("	return field.Get", key, "()")
			g.P("}")
			g.P()

			switch alias.Value {
			case protoreflect.MessageKind:
				prefix, value := getValueType(alias.ValueType)
				g.P("func (x *", name, ") Value(i uint32) ", prefix, value, " {")
				g.P("	field := x.core.Value(i)")
				g.P("	return ", prefix, "AS_", value, "(field.GetObject())")
				g.P("}")
				g.P()
				g.P("func (x *", name, ") Find(key ", builinTypes[key], ") (", prefix, value, ", bool) {")
				g.P("	field := x.core.FindBy", key, "(key)")
				g.P("	return ", prefix, "AS_", value, "(field.GetObject()), field.IsValid()")
				g.P("}")
			case protoreflect.BytesKind:
				handleSimpleKV(name, key, "Bytes")
			case protoreflect.StringKind:
				handleSimpleKV(name, key, "String")
			case protoreflect.DoubleKind:
				handleSimpleKV(name, key, "Float64")
			case protoreflect.FloatKind:
				handleSimpleKV(name, key, "Float32")
			case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				handleSimpleKV(name, key, "Uint64")
			case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				handleSimpleKV(name, key, "Int64")
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
				handleSimpleKV(name, key, "Uint32")
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
				handleSimpleKV(name, key, "Int32")
			case protoreflect.BoolKind:
				handleSimpleKV(name, key, "Bool")
			case protoreflect.EnumKind:
				prefix, value := getValueType(alias.ValueType)
				g.P("func (x *", name, ") Value(i uint32) ", prefix, value, " {")
				g.P("	field := x.core.Value(i)")
				g.P("	return ", prefix, value, "(field.GetEnumValue())")
				g.P("}")
				g.P()
				g.P("func (x *", name, ") Find(key ", builinTypes[key], ") (", prefix, value, ", bool) {")
				g.P("	field := x.core.FindBy", key, "(key)")
				g.P("	return ", prefix, value, "(field.GetEnumValue()), field.IsValid()")
				g.P("}")
			default:
				panic("unsupported value type")
			}
		}
		g.P()
		g.P("func (x *", name, ") IsValid() bool {return x.core.IsValid()}")
		g.P()
		g.P("func (x *", name, ") Size() uint32 {return x.core.Size()}")
	}
	return nil
}

func exFieldName(field *protogen.Field) string {
	return "f" + field.GoName
}

func exSupportsField(desc protoreflect.FieldDescriptor) bool {
	return true
}

func exIsAliasMessage(desc protoreflect.MessageDescriptor) bool {
	_, ok := aliasBook[string(desc.FullName())]
	return ok
}

func exMapKeyType(desc protoreflect.FieldDescriptor) string {
	switch desc.Kind() {
	case protoreflect.StringKind:
		return "string"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	default:
		panic("unsupported map key type")
	}
}

func exNamedType(imports map[string]string, desc protoreflect.FieldDescriptor) string {
	switch desc.Kind() {
	case protoreflect.MessageKind:
		t := typeBook[string(desc.Message().FullName())]
		if pkg, got := imports[t.Package]; got {
			return pkg + "." + t.GoName
		}
		return t.GoName
	case protoreflect.EnumKind:
		t := typeBook[string(desc.Enum().FullName())]
		if pkg, got := imports[t.Package]; got {
			return pkg + "." + t.GoName
		}
		return t.GoName
	default:
		panic("descriptor has no named type")
	}
}

func exAliasType(imports map[string]string, desc protoreflect.MessageDescriptor) string {
	t := typeBook[string(desc.FullName())]
	if pkg, got := imports[t.Package]; got {
		return pkg + "." + t.GoName + "EX"
	}
	return t.GoName + "EX"
}

func exAliasCtor(imports map[string]string, desc protoreflect.MessageDescriptor) string {
	t := typeBook[string(desc.FullName())]
	name := "TO_" + t.GoName + "EX"
	if pkg, got := imports[t.Package]; got {
		return pkg + "." + name
	}
	return name
}

func exAliasSerializeExpr(access string) string {
	return "func() ([]uint32, error) { data, err := " + access + ".Serialize(); if err != nil { return nil, err }; return protocache.BytesToWords(data), nil }()"
}

func exGoType(imports map[string]string, field *protogen.Field) string {
	desc := field.Desc
	if desc.IsMap() {
		keyType := exMapKeyType(desc.MapKey())
		valueField := *field
		valueField.Desc = desc.MapValue()
		return "map[" + keyType + "]" + exGoType(imports, &valueField)
	}
	if desc.IsList() {
		switch desc.Kind() {
		case protoreflect.MessageKind:
			if exIsAliasMessage(desc.Message()) {
				return "[]" + exAliasType(imports, desc.Message())
			}
			return "[]*" + exNamedType(imports, desc) + "EX"
		case protoreflect.BytesKind:
			return "[][]byte"
		case protoreflect.StringKind:
			return "[]string"
		case protoreflect.DoubleKind:
			return "[]float64"
		case protoreflect.FloatKind:
			return "[]float32"
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return "[]uint64"
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return "[]int64"
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return "[]uint32"
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return "[]int32"
		case protoreflect.BoolKind:
			return "[]bool"
		case protoreflect.EnumKind:
			return "[]" + exNamedType(imports, desc)
		default:
			panic("unsupported ex field type")
		}
	}
	switch desc.Kind() {
	case protoreflect.MessageKind:
		if exIsAliasMessage(desc.Message()) {
			return exAliasType(imports, desc.Message())
		}
		return "*" + exNamedType(imports, desc) + "EX"
	case protoreflect.BytesKind:
		return "[]byte"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.EnumKind:
		return exNamedType(imports, desc)
	default:
		panic("unsupported ex field type")
	}
}

func exEncodeExpr(field *protogen.Field, access string) string {
	if field.Desc.IsList() {
		switch field.Desc.Kind() {
		case protoreflect.MessageKind:
			if exIsAliasMessage(field.Message.Desc) {
				return "protocache.EncodeObjectArray(len(" + access + "), func(i int) ([]uint32, error) { return " + exAliasSerializeExpr(access+"[i]") + " })"
			}
			return "protocache.EncodeObjectArray(len(" + access + "), func(i int) ([]uint32, error) { if " + access + "[i] == nil { return []uint32{0}, nil }; return " + access + "[i].serializeWords() })"
		case protoreflect.BytesKind:
			return "protocache.EncodeBytesArray(" + access + ")"
		case protoreflect.StringKind:
			return "protocache.EncodeStringArray(" + access + ")"
		case protoreflect.DoubleKind:
			return "protocache.EncodeFloat64Array(" + access + ")"
		case protoreflect.FloatKind:
			return "protocache.EncodeFloat32Array(" + access + ")"
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return "protocache.EncodeUint64Array(" + access + ")"
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return "protocache.EncodeInt64Array(" + access + ")"
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return "protocache.EncodeUint32Array(" + access + ")"
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return "protocache.EncodeInt32Array(" + access + ")"
		case protoreflect.BoolKind:
			return "protocache.EncodeBoolArray(" + access + ")"
		case protoreflect.EnumKind:
			return "protocache.EncodeEnumArray(" + access + ")"
		default:
			panic("unsupported ex encode type")
		}
	}
	switch field.Desc.Kind() {
	case protoreflect.MessageKind:
		if exIsAliasMessage(field.Message.Desc) {
			return exAliasSerializeExpr(access)
		}
		return access + ".serializeWords()"
	case protoreflect.BytesKind:
		return "protocache.EncodeBytes(" + access + ")"
	case protoreflect.StringKind:
		return "protocache.EncodeString(" + access + ")"
	case protoreflect.DoubleKind:
		return "protocache.EncodeFloat64(" + access + "), nil"
	case protoreflect.FloatKind:
		return "protocache.EncodeFloat32(" + access + "), nil"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "protocache.EncodeUint64(" + access + "), nil"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "protocache.EncodeInt64(" + access + "), nil"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "protocache.EncodeUint32(" + access + "), nil"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "protocache.EncodeInt32(" + access + "), nil"
	case protoreflect.BoolKind:
		return "protocache.EncodeBool(" + access + "), nil"
	case protoreflect.EnumKind:
		return "protocache.EncodeInt32(int32(" + access + ")), nil"
	default:
		panic("unsupported ex encode type")
	}
}

func exMapValueField(field *protogen.Field) *protogen.Field {
	valueField := *field
	valueField.Desc = field.Desc.MapValue()
	return &valueField
}

func exMapKeyEncodeExpr(desc protoreflect.FieldDescriptor, access string) string {
	switch desc.Kind() {
	case protoreflect.StringKind:
		return "protocache.EncodeString(" + access + ")"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "protocache.EncodeUint64(" + access + "), nil"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "protocache.EncodeInt64(" + access + "), nil"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "protocache.EncodeUint32(" + access + "), nil"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "protocache.EncodeInt32(" + access + "), nil"
	default:
		panic("unsupported map key type")
	}
}

func exGetExpr(field *protogen.Field, imports map[string]string, name string) []string {
	switch {
	case field.Desc.IsMap():
		valueField := exMapValueField(field)
		lines := []string{
			"pack := field.GetMap()",
			"m." + name + " = make(" + exGoType(imports, field) + ", int(pack.Size()))",
			"for i := uint32(0); i < pack.Size(); i++ {",
			"	keyField := pack.Key(i)",
		}
		switch field.Desc.MapKey().Kind() {
		case protoreflect.StringKind:
			lines = append(lines, "	key := keyField.GetString()")
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			lines = append(lines, "	key := keyField.GetUint64()")
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			lines = append(lines, "	key := keyField.GetInt64()")
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			lines = append(lines, "	key := keyField.GetUint32()")
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			lines = append(lines, "	key := keyField.GetInt32()")
		}
		lines = append(lines, "	valField := pack.Value(i)")
		switch valueField.Desc.Kind() {
		case protoreflect.MessageKind:
			call := exAliasCtor(imports, valueField.Desc.Message())
			if !exIsAliasMessage(valueField.Desc.Message()) {
				call = "TO_" + exNamedType(imports, valueField.Desc) + "EX"
			}
			lines = append(lines, "	m."+name+"[key] = "+call+"(valField.GetObject())")
		case protoreflect.BytesKind:
			lines = append(lines,
				"	if data := valField.GetBytes(); data != nil {",
				"		m."+name+"[key] = append([]byte(nil), data...)",
				"	} else {",
				"		m."+name+"[key] = nil",
				"	}",
			)
		case protoreflect.StringKind:
			lines = append(lines, "	m."+name+"[key] = valField.GetString()")
		case protoreflect.DoubleKind:
			lines = append(lines, "	m."+name+"[key] = valField.GetFloat64()")
		case protoreflect.FloatKind:
			lines = append(lines, "	m."+name+"[key] = valField.GetFloat32()")
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			lines = append(lines, "	m."+name+"[key] = valField.GetUint64()")
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			lines = append(lines, "	m."+name+"[key] = valField.GetInt64()")
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			lines = append(lines, "	m."+name+"[key] = valField.GetUint32()")
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			lines = append(lines, "	m."+name+"[key] = valField.GetInt32()")
		case protoreflect.BoolKind:
			lines = append(lines, "	m."+name+"[key] = valField.GetBool()")
		case protoreflect.EnumKind:
			lines = append(lines, "	m."+name+"[key] = "+exGoType(imports, valueField)+"(valField.GetEnumValue())")
		default:
			panic("unsupported map value type")
		}
		lines = append(lines, "}")
		return lines
	case field.Desc.IsList():
		switch field.Desc.Kind() {
		case protoreflect.MessageKind:
			call := exAliasCtor(imports, field.Message.Desc)
			if !exIsAliasMessage(field.Message.Desc) {
				call = "TO_" + field.Message.GoIdent.GoName + "EX"
				if pkg, got := imports[string(field.Message.GoIdent.GoImportPath)]; got {
					call = pkg + "." + call
				}
			}
			return []string{
				"arr := field.GetArray()",
				"m." + name + " = make(" + exGoType(imports, field) + ", int(arr.Size()))",
				"for i := uint32(0); i < arr.Size(); i++ {",
				"	elem := arr.Get(i)",
				"	m." + name + "[i] = " + call + "(elem.GetObject())",
				"}",
			}
		case protoreflect.BytesKind:
			return []string{
				"arr := protocache.AsBytesArray(field.GetObject())",
				"m." + name + " = make([][]byte, int(arr.Size()))",
				"for i := uint32(0); i < arr.Size(); i++ {",
				"	if data := arr.Get(i); data != nil {",
				"		m." + name + "[i] = append([]byte(nil), data...)",
				"	}",
				"}",
			}
		case protoreflect.StringKind:
			return []string{
				"arr := protocache.AsStringArray(field.GetObject())",
				"m." + name + " = make([]string, int(arr.Size()))",
				"for i := uint32(0); i < arr.Size(); i++ {",
				"	m." + name + "[i] = arr.Get(i)",
				"}",
			}
		case protoreflect.DoubleKind:
			return []string{"m." + name + " = append([]float64(nil), field.GetFloat64Array()...)"}
		case protoreflect.FloatKind:
			return []string{"m." + name + " = append([]float32(nil), field.GetFloat32Array()...)"}
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return []string{"m." + name + " = append([]uint64(nil), field.GetUint64Array()...)"}
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return []string{"m." + name + " = append([]int64(nil), field.GetInt64Array()...)"}
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return []string{"m." + name + " = append([]uint32(nil), field.GetUint32Array()...)"}
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return []string{"m." + name + " = append([]int32(nil), field.GetInt32Array()...)"}
		case protoreflect.BoolKind:
			return []string{"m." + name + " = append([]bool(nil), field.GetBoolArray()...)"}
		case protoreflect.EnumKind:
			return []string{"m." + name + " = append(" + exGoType(imports, field) + "(nil), protocache.CastEnumArray[" + strings.TrimPrefix(exGoType(imports, field), "[]") + "](field.GetEnumValueArray())...)"}
		}
		panic("unsupported repeated ex type")
	case field.Desc.Kind() == protoreflect.MessageKind:
		call := exAliasCtor(imports, field.Message.Desc)
		if !exIsAliasMessage(field.Message.Desc) {
			call = "TO_" + field.Message.GoIdent.GoName + "EX"
			if pkg, got := imports[string(field.Message.GoIdent.GoImportPath)]; got {
				call = pkg + "." + call
			}
		}
		return []string{"m." + name + " = " + call + "(field.GetObject())"}
	case field.Desc.Kind() == protoreflect.BytesKind:
		return []string{
			"if data := field.GetBytes(); data != nil {",
			"	m." + name + " = append([]byte(nil), data...)",
			"}",
		}
	case field.Desc.Kind() == protoreflect.StringKind:
		return []string{"m." + name + " = field.GetString()"}
	case field.Desc.Kind() == protoreflect.DoubleKind:
		return []string{"m." + name + " = field.GetFloat64()"}
	case field.Desc.Kind() == protoreflect.FloatKind:
		return []string{"m." + name + " = field.GetFloat32()"}
	case field.Desc.Kind() == protoreflect.Uint64Kind || field.Desc.Kind() == protoreflect.Fixed64Kind:
		return []string{"m." + name + " = field.GetUint64()"}
	case field.Desc.Kind() == protoreflect.Int64Kind || field.Desc.Kind() == protoreflect.Sint64Kind || field.Desc.Kind() == protoreflect.Sfixed64Kind:
		return []string{"m." + name + " = field.GetInt64()"}
	case field.Desc.Kind() == protoreflect.Uint32Kind || field.Desc.Kind() == protoreflect.Fixed32Kind:
		return []string{"m." + name + " = field.GetUint32()"}
	case field.Desc.Kind() == protoreflect.Int32Kind || field.Desc.Kind() == protoreflect.Sint32Kind || field.Desc.Kind() == protoreflect.Sfixed32Kind:
		return []string{"m." + name + " = field.GetInt32()"}
	case field.Desc.Kind() == protoreflect.BoolKind:
		return []string{"m." + name + " = field.GetBool()"}
	case field.Desc.Kind() == protoreflect.EnumKind:
		return []string{"m." + name + " = " + exGoType(imports, field) + "(field.GetEnumValue())"}
	default:
		panic("unsupported ex type")
	}
}

func exAliasDecodeExpr(field *protogen.Field, imports map[string]string, target, data string) []string {
	switch {
	case field.Desc.IsMap():
		valueField := exMapValueField(field)
		lines := []string{
			"pack := protocache.AsMap(" + data + ")",
			target + " = make(" + exGoType(imports, field) + ", int(pack.Size()))",
			"for i := uint32(0); i < pack.Size(); i++ {",
			"	keyField := pack.Key(i)",
		}
		switch field.Desc.MapKey().Kind() {
		case protoreflect.StringKind:
			lines = append(lines, "	key := keyField.GetString()")
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			lines = append(lines, "	key := keyField.GetUint64()")
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			lines = append(lines, "	key := keyField.GetInt64()")
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			lines = append(lines, "	key := keyField.GetUint32()")
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			lines = append(lines, "	key := keyField.GetInt32()")
		}
		lines = append(lines, "	valField := pack.Value(i)")
		switch valueField.Desc.Kind() {
		case protoreflect.MessageKind:
			lines = append(lines, target+"[key] = "+exAliasCtor(imports, valueField.Desc.Message())+"(valField.GetObject())")
		case protoreflect.BytesKind:
			lines = append(lines,
				"	if raw := valField.GetBytes(); raw != nil {",
				"		"+target+"[key] = append([]byte(nil), raw...)",
				"	} else {",
				"		"+target+"[key] = nil",
				"	}",
			)
		case protoreflect.StringKind:
			lines = append(lines, target+"[key] = valField.GetString()")
		case protoreflect.DoubleKind:
			lines = append(lines, target+"[key] = valField.GetFloat64()")
		case protoreflect.FloatKind:
			lines = append(lines, target+"[key] = valField.GetFloat32()")
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			lines = append(lines, target+"[key] = valField.GetUint64()")
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			lines = append(lines, target+"[key] = valField.GetInt64()")
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			lines = append(lines, target+"[key] = valField.GetUint32()")
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			lines = append(lines, target+"[key] = valField.GetInt32()")
		case protoreflect.BoolKind:
			lines = append(lines, target+"[key] = valField.GetBool()")
		case protoreflect.EnumKind:
			lines = append(lines, target+"[key] = "+exGoType(imports, valueField)+"(valField.GetEnumValue())")
		default:
			panic("unsupported alias map value type")
		}
		lines = append(lines, "}")
		return lines
	case field.Desc.IsList():
		switch field.Desc.Kind() {
		case protoreflect.MessageKind:
			return []string{
				"arr := protocache.AsArray(" + data + ")",
				target + " = make(" + exGoType(imports, field) + ", int(arr.Size()))",
				"for i := uint32(0); i < arr.Size(); i++ {",
				"	elem := arr.Get(i)",
				"	" + target + "[i] = " + exAliasCtor(imports, field.Message.Desc) + "(elem.GetObject())",
				"}",
			}
		case protoreflect.BytesKind:
			return []string{
				"arr := protocache.AsBytesArray(" + data + ")",
				target + " = make(" + exGoType(imports, field) + ", int(arr.Size()))",
				"for i := uint32(0); i < arr.Size(); i++ {",
				"	if raw := arr.Get(i); raw != nil {",
				"		" + target + "[i] = append([]byte(nil), raw...)",
				"	}",
				"}",
			}
		case protoreflect.StringKind:
			return []string{
				"arr := protocache.AsStringArray(" + data + ")",
				target + " = make(" + exGoType(imports, field) + ", int(arr.Size()))",
				"for i := uint32(0); i < arr.Size(); i++ {",
				"	" + target + "[i] = arr.Get(i)",
				"}",
			}
		case protoreflect.DoubleKind:
			return []string{"arr := protocache.AsFloat64Array(" + data + ")", target + " = append(" + exGoType(imports, field) + "(nil), arr.Raw()...)"}
		case protoreflect.FloatKind:
			return []string{"arr := protocache.AsFloat32Array(" + data + ")", target + " = append(" + exGoType(imports, field) + "(nil), arr.Raw()...)"}
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return []string{"arr := protocache.AsUint64Array(" + data + ")", target + " = append(" + exGoType(imports, field) + "(nil), arr.Raw()...)"}
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return []string{"arr := protocache.AsInt64Array(" + data + ")", target + " = append(" + exGoType(imports, field) + "(nil), arr.Raw()...)"}
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return []string{"arr := protocache.AsUint32Array(" + data + ")", target + " = append(" + exGoType(imports, field) + "(nil), arr.Raw()...)"}
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return []string{"arr := protocache.AsInt32Array(" + data + ")", target + " = append(" + exGoType(imports, field) + "(nil), arr.Raw()...)"}
		case protoreflect.BoolKind:
			return []string{"arr := protocache.AsBoolArray(" + data + ")", target + " = append(" + exGoType(imports, field) + "(nil), arr.Raw()...)"}
		case protoreflect.EnumKind:
			return []string{target + " = append(" + exGoType(imports, field) + "(nil), protocache.AsEnumArray[" + strings.TrimPrefix(exGoType(imports, field), "[]") + "](" + data + ").Raw()...)"}
		default:
			panic("unsupported alias repeated type")
		}
	default:
		panic("alias must be list or map")
	}
}

func exAliasEncodeExpr(field *protogen.Field, imports map[string]string, access string) []string {
	switch {
	case field.Desc.IsMap():
		valueField := exMapValueField(field)
		lines := []string{
			"if len(" + access + ") == 0 {",
			"	return []uint32{5 << 28}, nil",
			"}",
			"keys := make([][]uint32, 0, len(" + access + "))",
			"vals := make([][]uint32, 0, len(" + access + "))",
			"for k, v := range " + access + " {",
			"	keyPart, err := " + exMapKeyEncodeExpr(field.Desc.MapKey(), "k"),
			"	if err != nil {",
			"		return nil, err",
			"	}",
			"	valPart, err := " + exEncodeExpr(valueField, "v"),
			"	if err != nil {",
			"		return nil, err",
			"	}",
			"	keys = append(keys, keyPart)",
			"	vals = append(vals, valPart)",
			"}",
			"return protocache.EncodeMapParts(keys, vals, " + fmt.Sprint(field.Desc.MapKey().Kind() == protoreflect.StringKind) + ")",
		}
		return lines
	case field.Desc.IsList():
		switch field.Desc.Kind() {
		case protoreflect.MessageKind:
			return []string{
				"if len(" + access + ") == 0 {",
				"	return []uint32{1}, nil",
				"}",
				"return protocache.EncodeObjectArray(len(" + access + "), func(i int) ([]uint32, error) { return " + exAliasSerializeExpr(access+"[i]") + " })",
			}
		case protoreflect.BytesKind:
			return []string{
				"if len(" + access + ") == 0 {",
				"	return []uint32{1}, nil",
				"}",
				"return protocache.EncodeBytesArray([][]byte(" + access + "))",
			}
		case protoreflect.StringKind:
			return []string{
				"if len(" + access + ") == 0 {",
				"	return []uint32{1}, nil",
				"}",
				"return protocache.EncodeStringArray([]string(" + access + "))",
			}
		case protoreflect.DoubleKind:
			return []string{"if len(" + access + ") == 0 { return []uint32{1}, nil }", "return protocache.EncodeFloat64Array([]float64(" + access + "))"}
		case protoreflect.FloatKind:
			return []string{"if len(" + access + ") == 0 { return []uint32{1}, nil }", "return protocache.EncodeFloat32Array([]float32(" + access + "))"}
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return []string{"if len(" + access + ") == 0 { return []uint32{1}, nil }", "return protocache.EncodeUint64Array([]uint64(" + access + "))"}
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return []string{"if len(" + access + ") == 0 { return []uint32{1}, nil }", "return protocache.EncodeInt64Array([]int64(" + access + "))"}
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return []string{"if len(" + access + ") == 0 { return []uint32{1}, nil }", "return protocache.EncodeUint32Array([]uint32(" + access + "))"}
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return []string{"if len(" + access + ") == 0 { return []uint32{1}, nil }", "return protocache.EncodeInt32Array([]int32(" + access + "))"}
		case protoreflect.BoolKind:
			return []string{"if len(" + access + ") == 0 { return []uint32{1}, nil }", "return protocache.EncodeBoolArray([]bool(" + access + "))"}
		case protoreflect.EnumKind:
			return []string{"if len(" + access + ") == 0 { return []uint32{1}, nil }", "return protocache.EncodeEnumArray([]" + strings.TrimPrefix(exGoType(imports, field), "[]") + "(" + access + "))"}
		default:
			panic("unsupported alias repeated type")
		}
	default:
		panic("alias must be list or map")
	}
}

func GenEXAlias(g *protogen.GeneratedFile, imports map[string]string, one *protogen.Message) error {
	field := one.Fields[0]
	typeName := one.GoIdent.GoName + "EX"
	g.P()
	g.P("type ", typeName, " ", exGoType(imports, field))
	g.P()
	g.P("func TO_", one.GoIdent.GoName, "EX(data []byte) ", typeName, " {")
	g.P("	var out ", typeName)
	for _, line := range exAliasDecodeExpr(field, imports, "out", "data") {
		g.P("	", line)
	}
	g.P("	return out")
	g.P("}")
	g.P()
	g.P("func (x ", typeName, ") Serialize() ([]byte, error) {")
	g.P("	words, err := serialize", one.GoIdent.GoName, "EX(x)")
	g.P("	if err != nil {")
	g.P("		return nil, err")
	g.P("	}")
	g.P("	return protocache.WordsToBytes(words), nil")
	g.P("}")
	g.P()
	g.P("func serialize", one.GoIdent.GoName, "EX(x ", typeName, ") ([]uint32, error) {")
	for _, line := range exAliasEncodeExpr(field, imports, "x") {
		g.P("	", line)
	}
	g.P("}")
	return nil
}

func exSetExpr(field *protogen.Field, imports map[string]string, name string) []string {
	switch {
	case field.Desc.IsMap() && field.Desc.MapValue().Kind() == protoreflect.BytesKind:
		return []string{
			"if v == nil {",
			"	m." + name + " = nil",
			"} else {",
			"	m." + name + " = make(" + exGoType(imports, field) + ", len(v))",
			"	for k, one := range v {",
			"		if one != nil {",
			"			m." + name + "[k] = append([]byte(nil), one...)",
			"		}",
			"	}",
			"}",
		}
	case field.Desc.IsMap():
		return []string{
			"if v == nil {",
			"	m." + name + " = nil",
			"} else {",
			"	m." + name + " = make(" + exGoType(imports, field) + ", len(v))",
			"	for k, one := range v {",
			"		m." + name + "[k] = one",
			"	}",
			"}",
		}
	case field.Desc.IsList() && field.Desc.Kind() == protoreflect.BytesKind:
		return []string{
			"if v == nil {",
			"	m." + name + " = nil",
			"} else {",
			"	m." + name + " = make([][]byte, len(v))",
			"	for i := range v {",
			"		if v[i] != nil {",
			"			m." + name + "[i] = append([]byte(nil), v[i]...)",
			"		}",
			"	}",
			"}",
		}
	case field.Desc.IsList():
		return []string{
			"if v == nil {",
			"	m." + name + " = nil",
			"} else {",
			"	m." + name + " = append(m." + name + "[:0], v...)",
			"}",
		}
	case field.Desc.Kind() == protoreflect.BytesKind:
		return []string{
			"if v == nil {",
			"	m." + name + " = nil",
			"} else {",
			"	m." + name + " = append([]byte(nil), v...)",
			"}",
		}
	default:
		return []string{"m." + name + " = v"}
	}
}

func GenEXMessages(g *protogen.GeneratedFile, imports map[string]string, list []*protogen.Message) error {
	for _, one := range list {
		if one.Desc.IsMapEntry() {
			continue
		}
		if err := GenEXMessages(g, imports, one.Messages); err != nil {
			return err
		}
		if _, ok := aliasBook[string(one.Desc.FullName())]; ok {
			if err := GenEXAlias(g, imports, one); err != nil {
				return err
			}
			continue
		}

		fields := make([]*protogen.Field, len(one.Fields))
		copy(fields, one.Fields)
		order := slices.Order[*protogen.Field]{
			Less: func(a, b *protogen.Field) bool {
				return a.Desc.Number() < b.Desc.Number()
			},
		}
		order.Sort(fields)
		maxID := protoreflect.FieldNumber(0)
		for _, field := range fields {
			if maxID < field.Desc.Number() {
				maxID = field.Desc.Number()
			}
		}

		g.P()
		g.P("type ", one.GoIdent.GoName, "EX struct {")
		g.P("	__ protocache.MessageEX")
		for _, field := range fields {
			if !exSupportsField(field.Desc) {
				continue
			}
			g.P("	", exFieldName(field), " ", exGoType(imports, field))
		}
		g.P("}")
		g.P()
		g.P("func TO_", one.GoIdent.GoName, "EX(data []byte) *", one.GoIdent.GoName, "EX {")
		g.P("	out := &", one.GoIdent.GoName, "EX{}")
		g.P("	out.__.Init(data)")
		g.P("	return out")
		g.P("}")
		g.P()
		g.P("func (m *", one.GoIdent.GoName, "EX) HasBase() bool { return m.__.HasBase() }")
		g.P()
		g.P("func (m *", one.GoIdent.GoName, "EX) IsValid() bool { return m.HasBase() }")
		g.P()
		g.P("func (m *", one.GoIdent.GoName, "EX) Serialize() ([]byte, error) {")
		g.P("	words, err := m.serializeWords()")
		g.P("	if err != nil {")
		g.P("		return nil, err")
		g.P("	}")
		g.P("	return protocache.WordsToBytes(words), nil")
		g.P("}")
		g.P()
		g.P("func (m *", one.GoIdent.GoName, "EX) serializeWords() ([]uint32, error) {")
		g.P("	if m == nil {")
		g.P("		return []uint32{0}, nil")
		g.P("	}")
		g.P("	parts := make([][]uint32, ", maxID, ")")
		for _, field := range fields {
			id := field.Desc.Number() - 1
			if exSupportsField(field.Desc) {
				name := exFieldName(field)
				g.P("	if m.__.IsVisited(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ", _FIELD_TOTAL_", one.GoIdent.GoName, ") {")
				if field.Desc.IsMap() {
					valField := exMapValueField(field)
					g.P("		if len(m.", name, ") != 0 {")
					g.P("			keys := make([][]uint32, 0, len(m.", name, "))")
					g.P("			vals := make([][]uint32, 0, len(m.", name, "))")
					g.P("			for k, v := range m.", name, " {")
					g.P("				keyPart, err := ", exMapKeyEncodeExpr(field.Desc.MapKey(), "k"))
					g.P("				if err != nil {")
					g.P("					return nil, err")
					g.P("				}")
					g.P("				valPart, err := ", exEncodeExpr(valField, "v"))
					g.P("				if err != nil {")
					g.P("					return nil, err")
					g.P("				}")
					if valField.Desc.Kind() == protoreflect.MessageKind {
						g.P("				if len(valPart) <= 1 {")
						g.P("					valPart = nil")
						g.P("				}")
					}
					g.P("				keys = append(keys, keyPart)")
					g.P("				vals = append(vals, valPart)")
					g.P("			}")
					g.P("			part, err := protocache.EncodeMapParts(keys, vals, ", field.Desc.MapKey().Kind() == protoreflect.StringKind, ")")
					g.P("			if err != nil {")
					g.P("				return nil, err")
					g.P("			}")
					g.P("			parts[", id, "] = part")
					g.P("		}")
				} else if field.Desc.IsList() {
					g.P("		if len(m.", name, ") != 0 {")
					g.P("			part, err := ", exEncodeExpr(field, "m."+name))
					g.P("			if err != nil {")
					g.P("				return nil, err")
					g.P("			}")
					g.P("			parts[", id, "] = part")
					g.P("		}")
				} else if field.Desc.Kind() == protoreflect.MessageKind {
					g.P("		if m.", name, " != nil {")
					g.P("			part, err := ", exEncodeExpr(field, "m."+name))
					g.P("			if err != nil {")
					g.P("				return nil, err")
					g.P("			}")
					g.P("			if len(part) > 1 {")
					g.P("				parts[", id, "] = part")
					g.P("			}")
					g.P("		}")
				} else {
					g.P("		part, err := ", exEncodeExpr(field, "m."+name))
					g.P("		if err != nil {")
					g.P("			return nil, err")
					g.P("		}")
					g.P("		parts[", id, "] = part")
				}
				g.P("	} else if raw := m.__.RawPart(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), "); len(raw) != 0 {")
				g.P("		parts[", id, "] = protocache.BytesToWords(raw)")
				g.P("	}")
			} else {
				g.P("	if raw := m.__.RawPart(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), "); len(raw) != 0 {")
				g.P("		parts[", id, "] = protocache.BytesToWords(raw)")
				g.P("	}")
			}
		}
		g.P("	return protocache.EncodeMessageParts(parts)")
		g.P("}")

		for _, field := range fields {
			if !exSupportsField(field.Desc) {
				continue
			}
			name := exFieldName(field)
			g.P()
			g.P("func (m *", one.GoIdent.GoName, "EX) Get", field.GoName, "() ", exGoType(imports, field), " {")
			g.P("	if m.__.IsVisited(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ", _FIELD_TOTAL_", one.GoIdent.GoName, ") {")
			g.P("		return m.", name)
			g.P("	}")
			g.P("	field := m.__.RawField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
			for _, line := range exGetExpr(field, imports, name) {
				g.P("	", line)
			}
			g.P("	m.__.Visit(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ", _FIELD_TOTAL_", one.GoIdent.GoName, ")")
			g.P("	return m.", name)
			g.P("}")
			g.P()
			g.P("func (m *", one.GoIdent.GoName, "EX) Set", field.GoName, "(v ", exGoType(imports, field), ") {")
			for _, line := range exSetExpr(field, imports, name) {
				g.P("	", line)
			}
			g.P("	m.__.Visit(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ", _FIELD_TOTAL_", one.GoIdent.GoName, ")")
			g.P("}")
		}
	}
	return nil
}

func GenEXFile(gen *protogen.Plugin, file *protogen.File) error {
	imports := make(map[string]string)
	CollectImports(string(file.GoImportPath), file.Messages, imports)
	pkgs := make([]string, 0, len(imports))
	for one := range imports {
		pkgs = append(pkgs, one)
	}
	slices.Sort(pkgs)
	for i, one := range pkgs {
		imports[one] = fmt.Sprintf("p%d", i+1)
	}

	filename := file.GeneratedFilenamePrefix + ".pc-ex.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("import (")
	g.P(`	"github.com/peterrk/protocache-go"`)
	for name, mark := range imports {
		g.P("	", mark, ` "`, name, `"`)
	}
	g.P(")")
	return GenEXMessages(g, imports, file.Messages)
}
