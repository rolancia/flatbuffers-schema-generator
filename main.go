package main

import (
	"flag"
	"fmt"
	"github.com/rolancia/flatbuffers-schema-generator/lib"
	"os"
	"path"
	"strings"
)

func main() {
	args := os.Args[1:]

	var (
		dbPath    = args[0]
		namespace = ""
	)

	flag.StringVar(&namespace, "ns", "Default", "namespace")

	result, err := lib.Generate(dbPath, namespace)
	if err != nil {
		panic(err)
	}

	targetDir, targetName := path.Split(dbPath)
	ext := path.Ext(targetName)
	name := strings.Replace(targetName, ext, "", -1)
	outPath := path.Join(targetDir, fmt.Sprintf("%s.%s", name, "fbs"))

	if err := save(outPath, result); err != nil {
		panic(err)
	}
}

func save(name string, data string) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
