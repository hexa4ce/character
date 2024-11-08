package character_test

import (
	"fmt"
	"github.com/hexa4ce/character"
	"os"
)

func Example() {
	// Read character card PNG file
	data, err := os.ReadFile("testdata/example.png")
	if err != nil {
		panic(err)
	}

	// Parse character data
	char, err := character.FromFile(data)
	if err != nil {
		panic(err)
	}

	// Access character information
	fmt.Printf("Name: %s\n", char.Name())
	fmt.Printf("Has description: %v\n", char.Description() != "")
	fmt.Printf("Has avatar: %v\n", char.Avatar() != "")
}
