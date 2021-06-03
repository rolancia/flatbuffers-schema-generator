package lib

import "testing"

func TestGenerate(t *testing.T) {
	if result, err := Generate("nine.db", "Test"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(result)
	}
}
