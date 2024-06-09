package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"

	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"
	"time"
)

const apiHasPathParamsTemplate = "_{{_ .description _}}_" +
	"export const _{{_ .functionName _}}_  = async (params?: any) => {\n" +
	"    _{{_ .extractPathParams _}}_\n" +
	"    return request(_{{_ .globalAPI _}}_, { method: '_{{_ .method _}}_' _{{_ if .queryParams _}}_ , params: rest _{{_ end _}}_ _{{_ if .body _}}_ , body: rest _{{_ end _}}_ });\n" +
	"};\n"

const GlobalApiName = "GlobalApiName"

type RestAPI struct {
	// 从接口中自动提取
	structTypes        map[string]reflect.Type
	routes             []restful.Route
	apis               map[string]ApiData
	globalApiName      string
	generateTypescript bool
	files              map[string]string
}

type ApiData struct {
	DocumentName string
	Name         string
	Doc          string
	Notes        string
	Path         string
	Method       string
	Parameters   map[string]Parameters
	Response     map[int]string
}

func (api ApiData) String() string {
	return api.Method + " " + api.Path

}

type Parameters struct {
	Name        string // 参数名
	DataType    string
	Position    string // query path
	Description string
	Required    bool
	Enum        []string
	Default     string
}

func NewRestAPI(frontApiName string, generateTypescript bool) *RestAPI {
	api := &RestAPI{
		structTypes:        make(map[string]reflect.Type),
		apis:               make(map[string]ApiData),
		globalApiName:      frontApiName,
		generateTypescript: generateTypescript,
		files:              make(map[string]string),
	}
	if len(api.globalApiName) == 0 {
		api.globalApiName = GlobalApiName
	}
	return api
}
func (rest *RestAPI) AddRoute(route restful.Route) {
	rest.routes = append(rest.routes, route)
}
func GetStructFieldDescription(item reflect.Type) string {
	result, _ := json.Marshal(ExtractStructFieldDescription(item))
	return string(result)
}
func ExtractStructFieldDescription(item reflect.Type) (result map[string]string) {
	result = make(map[string]string)
	for i := 0; i < item.NumField(); i++ {
		jsonName := item.Field(i).Tag.Get("json")
		if len(jsonName) == 0 {
			jsonName = item.Field(i).Name
		}
		description := item.Field(i).Tag.Get("description")
		if len(description) == 0 {
			description = item.Field(i).Name
		}
		result[jsonName] = description
	}
	return
}

func (rest *RestAPI) GenerateToDir(dir string) {
	apiContent, typeContent := rest.Generate()
	if len(apiContent) > 0 {
		_ = os.WriteFile(path.Join(dir, time.Now().Format(time.DateOnly)+".api.ts"), []byte(apiContent), os.ModePerm)
	}
	if len(typeContent) > 0 {
		_ = os.WriteFile(path.Join(dir, time.Now().Format(time.DateOnly)+".types.d.ts"), []byte(apiContent), os.ModePerm)
	}
}
func (rest *RestAPI) Generate() (apiContent, typeContent string) {

	rest.ParserRoutes()
	// 生成typescript定义
	if rest.generateTypescript {
		typs := NewTypeScript()
		for _, item := range rest.structTypes {
			typs.AddStruct(item)
		}
		typeContent = typs.Generate()
	}
	// 生成api
	for _, api := range rest.apis {
		rest.files[api.DocumentName] += rest.generateOneApi(api)
	}

	return
}
func (rest *RestAPI) generateOneApi(api ApiData) (content string) {
	var description string
	description += fmt.Sprintf("// %s\n", api.Doc)
	description += fmt.Sprintf("// %s\n", api.Notes)
	for code, data := range api.Response {
		description += fmt.Sprintf("// 响应码: %d  响应数据: %s\n", code, data)
	}
	var pathParams []string
	for name, item := range api.Parameters {
		required := "否"
		if item.Required {
			required = "是"
		}
		other := ""
		if len(item.Default) > 0 {
			other += fmt.Sprintf("默认值: %s", item.Default)
		}
		if len(item.Enum) > 0 {
			other += fmt.Sprintf(" 可选值: %s", strings.Join(item.Enum, ";"))
		}
		description += fmt.Sprintf("// 参数名: %s 参数类型: %s 参数位置: %s 是否必须: %s 参数说明: %s %s\n", name, item.DataType, item.Position, required, item.Description, other)
		if item.Position == "path" {
			pathParams = append(pathParams, item.Name)
		}
	}
	params := make(map[string]interface{})
	params["description"] = description
	params["functionName"] = api.Name
	params["method"] = strings.ToLower(api.Method)
	if len(pathParams) > 0 {
		params["extractPathParams"] = fmt.Sprintf("const { %s, ...rest } = params;", strings.Join(pathParams, ", "))
	}
	switch api.Method {
	case http.MethodPost, http.MethodPut:
		params["body"] = true
	case http.MethodGet, http.MethodDelete:
		params["queryParams"] = true
	}
	t, _ := template.New(time.Now().String()).Delims("_{{_", "_}}_").Parse(apiHasPathParamsTemplate)
	b := new(bytes.Buffer)
	err := t.Execute(b, params)
	if err == nil {
		content += b.String()
	}
	return content
}
func (rest *RestAPI) ParserRoutes() {
	for _, route := range rest.routes {
		var api ApiData
		api.Parameters = make(map[string]Parameters)
		api.Response = make(map[int]string)
		api.Doc = route.Doc
		api.Notes = route.Notes
		api.Path = route.Path
		api.Method = route.Method

		if name, exist := route.Metadata[rest.globalApiName]; exist {
			api.Name = fmt.Sprintf("%v", name)
		} else {
			// todo 驼峰
			api.Name = fmt.Sprintf("%s%s", strings.ToLower(route.Method), route.Operation)
		}
		if doc, ex := route.Metadata[restfulspec.KeyOpenAPITags]; ex {
			api.DocumentName = strings.ReplaceAll(fmt.Sprintf("%v", doc), "-", "_")
		} else {
			api.DocumentName = "api"
		}
		for _, param := range route.ParameterDocs {
			var p Parameters
			p.Description = param.Data().Description
			p.Name = param.Data().Name
			p.DataType = param.Data().DataType
			p.Required = param.Data().Required
			p.Default = param.Data().DefaultValue
			p.Enum = param.Data().PossibleValues
			switch param.Kind() {
			case restful.PathParameterKind:
				p.Position = "path"
			case restful.QueryParameterKind:
				p.Position = "query"
			case restful.BodyParameterKind:
				p.Position = "body"
			case restful.HeaderParameterKind:
				p.Position = "header"
			case restful.FormParameterKind:
				p.Position = "form"
			case restful.MultiPartFormParameterKind:
				p.Position = "multipart/form-data"
			}
			api.Parameters[p.Name] = p
			if route.ReadSample != nil {
				read := reflect.TypeOf(route.ReadSample)
				if read != nil {
					rest.structTypes[read.Name()] = read
				}
			}
			if route.WriteSample != nil {
				write := reflect.TypeOf(route.WriteSample)
				if write != nil {
					rest.structTypes[write.Name()] = write
				}
			}
			for _, res := range route.ResponseErrors {
				if res.Code == http.StatusOK || res.Code == http.StatusCreated {
					successRes := reflect.TypeOf(res.Model)
					if successRes != nil {
						rest.structTypes[successRes.Name()] = successRes
						api.Response[res.Code] = successRes.Name()
					}
				} else {
					if res.Model != nil && reflect.TypeOf(res.Model) != nil {
						api.Response[res.Code] = GetStructFieldDescription(reflect.TypeOf(res.Model))

					}
				}
			}
		}
		rest.apis[api.String()] = api
	}
}
