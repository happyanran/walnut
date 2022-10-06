package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestPath(t *testing.T) {
	p := strings.TrimRight(filepath.Join("./data", "123", "*"), "*")
	fmt.Println(p)
	fmt.Println(len(p))

	fmt.Println("abc" + "%")
}
