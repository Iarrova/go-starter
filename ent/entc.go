//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	if err := entc.Generate("./schema", &gen.Config{
		Target:  "./generated",
		Package: "github.com/Iarrova/go-starter/ent/generated",
	}); err != nil {
		log.Fatal("running ent codegen:", err)
	}
}
