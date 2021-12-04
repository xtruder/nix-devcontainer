package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func contains(s []interface{}, e string) bool {
	for _, a := range s {
		if a.(string) == e {
			return true
		}
	}
	return false
}

func envOrDefault(envName string, defaultVal string) string {
	if val := os.Getenv(envName); val != "" {
		return val
	}

	return defaultVal
}

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger

	verbose           bool
	preloadExtensionsStr string
	extensionsPath string
)

func main() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	vscodeDir := path.Join(home, ".vscode-server")
	if (os.Getenv("CODESPACES") == "true") {
		vscodeDir = path.Join(home, ".vscode-remote")
	}

	defaultExtensionsPath:= envOrDefault("VSCODE_EXTENSIONS_PATH", path.Join(vscodeDir, "extensions"))
	defaultPreloadExtensions := envOrDefault("PRELOAD_EXTENSIONS", "")

	flag.BoolVar(&verbose, "verbose", false,
		"Enable verbose logging")
	flag.StringVar(&extensionsPath, "path", defaultExtensionsPath,
		"VSCode extensions path")
	flag.StringVar(&preloadExtensionsStr, "ext", defaultPreloadExtensions,
		"Comma separated list of extensions to preload")

	flag.Parse()

	logFile, _ := os.OpenFile("/tmp/ext-preloader.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	logWriter := io.MultiWriter(os.Stdout, logFile)

	debugLogger = log.New(logWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(logWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(logWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	if !verbose {
		debugLogger.SetOutput(ioutil.Discard)
	}

	if preloadExtensionsStr == "" {
		errorLogger.Fatalln("missing preload extensions")
	}

	preloadExtensions := strings.Split(preloadExtensionsStr, ",")

	infoLogger.Printf("starting extension preloader, extensions to preload: %s", preloadExtensions)

	// make sure extensions path exists
	if err := os.MkdirAll(extensionsPath, 0755); err != nil {
		panic(err)
	}

	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	infoLogger.Printf("extensions path: %s", extensionsPath)

	infoLogger.Println("adding extension dir watcher")
	if err := watcher.Add(extensionsPath); err != nil {
		panic(err)
	}

	done := make(chan bool)
	go func() {
		files, err := ioutil.ReadDir(extensionsPath)
		if err != nil {
			errorLogger.Fatalf("error listing extensions: %v", err)
		}

		for _, f := range files {
			if !f.IsDir() {
				continue
			}

			extPath := path.Join(extensionsPath, f.Name())
			pkgJSONPath := path.Join(extPath, "package.json")

			if fixed, extName, err := modifyPackageJSON(pkgJSONPath, preloadExtensions); err != nil {
				errorLogger.Printf("error fixing package.json for: %s, %v", extPath, err)

				debugLogger.Printf("adding watcher for extension: %s", extPath)
				if err := watcher.Add(extPath); err != nil {
					debugLogger.Printf("error adding watcher for: %s, %v", extPath, err)
				}
			} else if fixed {
				infoLogger.Printf("modified package.json for extension '%s' or '%s'", extName, extPath)
			} else {
				debugLogger.Printf("extension already modified '%s'", extName)
			}
		}

		for {
			select {
			case event := <-watcher.Events:
				debugLogger.Printf("EVENT! %#v\n", event)

				if event.Op == fsnotify.Create && filepath.Dir(event.Name) == extensionsPath {
					debugLogger.Printf("adding watcher for extension path: %s", event.Name)

					if err := watcher.Add(event.Name); err != nil {
						debugLogger.Printf("error adding watcher for extension path: %s, %v", event.Name, err)
					}

					break
				}

				if event.Op == fsnotify.Create || event.Op == fsnotify.Write {
					filePath := event.Name

					pkgJSONPath := path.Join(extensionsPath, filepath.Base(filepath.Dir(filePath)), "package.json")

					if filePath == pkgJSONPath {
						debugLogger.Printf("modifing package.json: %s", pkgJSONPath)

						if fixed, extName, err := modifyPackageJSON(pkgJSONPath, preloadExtensions); err != nil {
							if !strings.Contains(err.Error(), "unexpected end of JSON input") {
								errorLogger.Printf("error modifing package.json for: %s, %v", pkgJSONPath, err)
							}
						} else if fixed {
							infoLogger.Printf("modified package.json for extension '%s' on path '%s'", extName, filePath)
						} else {
							debugLogger.Printf("extension already modified '%s'", extName)
						}
					}
				}
			case err := <-watcher.Errors:
				debugLogger.Println("ERROR", err)
			case <-sigs:
				done <- true
				return
			}
		}
	}()

	<-done
	infoLogger.Println("exiting")
}

// modifies package.json with added extensionsDependencies
func modifyPackageJSON(path string, preloadExtensions []string) (modified bool, extFullName string, err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return
	}

	extName, ok := result["name"].(string)
	if !ok {
		return false, "", fmt.Errorf("error parsing ext name")
	}

	extPublisher, ok := result["publisher"].(string)
	if !ok {
		return false, "", fmt.Errorf("error parsing ext publisher")
	}

	extFullName = extPublisher + "." + extName

	// do not modify itself
	for _, name := range preloadExtensions {
		if name == extFullName {
			return
		}
	}

	extensionDeps := []string{}
	extensionDeps = append(extensionDeps, preloadExtensions...)

	if val, ok := result["extensionDependencies"]; ok {
		if existingExtDeps, ok := val.([]interface{}); ok {
			for _, name := range preloadExtensions {
				if contains(existingExtDeps, name) {
					return
				}
			}

			for _, name := range existingExtDeps {
				extensionDeps = append(extensionDeps, name.(string))
			}
		}
	}

	result["extensionDependencies"] = extensionDeps

	byteValue, err = json.MarshalIndent(result, "", "  ")
	if err != nil {
		return
	}

	return true, extFullName, ioutil.WriteFile(path, byteValue, 0644)
}
