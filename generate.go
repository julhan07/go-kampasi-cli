package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type ModelField struct {
	Name     string
	Type     string
	JsonTag  string
	Validate string
}

type SQLField struct {
	Name        string
	Type        string
	Constraints string
}
type TemplateData struct {
	Name        string
	TableName   string
	LowerName   string
	PackagePath string
	ModelField  []ModelField
	SQLField    []SQLField
	Routers     []string
	Add         func(int, int) int
}

func getGOPATH() (string, error) {
	cmd := exec.Command("go", "env", "GOPATH")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func add(a, b int) int {
	return a + b
}
func generateFile(templatePath, outputPath string, data TemplateData) error {
	funcMap := template.FuncMap{
		"add": add,
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	// Create directory if not exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

var generateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate code for model, repository, service, interface, or handler",
	Long: `Generate boilerplate code for different components of your application,
including models, repositories, services, interfaces, and handlers.`,
	Args: cobra.MinimumNArgs(2),
}

func generatorExecute(packagePath string) {
	generateCmd.Run = func(cmd *cobra.Command, args []string) {
		typeName := args[0]
		name := args[1]
		var lowerName = strings.ToLower(name)

		data := TemplateData{
			Name:        name,
			TableName:   fmt.Sprintf("tb_%s", lowerName),
			PackagePath: packagePath,
			LowerName:   lowerName,
		}

		templatePath := []string{
			filepath.Join("templates", "routes", "root.tmpl"),
		}
		outputPath := []string{
			filepath.Join("routes", "root.go"),
		}
		gopath, err := getGOPATH()
		if err != nil {
			fmt.Println("Error getting GOPATH:", err)
			return
		}

		if gopath == "" {
			fmt.Println("GOPATH is not set")
			return
		}

		modelPath := filepath.Join("templates", "models", "model.tmpl")
		repoPath := filepath.Join("templates", "repository", "repository.tmpl")
		servicePath := filepath.Join("templates", "service", "service.tmpl")
		interfacePath := filepath.Join("templates", "interface", "interface.tmpl")
		handlerPath := filepath.Join("templates", "handler", "handler.tmpl")
		routesPath := filepath.Join("templates", "routes", "routes.tmpl")
		sqlPath := filepath.Join("templates", "sql", "table.tmpl")

		modelOutput := filepath.Join("models", fmt.Sprintf("%s.go", lowerName))
		repoOutput := filepath.Join("repository", fmt.Sprintf("%s_repository.go", lowerName))
		serviceOutput := filepath.Join("service", fmt.Sprintf("%s_service.go", lowerName))
		interfaceOutput := filepath.Join("interfaces", fmt.Sprintf("%s_interface.go", lowerName))
		handlerOutput := filepath.Join("handlers", fmt.Sprintf("%s_handler.go", lowerName))
		routesOutput := filepath.Join("routes", fmt.Sprintf("%s_routes.go", lowerName))
		sqlOutput := filepath.Join("sql", fmt.Sprintf("%s_table.sql", lowerName))

		fields := []ModelField{}

		for _, arg := range os.Args[4:] {
			fieldParts := strings.Split(arg, ":")
			if len(fieldParts) < 2 {
				log.Fatalf("Invalid field format. Expected FieldName:FieldType,Validate")
			}

			name := fieldParts[0]
			typeAndValidation := strings.Split(fieldParts[1], ",")
			if len(typeAndValidation) == 0 {
				log.Fatalf("Invalid field format. Expected FieldType,Validate")
			}

			fieldType := typeAndValidation[0]
			var validate string
			if len(typeAndValidation) == 2 {
				validate = typeAndValidation[1]
			}

			fields = append(fields, ModelField{
				Name:     name,
				Type:     fieldType,
				JsonTag:  strings.ToLower(name),
				Validate: validate,
			})
		}

		data.ModelField = fields
		switch typeName {
		case "model":
			templatePath = append(templatePath, modelPath)
			outputPath = append(outputPath, modelOutput)
		case "repository":
			templatePath = append(templatePath, repoPath)
			outputPath = append(outputPath, repoOutput)
		case "service":
			templatePath = append(templatePath, servicePath)
			outputPath = append(outputPath, serviceOutput)
		case "interface":
			templatePath = append(templatePath, interfacePath)
			outputPath = append(outputPath, interfaceOutput)
		case "handler":
			templatePath = append(templatePath, handlerPath)
			outputPath = append(outputPath, handlerOutput)
		case "router":
			templatePath = append(templatePath, routesPath)
			outputPath = append(outputPath, routesOutput)
		case "sql":
			templatePath = append(templatePath, sqlPath)
			outputPath = append(outputPath, sqlOutput)
			data.SQLField = generateSQLFields(data.ModelField)
		case "api":
			templatePath = append(templatePath, modelPath, repoPath, servicePath, interfacePath, handlerPath, routesPath, sqlPath)
			outputPath = append(outputPath, modelOutput, repoOutput, serviceOutput, interfaceOutput, handlerOutput, routesOutput, sqlOutput)
			data.SQLField = generateSQLFields(data.ModelField)
		default:
			fmt.Println("Invalid type. Must be one of: model, repository, service, interface, handler")
			return
		}

		for i := 0; i < len(templatePath); i++ {
			fmt.Println("output", outputPath[i])
			if err := generateFile(templatePath[i], outputPath[i], data); err != nil {
				fmt.Println("Error generating file:", err)
			} else {
				fmt.Printf("Successfully generated %s at %s\n", typeName, outputPath[i])
			}
		}

	}

	rootCmd.AddCommand(generateCmd)
}

func goTypeToSQLType(goType string) string {
	switch goType {
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		return "INTEGER"
	case "float32", "float64":
		return "REAL"
	case "string":
		return "TEXT"
	case "time.Time":
		return "TIMESTAMP"
	default:
		return "TEXT" // Default to TEXT if type is unknown
	}
}

func parseValidationTags(tags string) string {
	constraints := []string{}
	for _, tag := range strings.Split(tags, ",") {
		switch tag {
		case "required":
			constraints = append(constraints, "NOT NULL")
		case "email":
			// You can add custom constraints or leave it for application logic
		}
	}
	return strings.Join(constraints, " ")
}

func generateSQLFields(modelFields []ModelField) []SQLField {
	var sqlFields []SQLField
	for _, field := range modelFields {
		sqlFields = append(sqlFields, SQLField{
			Name:        strings.ToLower(field.Name),
			Type:        goTypeToSQLType(field.Type),
			Constraints: parseValidationTags(field.Validate),
		})
	}
	return sqlFields
}
