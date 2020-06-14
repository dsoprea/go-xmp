package xmp

import (
	"os"
	"path"

	"io/ioutil"

	"github.com/dsoprea/go-logging"
)

var (
	testDataRelFilepath = "test.xmp"
)

var (
	moduleRootPath = ""
	assetsPath     = ""
)

// GetModuleRootPath returns the root-path of the module.
func GetModuleRootPath() string {
	if moduleRootPath == "" {
		moduleRootPath = os.Getenv("XMP_MODULE_ROOT_PATH")
		if moduleRootPath != "" {
			return moduleRootPath
		}

		currentWd, err := os.Getwd()
		log.PanicIf(err)

		currentPath := currentWd
		visited := make([]string, 0)

		for {
			tryStampFilepath := path.Join(currentPath, ".MODULE_ROOT")

			_, err := os.Stat(tryStampFilepath)
			if err != nil && os.IsNotExist(err) != true {
				log.Panic(err)
			} else if err == nil {
				break
			}

			visited = append(visited, tryStampFilepath)

			currentPath = path.Dir(currentPath)
			if currentPath == "/" {
				log.Panicf("could not find module-root: %v", visited)
			}
		}

		moduleRootPath = currentPath
	}

	return moduleRootPath
}

// GetTestAssetsPath returns the path of the test-assets.
func GetTestAssetsPath() string {
	if assetsPath == "" {
		moduleRootPath := GetModuleRootPath()
		assetsPath = path.Join(moduleRootPath, "assets")
	}

	return assetsPath
}

// GetTestDataFilepath returns the file-path of the common test-data.
func GetTestDataFilepath() string {
	assetsPath := GetTestAssetsPath()
	filepath := path.Join(assetsPath, testDataRelFilepath)

	return filepath
}

// GetTestData returns the common test-data.
func GetTestData() []byte {
	filepath := GetTestDataFilepath()

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	return data
}
