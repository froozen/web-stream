package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/froozen/go-helpers"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	addressString, root    string
	delay                  int
	args, hooks, fileTypes []string

	dataDir = os.Getenv("HOME") + "/.web-stream/"
)

// Load config loads the configuration from the configuration file
func LoadConfig() {
	// Open the config file
	file, err := os.Open(dataDir + "config.json")
	if err != nil {
		helpers.Fail("Error:", dataDir+"config.json", "doesn't exist")
	}
	defer file.Close()

	// Parse the json
	json, err := simplejson.NewFromReader(file)
	if err != nil {
		helpers.Fail("Error: Parsing the configuration failed:", err)
	}

	// Load the "port" value and append it to ":" to get a valid address
	// for listening and serving
	addressString = ":" + strconv.Itoa(json.Get("port").MustInt(2223))

	// Get the "root" value
	root = json.Get("root").MustString(os.Getenv("HOME") + "/Videos/")
	if !helpers.IsDir(root) {
		helpers.Fail("Error:", root, "is no directory")
	}
	// Make sure root ends with a "/"
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}

	// Load the "delay" value
	delay = json.Get("delay").MustInt(3)

	// Load the "args" value
	args = StringSlice(json.Get("args"))

	// Load the "hooks" value
	hooks = StringSlice(json.Get("hooks"))

	// Load the "filetypes" value
	fileTypes = StringSlice(json.Get("filetypes"))
}

// StringSlice creates a string slice from a simplejson.Json value
func StringSlice(json *simplejson.Json) (slice []string) {
	// Extract the array
	array, err := json.Array()
	if err != nil {
		return
	}

	// Add all strings to the slice
	for _, item := range array {
		v := reflect.ValueOf(item)
		if v.Kind() == reflect.String {
			slice = append(slice, v.String())
		}
	}
	return
}
