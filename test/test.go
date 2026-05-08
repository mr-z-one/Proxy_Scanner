package test

import (
	"fmt"
	"strings"
)

func main() {

	a := "application/json; charset=utf-8"

	fmt.Println(strings.TrimSpace(strings.Split(a, ";")[0]))

}
