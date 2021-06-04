package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GenerateMaker(fbsPath string) (string, error) {
	f, err := os.Open(fbsPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b := &strings.Builder{}
	writeStringL(b, importStr())
	writeStringL(b, i64Str())
	writeStringL(b, f64Str())
	writeStringL(b, strStr())

	//
	ns := ""

	//
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "namespace") {
			n := strings.Split(line, " ")[1]
			n = strings.ReplaceAll(n, ";", "")
			ns = n
		} else if strings.HasPrefix(line, "table") {
			tableName := strings.Split(line, " ")[1]
			nonScala := []fbType{}
			scala := []fbType{}

			for sc.Scan() {
				line := strings.ReplaceAll(sc.Text(), ";", "")
				vt := strings.Split(line, ":")
				if len(vt) < 2 {
					break
				}
				name := strings.TrimSpace(vt[0])
				ty := strings.TrimSpace(vt[1])
				fbTy := fbType{
					name: name,
					ty:   ty,
				}
				if ty == "string" {
					nonScala = append(nonScala, fbTy)
				} else {
					scala = append(scala, fbTy)
				}
			}

			b.WriteString(strMaker(ns, tableName, nonScala, scala))
		}
	}

	return b.String(), nil
}

//
func importStr() string {
	return `
import (
	flatbuffers "github.com/google/flatbuffers/go"
)
`
}

func i64Str() string {
	return `
func i64(v interface{}) int64 {
	vv, ok := v.(int64)
	if ok == false {
		vv = 0
	}
	return vv
}
`
}

func f64Str() string {
	return `
func f64(v interface{}) float64 {
	vv, ok := v.(float64)
	if ok == false {
		vv = 0.0
	}
	return vv
}
`
}

func strStr() string {
	return `
func str(b *flatbuffers.Builder, v interface{}) flatbuffers.UOffsetT {
	vv, ok := v.(string)
	if ok == false {
		vv = ""
	}
	p := b.CreateByteString([]byte(vv))
	return p
}
`
}

func strMaker(ns string, tableName string, nonScala []fbType, scala []fbType) string {
	tableName = strings.Title(tableName)

	b := strings.Builder{}
	writeStringL(&b, fmt.Sprintf("func Make%s (b *flatbuffers.Builder, j map[string]interface{}) flatbuffers.UOffsetT {", tableName))

	for i := range nonScala {
		switch nonScala[i].ty {
		case "string":
			writeStringTL(&b, fmt.Sprintf("%s := str(b, j[\"%s\"])", "__"+nonScala[i].name, nonScala[i].name))
		}
	}

	writeStringTL(&b, fmt.Sprintf("%s.%sStart(b)", ns, tableName))

	for i := range nonScala {
		switch nonScala[i].ty {
		case "string":
			writeStringTL(&b, fmt.Sprintf("%s.%sAdd%s(b, %s)", ns, tableName, camel(nonScala[i].name), "__"+nonScala[i].name))
		default:
			panic(fmt.Errorf("not supported type %s", nonScala[i].ty))
		}
	}

	for i := range scala {
		var fnName string
		switch scala[i].ty {
		case "int64":
			fnName = "i64"
		case "float64":
			fnName = "f64"
		default:
			panic(fmt.Errorf("not supported type %s", scala[i].ty))
		}

		writeStringTL(&b, fmt.Sprintf("%s.%sAdd%s(b, %s(j[\"%s\"]))", ns, tableName, camel(scala[i].name), fnName, scala[i].name))
	}

	writeStringTL(&b, fmt.Sprintf("return %s.%sEnd(b)\n}\n", ns, tableName))

	return b.String()
}

func camel(str string) string {
	spd := strings.Split(str, "_")
	for i := range spd {
		spd[i] = strings.Title(spd[i])
	}
	return strings.Join(spd, "")
}
