package remote

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/usi-lfkeitel/saturn/src/utils"
)

func GenerateScript(config *utils.Config, modules []string) (string, error) {
	if !utils.FileExists(config.Core.TempDir) {
		if err := os.MkdirAll(config.Core.TempDir, 0755); err != nil {
			return "", err
		}
	}

	tempFile, err := ioutil.TempFile(config.Core.TempDir, "")
	if err != nil {
		return "", err
	}
	tempFileName := tempFile.Name()

	if config.Core.Debug {
		log.Printf("Generated Script: %s", tempFileName)
	}

	if err := generateRemoteScript(tempFile, config.Core.ModuleDir, modules); err != nil {
		tempFile.Close()
		return tempFileName, err
	}
	tempFile.Close()

	if err := os.Chmod(tempFileName, 0755); err != nil {
		return tempFileName, err
	}

	return tempFileName, nil
}

func generateRemoteScript(file *os.File, modulesDir string, modules []string) error {
	file.WriteString("#!/bin/bash\n\n")

	file.WriteString("MODULES=")
	file.WriteString(`(` + strings.Join(modules, " ") + `)`)

	file.WriteString(`
main() {
	echo -n '{'

	i=1
	for var in "${MODULES[@]}"; do
		echo -n "\"$var\": "
		echo -n $($var)
		if [ $i -lt ${#MODULES[@]} ]; then
				i=$[i + 1]
				echo -n ', '
		fi
	done

	echo -n '}'
}

`)

	goodModules := make(map[string]bool)

	for _, module := range modules {
		moduleFile := filepath.Join(modulesDir, module+".sh")

		m, err := getBinData(moduleFile)
		if err != nil {
			return fmt.Errorf("Module not found %s", module)
		}

		goodModules[module] = true

		file.WriteString(module + "() {\n")
		file.Write(m)
		file.WriteString("\n}\n\n")
	}

	file.WriteString(`main

if [ "$1" = "-d" ]; then
	rm "$0"
fi
`)

	return nil
}

func GetModuleList() []string {
	m := getFileList()
	for i, v := range m {
		m[i] = strings.Split(v, "/")[1]
		m[i] = m[i][:len(m[i])-3]
	}
	return m
}
