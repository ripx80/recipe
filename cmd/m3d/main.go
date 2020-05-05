package main

import (
	"flag"

	"github.com/ripx80/recipe/pkgs/m3w"
)

func main() {
	outdir := flag.String("output", "recipes", "output dir. if not exists it will be created")
	flag.Parse()
	m3w.Down("https://www.maischemalzundmehr.de", *outdir)

}
