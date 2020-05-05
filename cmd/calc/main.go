package main

import (
	"fmt"

	"github.com/ripx80/recipe"
)

func main() {
	rec, err := recipe.LoadFile("../../testdata/ipa.json", &recipe.M3{})
	if err != nil {
		fmt.Println(err)
		return
	}
	rScale, err := rec.Scale(30.0, 70.0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rScale)
}
