package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	allLintersFile, err := os.Open("allLinters.txt")
	if err != nil {
		log.Println(err)
	}

	allLintersString, err := io.ReadAll(allLintersFile)
	if err != nil {
		log.Println(err)
	}

	deprecatedLinters := []string{
		"scopelint", "maligned", "golint", "interfacer", "exhaustivestruct",
	}
	goOneEighteenFailing := []string{
		"bodyclose", "contextcheck", "interfacer", "nilerr", "noctx",
		"rowserrcheck", "sqlclosecheck", "structcheck", "tparallel", "unparam",
		"wastedassign",
	}

	strLinters := strings.Split(string(allLintersString), "\n")

	linters := make(map[string]interface{})

	// the following is made to ensure that all linters are unique.
	for _, strLinter := range strLinters {
		if strings.TrimSpace(strLinter) == "" {
			continue
		}

		linters[strLinter] = nil
	}

	goodLinters := make([]string, 0, len(linters))

	for linter := range linters {
		if len(deprecatedLinters) > 0 {
			if isIn, pos := stringIsOneOf(linter, deprecatedLinters); isIn {
				deprecatedLinters = removeElement(deprecatedLinters, pos)

				continue
			}
		}

		if len(goOneEighteenFailing) > 0 {
			if isIn, pos := stringIsOneOf(linter, goOneEighteenFailing); isIn {
				goOneEighteenFailing = removeElement(goOneEighteenFailing, pos)

				continue
			}
		}

		goodLinters = append(goodLinters, linter)
	}

	err = printLinters(goodLinters)
	if err != nil {
		log.Println(err)
	}
}

// printLinters writes good linters for use in CL, vim and yaml config file.
func printLinters(linters []string) error {
	f, err := os.Create("configurations.txt")
	if err != nil {
		return fmt.Errorf("file creation: %w", err)
	}

	writeConfigName("cl run flags", f)
	f.WriteString("\n")
	writeConfiguration(linters, "golangci-lint run --disable-all ", "", "-E %s", " ", f)

	return nil
}

// writeConfiguration prints linter configuration to os.Writer
// with specified format.
func writeConfiguration(linters []string, header, footer, style, sep string, f *os.File) error {
	_, err := f.WriteString(header)
	if err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	for i := 0; i < len(linters); i++ {
		_, err = f.WriteString(fmt.Sprintf(style, linters[i]))
		if err != nil {
			return fmt.Errorf("writing linter: %w", err)
		}

		if i != len(linters)-1 {
			_, err = f.WriteString(sep)
			if err != nil {
				return fmt.Errorf("writing separator: %w", err)
			}
		}
	}

	_, err = f.WriteString(footer)
	if err != nil {
		return fmt.Errorf("writing footer: %w", err)
	}

	return nil
}

// writeConfigName prints beautified version of configuration to the file.
func writeConfigName(configName string, f *os.File) error {
	var err error

	width := 80

	for i := 0; i < width; i++ {
		_, err = f.WriteString("_")
		if err != nil {
			return fmt.Errorf("writing underlines: %w", err)
		}
	}

	_, err = f.WriteString("\n" + configName + "\n")
	if err != nil {
		return fmt.Errorf("writing config name: %w", err)
	}

	for i := 0; i < width; i++ {
		_, err = f.WriteString("‾")
		if err != nil {
			return fmt.Errorf("writing overlines: %w", err)
		}
	}

	return nil
}

// stringIsOneOf checks if the goal string is in the string slice.
func stringIsOneOf(goal string, options []string) (bool, int) {
	for id, opt := range options {
		if opt == goal {
			return true, id
		}
	}

	return false, 0
}

// removeElement returns string slice without element on specified position.
func removeElement(slice []string, pos int) []string {
	return append(slice[:pos], slice[pos+1:]...)
}
