package generate

import (
	"fmt"
	"github.com/efucloud/common"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

type TypeScript struct {
	structTypes map[string]reflect.Type
	kinds       map[string]string
	consts      []Const
	structMap   map[string]StructInfo
}
type Const struct {
	Model       string
	Raw         string
	Name        string
	Type        string
	Enum        []interface{}
	Description string
}
type StructInfo struct {
	Name        string
	Description string
	Fields      []FieldInfo
}
type FieldInfo struct {
	Name        string
	JsonName    string
	Kind        string
	Description string
	Required    bool
	EnumValues  []interface{}
	Default     string
	Enum        string
	Length      string
}

func (script *TypeScript) AddStruct(st reflect.Type) {
	script.structTypes[strings.TrimPrefix(st.Name(), "*")] = st
}
func NewTypeScript() *TypeScript {
	script := &TypeScript{}
	script.structTypes = make(map[string]reflect.Type)
	script.kinds = make(map[string]string)
	script.structMap = make(map[string]StructInfo)
	script.kinds[reflect.Bool.String()] = "boolean"
	script.kinds[reflect.Interface.String()] = "any"
	script.kinds[reflect.Int.String()] = "number"
	script.kinds[reflect.Int8.String()] = "number"
	script.kinds[reflect.Int16.String()] = "number"
	script.kinds[reflect.Int32.String()] = "number"
	script.kinds[reflect.Int64.String()] = "number"
	script.kinds[reflect.Uint.String()] = "number"
	script.kinds[reflect.Uint8.String()] = "number"
	script.kinds[reflect.Uint16.String()] = "number"
	script.kinds[reflect.Uint32.String()] = "number"
	script.kinds[reflect.Uint64.String()] = "number"
	script.kinds[reflect.Float32.String()] = "number"
	script.kinds[reflect.Float64.String()] = "number"
	script.kinds[reflect.String.String()] = "string"
	script.kinds["Time"] = "string"
	script.kinds["JSONMap"] = "{[key: string]: string}"
	script.kinds["ArrayString"] = "string[]"
	script.kinds["ArrayUint"] = "number[]"

	return script
}
func (script *TypeScript) GenerateToDir(dir string) {
	_ = os.WriteFile(path.Join(dir, time.Now().Format(time.DateOnly)+"types.d.ts"), []byte(script.Generate()), os.ModePerm)
}
func (script *TypeScript) Generate() (content string) {
	for _, item := range script.structTypes {
		script.structMap[item.Name()] = script.extractStructFields(item)

	}
	//生成常量
	for _, item := range script.consts {
		content += fmt.Sprintf("// %s \n", item.Description)
		if item.Type == "string" {
			var values []string
			for _, v := range item.Enum {
				values = append(values, fmt.Sprintf("'%v'", v))
			}
			content += fmt.Sprintf("type %s = %s ;\n", item.Name, strings.Join(values, "|"))

		} else {
			var values []string
			for _, v := range item.Enum {
				values = append(values, fmt.Sprintf("%v", v))
			}
			content += fmt.Sprintf("type %s = %s ;\n", item.Name, strings.Join(values, "|"))
		}
	}
	for _, item := range script.structMap {
		content += fmt.Sprintf("// %s %s \n", item.Name, item.Description)
		content += fmt.Sprintf("type %s = { \n", item.Name)
		for _, field := range item.Fields {
			if len(field.Description) > 0 {
				content += fmt.Sprintf("  // %s\n", field.Description)
			}
			if len(field.Default) > 0 {
				content += fmt.Sprintf("  // 默认值: %s\n", field.Default)
			}
			if len(field.EnumValues) > 0 {
				var ev []string
				for _, i := range field.EnumValues {
					ev = append(ev, fmt.Sprintf("%v", i))
				}
				content += fmt.Sprintf("  // 可选值: %v\n", strings.Join(ev, ";"))
			}
			if len(field.Length) > 0 {
				content += fmt.Sprintf("  // 长度: %s\n", field.Length)
			}
			content += fmt.Sprintf("  %s", field.JsonName)
			if !field.Required {
				content += "?"
			}
			content += fmt.Sprintf(": %s;\n", field.Kind)
		}
		content += "}; \n"
	}
	return content
}

func (script *TypeScript) extractStructFields(item reflect.Type) (structInfo StructInfo) {
	structInfo.Name = item.Name()
	for i := 0; i < item.NumField(); i++ {
		var field FieldInfo
		//获取描述
		if item.Field(i).Name == "Doc" {
			structInfo.Description = item.Field(i).Tag.Get("description")
			continue
		}

		field.Name = item.Field(i).Name
		field.JsonName = item.Field(i).Tag.Get("json")
		if field.JsonName == "-" {
			//忽略不输出的字段
			continue
		} else if field.JsonName == "" {
			field.JsonName = field.Name
		}
		field.Description = item.Field(i).Tag.Get("description")
		gorm := item.Field(i).Tag.Get("gorm")
		if len(gorm) > 0 && gorm != "" {
			for _, it := range strings.Split(gorm, ";") {
				sp := strings.Split(it, ":")
				if len(sp) == 2 {
					switch sp[0] {
					case "default":
						field.Default = sp[1]
					case "type":
						if strings.HasPrefix(sp[0], "varchar(") {
							lens := strings.TrimPrefix(sp[0], "varchar(")
							if len(lens) > 2 {
								field.Length = lens[:len(lens)-2]
							}
						}
					}
				}
			}
		}
		validate := item.Field(i).Tag.Get("validate")
		if len(validate) > 0 && validate != "-" {
			field.Required = true
		}
		enum := item.Field(i).Tag.Get("enum")
		if len(enum) > 0 {
			var con Const
			con.Type = item.Field(i).Type.String()
			con.Raw = enum
			con.Model = item.Name()
			con.Name = fmt.Sprintf("%s%s", item.Name(), item.Field(i).Name)
			field.Kind = con.Name
			con.Description = field.Description
			field.Enum = enum
			values := strings.Split(enum, "|")
			if item.Field(i).Type.String() == "string" {
				for _, i := range values {
					con.Enum = append(con.Enum, i)
					field.EnumValues = append(field.EnumValues, i)
				}
			} else if item.Field(i).Type.String() == "uint" {
				for _, i := range values {
					con.Enum = append(con.Enum, common.StringToInt(i))
					field.EnumValues = append(field.EnumValues, common.StringToInt(i))
				}
			}
			script.consts = append(script.consts, con)
		}
		if len(field.Kind) == 0 {
			field.Kind, _ = script.kinds[item.Field(i).Type.String()]
		}
		if len(field.Kind) == 0 {
			kindName := strings.TrimPrefix(item.Field(i).Type.String(), "*")
			if strings.Contains(kindName, ".") {
				sp := strings.Split(kindName, ".")
				spLen := len(sp)
				kindName = sp[spLen-1]
			}
			field.Kind, _ = script.kinds[kindName]
			if len(field.Kind) == 0 {
				switch item.Field(i).Type.Kind() {
				case reflect.Struct:
					field.Kind = kindName
					if _, exist := script.structTypes[field.Kind]; !exist {
						if _, ex := script.structMap[field.Kind]; !ex {
							script.structMap[field.Kind] = script.extractStructFields(item.Field(i).Type.Elem())
						}
						//script.structTypes[field.Kind] =
					}
				case reflect.Pointer:
					field.Kind = kindName
					if _, exist := script.structTypes[field.Kind]; !exist {
						if _, ex := script.structMap[field.Kind]; !ex {
							script.structMap[field.Kind] = script.extractStructFields(item.Field(i).Type.Elem())
						}
					}
				}
			}

		}
		structInfo.Fields = append(structInfo.Fields, field)
	}
	return
}