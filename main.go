package main

import (
	"fmt"
	"log"
	"round-trip/nested"
)

func main() {
	err := nested.ValidateYAMLMatchesStruct("nested/company.yaml", nested.Company{})
	if err != nil {
		log.Fatal("YAML validation failed:", err)
	}
	fmt.Println("YAML matches struct")
}

