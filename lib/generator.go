package lib

import (
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func Generate(dbPath string, namespace string, capital bool) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", err
	}

	j, err := ExportAsJson(db)
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	writeStringL(&b, newNamespace(namespace))

	sortedKeys := make([]string, 0, len(j))
	for k := range j {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i] < sortedKeys[j]
	})

	for i := range sortedKeys {
		tName := sortedKeys[i]
		tRow := j[tName]
		tys := inferType(tRow)
		writeStringL(&b, newStruct(tName, tys, capital))
	}

	return b.String(), nil
}

func writeStringL(b *strings.Builder, str string) {
	_, err := b.WriteString(fmt.Sprintf("%s\n", str))
	if err != nil {
		panic(err)
	}
}

func writeStringTL(b *strings.Builder, str string) {
	_, err := b.WriteString(fmt.Sprintf("\t%s\n", str))
	if err != nil {
		panic(err)
	}
}

func inferType(arr interface{}) []fbType {
	vof := reflect.ValueOf(arr).Index(0)
	keys := vof.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() > keys[j].String()
	})
	ret := make([]fbType, len(keys), len(keys))
	for i := range keys {
		k := keys[i].String()
		vt := vof.MapIndex(keys[i]).Elem().Kind()
		convertedTy := tyMap[vt]
		if convertedTy == "" {
			panic(fmt.Errorf("unknown type %s of %s", vt.String(), k))
		}

		ret[i] = fbType{
			name: k,
			ty:   convertedTy,
		}
	}

	return ret
}

func newNamespace(name string) string {
	return fmt.Sprintf("namespace %s;", name)
}

type fbType struct {
	name string
	ty   string
}

func newStruct(name string, tys []fbType, c bool) string {
	if c {
		name = strings.Title(name)
	}

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("table %s_e {\n", name))
	for i := range tys {
		b.WriteString(fmt.Sprintf("\t%s: %s;\n", tys[i].name, tys[i].ty))
	}
	b.WriteString(fmt.Sprintf("}\n"))
	return b.String()
}

var tyMap = map[reflect.Kind]string{
	reflect.String: "string",

	reflect.Bool: "bool",

	reflect.Int:    "int",
	reflect.Int8:   "int8",
	reflect.Int16:  "int16",
	reflect.Int32:  "int32",
	reflect.Int64:  "int64",
	reflect.Uint:   "uint",
	reflect.Uint8:  "uint8",
	reflect.Uint16: "uint16",
	reflect.Uint32: "uint32",
	reflect.Uint64: "uint64",

	reflect.Float32: "float32",
	reflect.Float64: "float64",
}
