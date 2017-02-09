package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	flag.StringVar(&inputDir, "i", "", "Input directory to process")
}

type module struct {
	jsonType string
	keys     []*moduleKey
	types    map[string][]*moduleKey
	name     string
}

type moduleKey struct {
	name    string
	keyType string
	suffix  string
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

	modules, err := getModules()
	if err != nil {
		fmt.Println(err.Error())
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

	if err := writeModules(file, modules); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := writeControl(file, modules); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getModules() ([]*module, error) {
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}

	var modules []*module

	for _, file := range files {
		filename := path.Join(inputDir, file.Name())
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		_, _, err = reader.ReadLine()
		if err != nil {
			f.Close()
			return nil, err
		}

		genLine, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		var mod *module

		if bytes.Equal(genLine[:12], []byte("#gen:module2")) {
			mod, err = parseNewFormat(genLine, reader, filename, file)
		} else if bytes.Equal(genLine[:11], []byte("#gen:module")) {
			mod, err = parseOldFormat(genLine, filename, file)
		} else {
			log.Printf("Skipping file %s", file.Name())
			continue
		}

		if err != nil {
			return nil, err
		}

		modules = append(modules, mod)
		f.Close()
	}
	return modules, nil
}

func parseNewFormat(firstline []byte, reader *bufio.Reader, filename string, file os.FileInfo) (*module, error) {
	genHeaderLine := bytes.SplitN(firstline, []byte(" "), 2)
	if len(genHeaderLine) != 2 {
		return nil, fmt.Errorf("Invalid gen line in file %s", filename)
	}

	mod := &module{
		name:     file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))],
		jsonType: string(genHeaderLine[1]),
		keys:     make([]*moduleKey, 0),
		types:    make(map[string][]*moduleKey),
	}

	mode := "key"
	t := ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line = bytes.TrimSpace(line)
		line = bytes.TrimLeft(line, "#")
		line = bytes.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if bytes.Equal(line[:12], []byte("!gen:module2")) {
			if mode != "key" {
				return nil, fmt.Errorf("Missing 'endtype' in module %s", mod.name)
			}
			break
		}

		lineParts := bytes.Split(line, []byte{' '})
		if mode == "key" { // root key mode
			if string(lineParts[0]) == "type" {
				if len(lineParts) != 2 {
					return nil, fmt.Errorf("Missing type name in module %s", mod.name)
				}
				t = string(lineParts[1])
				if _, exists := mod.types[t]; exists {
					return nil, fmt.Errorf("Duplicate type definition in module %s", mod.name)
				}
				mod.types[t] = make([]*moduleKey, 0, 1)
				mode = "type"
				continue
			} else if string(lineParts[0]) == "key" {
				if len(lineParts) < 3 {
					return nil, fmt.Errorf("Invalid key in module %s", mod.name)
				}
				key := &moduleKey{
					name:    string(lineParts[1]),
					keyType: string(lineParts[2]),
				}
				if len(lineParts) == 4 {
					key.suffix = string(lineParts[3])
				}
				mod.keys = append(mod.keys, key)
				continue
			} else {
				return nil, fmt.Errorf("Unknown symbol %s in module %s", lineParts[0], mod.name)
			}
		} else { // type mode
			if bytes.Equal(lineParts[0], []byte("endtype")) {
				t = ""
				mode = "key"
				continue
			}

			if len(lineParts) < 2 {
				return nil, fmt.Errorf("Invalid key in type %s in module %s", t, mod.name)
			}
			key := &moduleKey{
				name:    string(lineParts[0]),
				keyType: string(lineParts[1]),
			}
			if len(lineParts) == 3 {
				key.suffix = string(lineParts[2])
			}
			mod.types[t] = append(mod.types[t], key)
			continue
		}
	}
	return mod, nil
}

func parseOldFormat(line []byte, filename string, file os.FileInfo) (*module, error) {
	genLineParts := bytes.SplitN(line, []byte(" "), 3)
	if len(genLineParts) != 3 {
		return nil, fmt.Errorf("Invalid gen line in file %s", filename)
	}

	mod := &module{
		name:     file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))],
		jsonType: string(genLineParts[1]),
	}

	keys := bytes.Split(genLineParts[2], []byte(","))
	mod.keys = make([]*moduleKey, 0, len(keys))

	for _, key := range keys {
		keyParts := strings.Split(string(key), ":")
		newKey := &moduleKey{
			name:    keyParts[0],
			keyType: keyParts[1],
		}

		if len(keyParts) == 3 {
			newKey.suffix = keyParts[2]
		}

		mod.keys = append(mod.keys, newKey)
	}
	return mod, nil
}

func writeHeader(w io.Writer) error {
	_, err := w.Write([]byte(fmt.Sprintf(`package %s

// DO NOT EDIT. This file is generated by cmd/generateModuleTypes/main.go.

import (
	"fmt"
	"strings"
)

`, packageName)))
	return err
}

func writeModules(w io.Writer, modules []*module) error {
	for _, mod := range modules {
		modName := goName(mod.name)
		fmt.Fprintf(w, "type %s struct {\n", modName)
		for _, key := range mod.keys {
			ktype := key.keyType
			if ktype[:2] == "[]" {
				ktype = "[]*" + modName + goName(ktype[2:])
			}
			fmt.Fprintf(w, "	%s %s `json:\"%s\"`\n", goName(key.name), ktype, key.name)
		}
		fmt.Fprint(w, "}\n\n")

		for t, keys := range mod.types {
			modName := goName(mod.name)
			fmt.Fprintf(w, "type %s struct {\n", modName+goName(t))
			for _, key := range keys {
				ktype := key.keyType
				if ktype[:2] == "[]" {
					ktype = "[]*" + modName + goName(ktype[2:])
				}
				fmt.Fprintf(w, "	%s %s `json:\"%s\"`\n", goName(key.name), ktype, key.name)
			}
			fmt.Fprint(w, "}\n\n")
		}

		if mod.jsonType != "a" && mod.jsonType != "o" {
			return fmt.Errorf("Json type %s not supported", mod.jsonType)
		}

		writePrint(w, mod)
	}
	return nil
}

func writePrint(w io.Writer, mod *module) {
	srcName := goName(mod.name)
	fancyName := displayName(mod.name)

	headerType := "a []"
	if mod.jsonType == "o" {
		headerType = "o "
	}

	fmt.Fprintf(w, `func printLong%s(depth int, %s*%s) {
	indent := strings.Repeat(" ", depth*2)
	fmt.Printf("%%s%s:\n", indent)
`, srcName, headerType, srcName, fancyName)

	if mod.jsonType == "a" {
		fmt.Fprintf(w, "	for _, o := range a {\n")
	}

	for _, key := range mod.keys {
		if key.suffix == "%" {
			key.suffix = "%%"
		}

		ktype := key.keyType
		isSlice := false
		if ktype[:2] == "[]" {
			ktype = ktype[2:]
			isSlice = true
		}
		if _, exists := mod.types[ktype]; exists {
			fmt.Fprintf(w, "	fmt.Printf(\"%%s%s: \\n\", indent)\n", displayName(key.name))
			if isSlice {
				fmt.Fprintf(w, "	for _, p := range o.%s {\n", goName(key.name))
				fmt.Fprintf(w, "		printLong%s(depth+1, p)\n", srcName+goName(ktype))
				fmt.Fprintf(w, "	}\n")
			} else {
				fmt.Fprintf(w, "	printLong%s(depth+1, o.%s)\n", srcName+goName(ktype), goName(key.name))
			}
		} else {
			fmt.Fprintf(w, "	fmt.Printf(\"%%s%s: %s%s\\n\", indent, o.%s)\n",
				displayName(key.name),
				fmtType(ktype),
				key.suffix,
				goName(key.name),
			)
		}
	}
	if mod.jsonType == "a" {
		fmt.Fprint(w, "	fmt.Println(\"\")\n	}\n}\n\n")
	} else {
		fmt.Fprint(w, "	fmt.Println(\"\")\n}\n\n")
	}

	for t, keys := range mod.types {
		fullType := srcName + goName(t)
		fmt.Fprintf(w, `func printLong%s(depth int, o *%s) {
	indent := strings.Repeat(" ", depth*2)
`, fullType, fullType)

		for _, key := range keys {
			if key.suffix == "%" {
				key.suffix = "%%"
			}

			ktype := key.keyType
			if _, exists := mod.types[ktype]; exists {
				fmt.Fprintf(w, "	fmt.Printf(\"%%s%s: \\n\", indent)\n", displayName(key.name))
				fmt.Fprintf(w, "	printLong%s(depth+1, o.%s)\n", srcName+goName(ktype), goName(key.name))
			} else {
				fmt.Fprintf(w, "	fmt.Printf(\"%%s%s: %s%s\\n\", indent, o.%s)\n",
					displayName(key.name),
					fmtType(ktype),
					key.suffix,
					goName(key.name),
				)
			}
		}
	}

	if len(mod.types) > 0 {
		fmt.Fprint(w, "	fmt.Println(\"\")\n}\n\n")
	}
}

func writeControl(w io.Writer, modules []*module) error {
	fmt.Fprint(w, "type HostResponse struct {\n	Host *ConfigHost `json:\"host\"`\n")
	for _, mod := range modules {
		srcName := goName(mod.name)
		fmt.Fprintf(w, "	%s %s*%s `json:\"%s,omitempty\"`\n",
			srcName, srcType(mod.jsonType), srcName, mod.name,
		)
	}
	fmt.Fprint(w, `}

func (r *HostResponse) Print(short bool) {
	if r == nil {
		return
	}

	if short {
		r.printShort()
		return
	}
	r.printLong()
}

func (r *HostResponse) printShort() {
	// TODO
}

`)

	fmt.Fprint(w, "func (r *HostResponse) printLong() {\n")

	for _, mod := range modules {
		srcName := goName(mod.name)
		if mod.jsonType == "a" {
			fmt.Fprintf(w, `	if len(r.%s) > 0 {
		printLong%s(1, r.%s)
		fmt.Println("")
	}
`, srcName, srcName, srcName)
		} else if mod.jsonType == "o" {
			fmt.Fprintf(w, `	if r.%s != nil {
		printLong%s(1, r.%s)
		fmt.Println("")
	}
`, srcName, srcName, srcName)
		}
	}

	fmt.Fprint(w, "}\n")
	return nil
}

func srcType(t string) string {
	switch t {
	case "a":
		return "[]"
	}
	return ""
}

func fmtType(t string) string {
	switch t {
	case "string":
		return "%s"
	case "float64":
		return "%.2f"
	case "bool":
		return "%t"
	case "int":
		return "%d"
	}
	return ""
}

func displayName(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return s
}

func goName(s string) string {
	s = displayName(s)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	s = strings.Replace(s, "-", "", -1)
	return s
}
