package main

import (
	"fmt"

	// "github.com/easeq/go-service/cmd/protoc-gen-go-service/options"
	gs_generator "github.com/easeq/go-service/protoc-gen-go-service/generator"
	"github.com/easeq/go-service/protoc-gen-go-service/options"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func main() {
	protogen.Options{}.Run(run)
}

func run(gen *protogen.Plugin) error {
	for _, f := range gen.Files {
		if !f.Generate || len(f.Services) == 0 {
			continue
		}

		skip := true
		ss := make(map[string]int)
		registryTags := map[string]*options.RegistryTag{}
		imports := make(map[string]bool)
		for _, service := range f.Services {
			// Skip file generation if there are no methods
			if len(service.Methods) != 0 {
				skip = false
			}

			// Create a streaming service index map
			index := 0
			for _, method := range service.Methods {
				if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
					ss[fmt.Sprintf("%s%s", service.GoName, method.GoName)] = index
					index++
				}
			}

			// Get service registry tags for traefik
			opts := service.Desc.Options()
			if !opts.ProtoReflect().IsValid() {
				continue
			}

			tags := proto.GetExtension(opts, options.E_RegistryTag)
			registryTag, ok := tags.(*options.RegistryTag)
			if !ok {
				panic("Invalid registry tags")
			}

			registryTags[service.GoName] = registryTag
			imports["os"] = true
			imports["github.com/easeq/go-service/server"] = true
		}

		if skip {
			continue
		}

		generator := &gs_generator.Generator{
			FilenamePrefix: f.GeneratedFilenamePrefix,
			GoPackageName:  f.GoPackageName,
			GoImportPath:   f.GoImportPath,
			Gen:            gen,
			Streams:        ss,
			RegistryTags:   registryTags,
			Services:       f.Services,
			Imports:        imports,
		}
		generator.GenerateFile()
	}
	return nil
}
