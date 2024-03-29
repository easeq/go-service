package generator

var tmpl = `
// Code generated by protoc-gen-go_service. DO NOT EDIT.
// versions:
//	protoc	v{{.Gen.Request.CompilerVersion.Major}}.{{.Gen.Request.CompilerVersion.Minor}}.{{.Gen.Request.CompilerVersion.Patch}}

package {{.GoPackageName}}

import (
	"context"

	"github.com/easeq/go-service/client"
	{{range $import, $val := .Imports -}}
	"{{$import}}"
	{{end}}
)

{{$pkg := .GoPackageName}}
{{$streams := .Streams}}

{{range .Services}}
{{$serviceName := .GoName}}
{{$serviceNameCamel := (camelCase $serviceName)}}
{{$serviceFullName := .Desc.FullName}}
type {{$serviceName}}GSClient interface {
	{{range .Methods -}}
	{{if and (not .Desc.IsStreamingServer) (not .Desc.IsStreamingClient) -}}
	{{.GoName}}(ctx context.Context, in *{{.Input.GoIdent.GoName}}, opts ...client.CallOption) (*{{.Output.GoIdent.GoName}}, error)
	{{else -}}
	{{.GoName}}(ctx context.Context, in *{{.Input.GoIdent.GoName}}, opts ...client.CallOption) (client.StreamClient, error)
	{{end}}
	{{- end}}
}

type {{$serviceNameCamel}}GSClient struct {
	name 	string
	opts   	[]client.DialOption
	client.Client
}

func New{{$serviceName}}GSClient(name string, client client.Client, opts ...client.DialOption) {{$serviceName}}GSClient {
	return &{{$serviceNameCamel}}GSClient{name, opts, client}
}


func (sc *{{$serviceNameCamel}}GSClient) GetServiceName() string {
	return sc.name
}

func (sc *{{$serviceNameCamel}}GSClient) GetDialOptions() []client.DialOption {
	return sc.opts
}

{{range .Methods}}
{{$methodName := .GoName}}
{{$inputName := .Input.GoIdent.GoName}}
{{$outputName := .Output.GoIdent.GoName}}
{{if and (not .Desc.IsStreamingServer) (not .Desc.IsStreamingClient)}}
func (sc *{{$serviceNameCamel}}GSClient) {{$methodName}}(ctx context.Context, in *{{$inputName}}, opts ...client.CallOption) (*{{$outputName}}, error) {
	res := new({{$outputName}})
	err := sc.Call(ctx, sc, "/{{$serviceFullName}}/{{$methodName}}", in, res, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
{{else}}
func (sc *{{$serviceNameCamel}}GSClient) {{$methodName}}(ctx context.Context, in *{{$inputName}}, opts ...client.CallOption) (client.StreamClient, error) {
	stream, err := sc.Stream(ctx, sc, &{{$serviceName}}_ServiceDesc.Streams[{{index $streams (printf "%s%s" .Parent.GoName .GoName)}}], "/{{$serviceFullName}}/{{$methodName}}", in, opts...)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
{{end}}
{{end}}
{{end}}

{{range $serviceName, $tag := .RegistryTags -}}
{{if $tag}}
func Add{{$serviceName}}ServerRegistryTags(server server.Server) {
	tags := []string{
		"traefik.http.routers.{{$tag.Name}}.rule=Host(` + "`" + `" + os.Getenv("{{$tag.Host}}") + "` + "`" + `) && PathPrefix(` + "`" + `{{$tag.Path}}` + "`" + `)",
		{{if $tag.Stripprefix -}}
		"traefik.http.middlewares.{{$tag.Name}}-stripprefix.stripprefix.prefixes={{$tag.Stripprefix}}",
		"traefik.http.routers.{{$tag.Name}}.middlewares={{$tag.Name}}-stripprefix@consulcatalog",
		{{- end}}
	}

	if len(tags) > 0 {
		tags = append(tags, "traefik.enable=true")
	}

	server.AddRegistryTags(tags...)
}
{{- end}}
{{- end}}
`
