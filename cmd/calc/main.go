package main

import (
	"fmt"
	"os"

	"github.com/ripx80/recipe"
	"github.com/ripx80/recipe/pkgs/plato"
)

const in string = "../../testdata/ipa.json"

/*
	Skalieren bei Änderung
	[ ] Schüttung
	[ ] Ausschlagwürze
	[ ] Sudausbeute

	Schüttung: SumMalt

	Ausschlagwürze: DecisiveSeasoning AM
	Sudausbeute: 	SudYield

	add tablewriter for output plato table
*/

func main() {
	recipe, err := recipe.LoadFile(in, &recipe.M3{})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	Plato := plato.New()

	//fmt.Println(recipe)
	_, err = Plato.SW(recipe.Global.OriginalWort)
	//fmt.Println(e)

	newam := 25
	fmt.Printf("original AM:\t %f\n", recipe.Global.DecisiveSeasoning)
	fmt.Printf("new AM:\t\t %d\n", newam)

	// e, err = Plato.SG(1.06846)
	// fmt.Println(e)
	// e, err = Plato.AF(17.35)
	// fmt.Println(e)

}
