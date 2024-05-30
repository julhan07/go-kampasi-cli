package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type TemplateData struct {
	Name        string
	LowerName   string
	PackagePath string
}

func generateFile(templatePath, outputPath string, data TemplateData) error {
	tmpl, err := template.ParseFiles(templatePath)
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
	Run: func(cmd *cobra.Command, args []string) {
		typeName := args[0]
		name := args[1]
		packagePath := "api-generator"
		var lowerName = strings.ToLower(name)

		data := TemplateData{
			Name:        name,
			PackagePath: packagePath,
			LowerName:   lowerName,
		}

		var templatePath, outputPath string

		switch typeName {
		case "model":
			templatePath = "templates/models/model.tmpl"
			outputPath = filepath.Join("models", fmt.Sprintf("%s.go", lowerName))
		case "repository":
			templatePath = "templates/repository/repository.tmpl"
			outputPath = filepath.Join("repository", fmt.Sprintf("%s_repository.go", lowerName))
		case "service":
			templatePath = "templates/service/service.tmpl"
			outputPath = filepath.Join("service", fmt.Sprintf("%s_service.go", lowerName))
		case "interface":
			templatePath = "templates/interface/interface.tmpl"
			outputPath = filepath.Join("interfaces", fmt.Sprintf("%s_interface.go", lowerName))
		case "handler":
			templatePath = "templates/handler/handler.tmpl"
			outputPath = filepath.Join("handlers", fmt.Sprintf("%s_handler.go", lowerName))
		default:
			fmt.Println("Invalid type. Must be one of: model, repository, service, interface, handler")
			return
		}

		fmt.Println("output", outputPath)
		fmt.Println("data name", data.Name)

		if err := generateFile(templatePath, outputPath, data); err != nil {
			fmt.Println("Error generating file:", err)
		} else {
			fmt.Printf("Successfully generated %s at %s\n", typeName, outputPath)
		}
	},
}

func generatorExecute() {
	rootCmd.AddCommand(generateCmd)
}
