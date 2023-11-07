/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type Dependency struct {
	GroupID    string
	ArtifactID string
	Packaging  string
	Version    string
	Scope      string
}

func extractDependencies(treeText string) []Dependency {
	dependencies := make([]Dependency, 0)
	lines := strings.Split(treeText, "\n")
	re := regexp.MustCompile(`^\s*\+-\s+(.*:.*:.*:.*:.*)`)

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			dependency := strings.Split(match[1], ":")
			if len(dependency) == 5 {
				dep := Dependency{
					GroupID:    dependency[0],
					ArtifactID: dependency[1],
					Packaging:  dependency[2],
					Version:    dependency[3],
					Scope:      dependency[4],
				}
				dependencies = append(dependencies, dep)
			}
		}
	}

	return dependencies
}
// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert Maven denpendencies txt file to xlsx format",
	Long: `You need to run 'mvn dependency:tree -DoutputFile=input.txt' first to generate the input.txt file. And then run 'maven-dependency2excel convert' to convert the input.txt file to xlsx format.`,
	Run: func(cmd *cobra.Command, args []string) {

		filePath := "input.txt"
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var mavenDependencyTree string
		for scanner.Scan() {
			mavenDependencyTree += scanner.Text() + "\n"
		}

		dependencies := extractDependencies(mavenDependencyTree)
		currentTime := time.Now()
		currentDatetime := currentTime.Format("2006-01-02_15-04")
		baseFileName := strings.TrimSuffix(filePath, ".txt")
		outputFileName := fmt.Sprintf("%s_dependencies_%s.xlsx", baseFileName, currentDatetime)

		xlsx := excelize.NewFile()
		sheetName := "Sheet1"
		index := xlsx.NewSheet(sheetName)
		xlsx.SetActiveSheet(index)

		titleStyle, _ := xlsx.NewStyle(`{"font":{"bold":true}}`)
		xlsx.SetCellStyle(sheetName, "A1", "E1", titleStyle)

		xlsx.SetCellValue(sheetName, "A1", "GroupID")
		xlsx.SetCellValue(sheetName, "B1", "ArtifactID")
		xlsx.SetCellValue(sheetName, "C1", "Packaging")
		xlsx.SetCellValue(sheetName, "D1", "Version")
		xlsx.SetCellValue(sheetName, "E1", "Scope")

		fmt.Println("convert called")
		for i, dep := range dependencies {
			rowIndex := i + 2
			xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), dep.GroupID)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), dep.ArtifactID)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), dep.Packaging)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), dep.Version)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), dep.Scope)
		}
	
		if err := xlsx.SaveAs(outputFileName); err != nil {
			fmt.Println("Error saving Excel file:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
