package main

import (
	"fmt"
	"github.com/froozen/go-helpers"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	// Make sure the environment is fit for the application
	FfmpegExists()

	// Load the configuration first
	LoadConfig()

	// Register the HandleFuncs
	http.HandleFunc("/style", FileServeFunc(dataDir+"web/style.css"))
	http.HandleFunc("/", ServeFunc)

	// Listen and serve using the HandleFuncs
	err := http.ListenAndServe(addressString, nil)
	helpers.ErrorCheck(err, "serving via http")
}

// FfmpegExists checks wether ffmpeg is accessible
func FfmpegExists() {
	if exec.Command("ffmpeg", "-version").Run() != nil {
		helpers.Fail("Error: Can't access ffmpeg.")
	}
}

// FileServerFunc retruns a file serving handle func
func FileServeFunc(filename string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}

// ServeFunc is the server's main HandleFunc
func ServeFunc(w http.ResponseWriter, r *http.Request) {
	// Replace the %20 with spaces
	name := strings.Replace(r.RequestURI, "%20", " ", -1)
	// Remove the leading "/"
	name = strings.TrimPrefix(name, "/")

	// Simply serve the dir page for a dir
	if helpers.IsDir(root + name) {

		// This is needed to be able to move into subdirs
		// properly. We don't want it to be added to the
		// root, though.
		if name != "" && !strings.HasSuffix(name, "/") {
			name += "/"
		}

		ServeDirPage(w, r, name)
		return
	}

	// Serve video for qualifying file
	if FileQualifies(name) {
		ServeVideo(w, r, root+name)
		return
	}

	fmt.Fprint(w, "Invalid query")
}

// ServerErrPage serves an error page
func ServeErrPage(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Fprint(w, "An  error occured: ", err)
}
