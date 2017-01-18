package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var (
	packageName string
	outputFile  string
	inputDir    string
	relativeTo  string
)

func init() {
	flag.StringVar(&packageName, "p", "", "Package name of generated source file")
	flag.StringVar(&outputFile, "o", "", "Filename for generated file")
	flag.StringVar(&inputDir, "i", "", "Input directory to compile")
	flag.StringVar(&relativeTo, "r", "", "Directory that names will be relative to")
}

func main() {
	flag.Parse()

	if packageName == "" {
		fmt.Println("Package name required")
		os.Exit(1)
	}
	if outputFile == "" {
		fmt.Println("Output file name required")
		os.Exit(1)
	}
	if inputDir == "" {
		fmt.Println("Input dir name required")
		os.Exit(1)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	if err := writeHeader(file); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := writeCommon(file); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := writeData(file); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func writeHeader(w io.Writer) error {
	_, err := w.Write([]byte(fmt.Sprintf(`
package %s

import (
	"fmt"
	"io/ioutil"
	"path"
)

var localMode = false
`, packageName)))
	return err
}

func writeCommon(w io.Writer) error {
	_, err := w.Write([]byte(`
func getBinData(name string) ([]byte, error) {
	name = path.Clean(name)

	if localMode {
		return getLocalData(name)
	}

	if _, ok := _binData[name]; !ok {
		return nil, fmt.Errorf("Static asset with name %s doesn't exist", name)
	}
	return _binData[name], nil
}

func getLocalData(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

`))
	return err
}

func writeData(w io.Writer) error {
	fmt.Fprintln(w, "var _binData = map[string][]byte{")

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := path.Join(inputDir, file.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		assetPath, _ := filepath.Rel(relativeTo, filename)
		fmt.Fprintf(w, `"%s": %#v,%s`, assetPath, data, "\n")
	}

	fmt.Fprintln(w, "}")
	return nil
}
