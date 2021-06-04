package main

import (
	"flag"
	"fmt"
	"github.com/rolancia/flatbuffers-schema-generator/lib"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		startsWithCapital = flag.Bool("sc", false, "whether to start with capitalized (title)")
		withMaker         = flag.Bool("maker", false, "whether to create makers")
		namespace         = flag.String("ns", "Default", "namespace of fbs")
		dbPath            = ""
	)
	flag.Parse()
	dbPath = flag.Arg(0)

	result, err := lib.Generate(dbPath, *namespace, *startsWithCapital)
	if err != nil {
		panic(err)
	}

	_, targetName := filepath.Split(dbPath)
	ext := filepath.Ext(targetName)
	name := strings.Replace(targetName, ext, "", -1)
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	outPath := filepath.Join(workDir, fmt.Sprintf("%s.%s", name, "fbs"))
	fmt.Println("out:", outPath)

	if err := save(outPath, result); err != nil {
		panic(err)
	}

	if *withMaker {
		makeOutPath := filepath.Join(workDir, fmt.Sprintf("%s.%s", name, "go"))
		genedMaker, err := lib.GenerateMaker(outPath)
		if err != nil {
			panic(err)
		}

		if err := save(makeOutPath, genedMaker); err != nil {
			panic(err)
		}
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
