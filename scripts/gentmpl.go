package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	out, err := os.Create("templates.go")
	if err != nil {
		log.Fatal(err)
	}

	files, err := filepath.Glob("templates/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(out, "// Genrated by %s\n\n", os.Args[0])
	fmt.Fprintf(out, "package main\n\n")
	names := make([]string, len(files))
	for i, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		name := path.Base(file)
		name = name[:len(name)-len(path.Ext(file))]
		fmt.Fprintf(out, "var %s = `%s`", name, string(data))
		names[i] = name
	}

	fmt.Fprintf(out, "\n\nvar Templates = map[string]string{\n")
	for _, name := range names {
		fmt.Fprintf(out, "\t\"%s\": %s,\n", name, name)
	}
	fmt.Fprintf(out, "}")
}
