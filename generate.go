package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

type ModelField struct {
	Name         string
	Type         string
	JsonTag      string
	Validate     string
	ValidateForm string
	ForeignKey   *ForeignKey
	AliasName    string
}

type ForeignKey struct {
	Table  string
	Column string
}
type FieldAlias struct {
	AliasTable  string
	AliasColumn string
	AliasName   string
}

type SQLField struct {
	Name        string
	Type        string
	Constraints string
	ForeignKey  string
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

	// Check if the file already exists
	if _, err := os.Stat(outputPath); err == nil {
		// File already exists, do not replace
		return nil
	} else if !os.IsNotExist(err) {
		// Error occurred while checking file existence
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

func generatorExecute(packagePath, templateName string) *cobra.Command {
	generateCmd.Run = func(cmd *cobra.Command, args []string) {
		typeName := args[0]
		name := args[1]
		lowerName := strings.ToLower(toSnakeCase(name))

		data := TemplateData{
			Name:        name,
			TableName:   fmt.Sprintf("tb_%s", lowerName),
			PackagePath: packagePath,
			LowerName:   lowerName,
		}

		templatePath := []string{
			filepath.Join(templateName, "app", "http", "handler", "auth_handler.tmpl"),
			filepath.Join(templateName, "app", "http", "handler", "profile_handler.tmpl"),
			filepath.Join(templateName, "app", "http", "handler", "menu_handler.tmpl"),
			filepath.Join(templateName, "app", "http", "handler", "user_handler.tmpl"),
			filepath.Join(templateName, "app", "http", "handler", "role_handler.tmpl"),

			filepath.Join(templateName, "app", "http", "middleware", "jwt_authentication.tmpl"),
			filepath.Join(templateName, "app", "http", "middleware", "jwt_request.tmpl"),

			filepath.Join(templateName, "app", "http", "models", "auth.tmpl"),
			filepath.Join(templateName, "app", "http", "models", "base.tmpl"),
			filepath.Join(templateName, "app", "http", "models", "menu.tmpl"),
			filepath.Join(templateName, "app", "http", "models", "role.tmpl"),
			filepath.Join(templateName, "app", "http", "models", "user.tmpl"),

			filepath.Join(templateName, "app", "http", "service", "auth_service.tmpl"),
			filepath.Join(templateName, "app", "http", "service", "menu_service.tmpl"),
			filepath.Join(templateName, "app", "http", "service", "role_service.tmpl"),
			filepath.Join(templateName, "app", "http", "service", "user_service.tmpl"),

			filepath.Join(templateName, "app", "interface", "auth_interface.tmpl"),
			filepath.Join(templateName, "app", "interface", "profile_interface.tmpl"),
			filepath.Join(templateName, "app", "interface", "menu_interface.tmpl"),
			filepath.Join(templateName, "app", "interface", "role_interface.tmpl"),
			filepath.Join(templateName, "app", "interface", "user_interface.tmpl"),

			filepath.Join(templateName, "app", "repository", "menu_repository.tmpl"),
			filepath.Join(templateName, "app", "repository", "role_repository.tmpl"),
			filepath.Join(templateName, "app", "repository", "user_repository.tmpl"),

			filepath.Join(templateName, "bootstrap", "app.tmpl"),
			filepath.Join(templateName, "config", "config.tmpl"),
			filepath.Join(templateName, "public", "assets", "asset.png"),
			filepath.Join(templateName, "public", "css", "sample.css"),

			filepath.Join(templateName, "routes", "auth_routes.tmpl"),
			filepath.Join(templateName, "routes", "profile_routes.tmpl"),
			filepath.Join(templateName, "routes", "menu_routes.tmpl"),
			filepath.Join(templateName, "routes", "role_routes.tmpl"),
			filepath.Join(templateName, "routes", "user_routes.tmpl"),
			filepath.Join(templateName, "routes", "root.tmpl"),

			filepath.Join(templateName, "utils", "validation.tmpl"),
			filepath.Join(templateName, "utils", "hash_password.tmpl"),

			filepath.Join(templateName, "pkg", "pgx.tmpl"),
			filepath.Join(templateName, "pkg", "redis.tmpl"),

			filepath.Join(templateName, "database", "migrations", "base.tmpl"),
		}

		outputPath := []string{
			filepath.Join("app", "http", "handlers", "auth_handler.go"),
			filepath.Join("app", "http", "handlers", "profile_handler.go"),
			filepath.Join("app", "http", "handlers", "menu_handler.go"),
			filepath.Join("app", "http", "handlers", "user_handler.go"),
			filepath.Join("app", "http", "handlers", "role_handler.go"),

			filepath.Join("app", "http", "middleware", "jwt_authentication.go"),
			filepath.Join("app", "http", "middleware", "jwt_request.go"),

			filepath.Join("app", "http", "models", "auth_model.go"),
			filepath.Join("app", "http", "models", "base_model.go"),
			filepath.Join("app", "http", "models", "menu_model.go"),
			filepath.Join("app", "http", "models", "role_model.go"),
			filepath.Join("app", "http", "models", "user_model.go"),

			filepath.Join("app", "http", "service", "auth_service.go"),
			filepath.Join("app", "http", "service", "menu_service.go"),
			filepath.Join("app", "http", "service", "role_service.go"),
			filepath.Join("app", "http", "service", "user_service.go"),

			filepath.Join("app", "interfaces", "auth_interface.go"),
			filepath.Join("app", "interfaces", "profile_interface.go"),
			filepath.Join("app", "interfaces", "menu_interface.go"),
			filepath.Join("app", "interfaces", "role_interface.go"),
			filepath.Join("app", "interfaces", "user_interface.go"),

			filepath.Join("app", "repository", "menu_repository.go"),
			filepath.Join("app", "repository", "role_repository.go"),
			filepath.Join("app", "repository", "user_repository.go"),

			filepath.Join("bootstrap", "app.go"),
			filepath.Join("config", "config.go"),
			filepath.Join("public", "assets", "asset.png"),
			filepath.Join("public", "css", "sample.css"),

			filepath.Join("routes", "auth_routes.go"),
			filepath.Join("routes", "profile_routes.go"),
			filepath.Join("routes", "menu_routes.go"),
			filepath.Join("routes", "role_routes.go"),
			filepath.Join("routes", "user_routes.go"),
			filepath.Join("routes", "root.go"),

			filepath.Join("utils", "validation.go"),
			filepath.Join("utils", "hash_password.go"),

			filepath.Join("pkg", "pgx.go"),
			filepath.Join("pkg", "redis.go"),

			filepath.Join("database", "migrations", "core_sql.sql"),
		}

		modelPath := filepath.Join(templateName, "app", "http", "models", "model.tmpl")
		repoPath := filepath.Join(templateName, "app", "repository", "repository.tmpl")
		servicePath := filepath.Join(templateName, "app", "http", "service", "service.tmpl")
		interfacePath := filepath.Join(templateName, "app", "interface", "interface.tmpl")
		handlerPath := filepath.Join(templateName, "app", "http", "handler", "handler.tmpl")
		routesPath := filepath.Join(templateName, "routes", "routes.tmpl")
		sqlPath := filepath.Join(templateName, "database", "migrations", "table.tmpl")

		modelOutput := filepath.Join("app/http/models", fmt.Sprintf("%s.go", lowerName))
		repoOutput := filepath.Join("app/repository", fmt.Sprintf("%s_repository.go", lowerName))
		serviceOutput := filepath.Join("app/http/service", fmt.Sprintf("%s_service.go", lowerName))
		interfaceOutput := filepath.Join("app/interfaces", fmt.Sprintf("%s_interface.go", lowerName))
		handlerOutput := filepath.Join("app/http/handlers", fmt.Sprintf("%s_handler.go", lowerName))
		routesOutput := filepath.Join("routes", fmt.Sprintf("%s_routes.go", lowerName))
		sqlOutput := filepath.Join("database/migrations", fmt.Sprintf("%s_table.sql", lowerName))

		fields := []ModelField{}

		if typeName != "boilerplate" {

			for _, arg := range os.Args[4:] {
				fieldParts := strings.Split(arg, ":")
				if len(fieldParts) == 0 {
					log.Fatalf("Invalid field format. Expected FieldName:FieldType,Validate")
				}

				name := fieldParts[0]
				var validate, fieldType, TableName, ColumnName string

				typeAndValidation := strings.Split(fieldParts[1], ",")
				if len(typeAndValidation) == 0 {
					validate = "string"
				}

				var aliasName []FieldAlias
				var validateList []string
				for idx, v := range typeAndValidation {
					if idx == 0 {
						fieldType = v
					} else {
						partsAlias := strings.Split(v, "--")
						frgnkey := strings.Split(v, "fk=")
						if len(partsAlias) == 2 {
							tbRef := strings.Split(partsAlias[1], "=")
							alsTb := tbRef[0]
							alsTBValue := tbRef[1]
							tbnm := strings.Split(alsTb, ".")
							aliasName = append(aliasName, FieldAlias{
								AliasTable:  tbnm[0],
								AliasColumn: tbnm[1],
								AliasName:   alsTBValue,
							})
						} else if len(frgnkey) == 2 {
							splt := strings.Split(frgnkey[1], ".")
							if len(splt) == 2 {
								TableName = splt[0]
								ColumnName = splt[1]
							}
						} else {
							validateList = append(validateList, v)
							validate = v
						}

					}

				}

				if TableName != "" && ColumnName != "" {
					fields = append(fields, ModelField{
						Name:         toCamelCase(name),
						Type:         fieldType,
						JsonTag:      toSnakeCase(name),
						Validate:     validate,
						ValidateForm: strings.Join(validateList, ","),
						ForeignKey: &ForeignKey{
							Table:  TableName,
							Column: ColumnName,
						},
					})
				} else {
					fields = append(fields, ModelField{
						Name:         toCamelCase(name),
						Type:         fieldType,
						JsonTag:      toSnakeCase(name),
						Validate:     validate,
						ValidateForm: strings.Join(validateList, ","),
					})
				}

				if len(aliasName) > 0 {
					for _, v := range aliasName {
						fields = append(fields, ModelField{
							Name:      toCamelCase(v.AliasName),
							Type:      "string",
							JsonTag:   toSnakeCase(v.AliasName),
							AliasName: fmt.Sprintf("%s.%s", v.AliasTable, v.AliasColumn),
						})
					}
				}

			}
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
		case "boilerplate":
			templatePath = append(templatePath, modelPath)
			outputPath = append(outputPath, modelOutput)
			data.SQLField = generateSQLFields(data.ModelField)
		default:
			log.Fatalf("Invalid type specified: %s", typeName)
		}

		data.Add = add

		for i := range templatePath {
			if err := generateFile(templatePath[i], outputPath[i], data); err != nil {
				log.Fatalf("Failed to generate file: %v", err)
			}
		}
	}
	return generateCmd
}

func generateSQLFields(modelFields []ModelField) []SQLField {
	var sqlFields []SQLField
	for _, field := range modelFields {
		var sqlType string
		switch field.Type {
		case "string":
			sqlType = "TEXT"
		case "int":
			sqlType = "INT"
		case "bool":
			sqlType = "BOOLEAN"
		case "interface{}":
			sqlType = "JSONB"
		default:
			sqlType = "TEXT"
		}

		switch field.Validate {
		case "required":
			field.Validate = "NOT NULL"
		case "unique":
			field.Validate = "UNIQUE"
		case "email":
			field.Validate = "UNIQUE"
		}

		if field.ForeignKey != nil {
			sqlFields = append(sqlFields, SQLField{
				Name:        toSnakeCase(field.Name),
				Type:        field.Type,
				Constraints: field.Validate,
				ForeignKey:  fmt.Sprintf("REFERENCES %s(id)", field.ForeignKey.Table),
			})
		} else {
			sqlFields = append(sqlFields, SQLField{
				Name:        toSnakeCase(field.Name),
				Type:        sqlType,
				Constraints: field.Validate,
			})
		}
	}
	return sqlFields
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(rune(s[i-1])) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1])))) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// Function to convert snake_case to CamelCase
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}
