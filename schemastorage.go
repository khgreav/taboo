package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kaptinlin/jsonschema"
)

type SchemaStorage struct {
	compiler *jsonschema.Compiler
}

var (
	schemaStorage *SchemaStorage
	ssOnce        sync.Once
)

func GetSchemaStorage() (*SchemaStorage, error) {
	var err error
	ssOnce.Do(func() {
		compiler := jsonschema.NewCompiler()
		schemaStorage = &SchemaStorage{
			compiler: compiler,
		}
		err = schemaStorage.loadSchemas("schemas/")
	})
	if err != nil {
		return nil, err
	}
	return schemaStorage, nil
}

func (ss *SchemaStorage) loadSchemas(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("Failed to read schema dir.", "dir", dir, "err", err)
		return fmt.Errorf("failed to read schema dir %s: %s", dir, "err")
	}

	for _, file := range files {
		// Skip directories and non-json files
		if file.IsDir() || !strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			continue
		}

		path := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			slog.Warn("Failed to read schema file.", "file", path, "err", err)
			continue
		}

		// compile schema
		schema, err := ss.compiler.Compile(content)
		if err != nil {
			slog.Warn("Failed to compile schema file.", "file", path, "err", err)
			continue
		}
		slog.Debug("Added new schema to storage.", "id", schema.ID)
	}

	return nil
}

func (ss SchemaStorage) validate(messageType MessageType, data []byte) error {
	schema, err := ss.compiler.GetSchema(string(messageType))
	if err != nil {
		return fmt.Errorf("could not get schema to validate %s message: %w", messageType, err)
	}
	result := schema.Validate(data)
	if !result.IsValid() {
		var errors []string
		for field, errMsg := range result.Errors {
			errors = append(errors, fmt.Sprintf("%s: %s", field, errMsg))
		}
		return fmt.Errorf("message is not valid against schema: %s", strings.Join(errors, "; "))
	}
	return nil
}
