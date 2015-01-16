package main

import (
	"fmt"
	"github.com/froozen/go-helpers"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

// ServerDirPage serves a dir page
func ServeDirPage(w http.ResponseWriter, r *http.Request, dir string) {
	// Generate the code
	code, err := GenerateDirCode(dir)
	if err != nil {
		ServeErrPage(w, r, err)
		return
	}

	// Load the thempalte
	template, err := helpers.ReadFile(dataDir + "web/page.html")
	if err != nil {
		ServeErrPage(w, r, err)
		return
	}

	// Generate the data needed to generate the final page
	data := map[string]string{
		"CODE":    code,
		"DIRNAME": path.Base(dir),
	}

	// Serve the finished page
	fmt.Fprint(w, FillTemplate(template, data))
}

// GenerateDirCode generates the code for a dir page
func GenerateDirCode(dir string) (string, error) {
	// Get the files and dirs
	files, dirs, err := ListFilesAndDirs(dir)
	if err != nil {
		return "", err
	}

	// Generate the dirCode
	dirCode, err := GenerateItemCode(dirs, dataDir+"web/dir.html")
	if err != nil {
		return "", err
	}

	// Generate the fileCode
	files = FilterFiles(files)
	fileCode, err := GenerateItemCode(files, dataDir+"web/file.html")
	if err != nil {
		return "", err
	}

	// Join the two code pieces
	return dirCode + "\n" + fileCode, nil
}

// GenerateItemCode generates the code for a list of itmes and a template
func GenerateItemCode(items []map[string]string, templateName string) (string, error) {
	var code string

	// Load the template
	template, err := helpers.ReadFile(templateName)
	if err != nil {
		return "", err
	}

	// Generate the code
	for _, data := range items {
		// Fill the tempalte and append to code
		code += FillTemplate(template, data)
	}

	return code, nil
}

// FillTemplate fills a template with data
func FillTemplate(template string, data map[string]string) string {
	code := template

	// Replace all the placeholders with data
	for key, val := range data {
		placeholder := fmt.Sprintf("<!-- Data:%s -->", key)
		code = strings.Replace(code, placeholder, val, -1)
	}

	return code
}

// FilterFiles filters all the unqualified files out
func FilterFiles(files []map[string]string) (filteredFiles []map[string]string) {
	for _, file := range files {
		// Check wether the itemname qualifies
		if FileQualifies(file["ITEMNAME"]) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return
}

// ListFilesAndFolders returns the files and dirs of a dir in two separate slices
func ListFilesAndDirs(dir string) (files, dirs []map[string]string, err error) {
	// Retrieve the file info
	fileInfo, err := ioutil.ReadDir(root + dir)
	if err != nil {
		return nil, nil, err
	}

	// Add the .. dir except in the root
	if dir != "" {
		dirs = append(dirs, map[string]string{
			"ITEMNAME": "..",
			"PATHNAME": dir,
		})
	}

	// Seperate the data
	for _, item := range fileInfo {
		data := map[string]string{
			"ITEMNAME": path.Base(item.Name()),
			"PATHNAME": dir,
		}

		if item.IsDir() {
			dirs = append(dirs, data)
		} else {
			files = append(files, data)
		}
	}
	return
}
