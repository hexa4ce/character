# character

A Go package for parsing character cards in PNG format.

## Installation

```bash
go get github.com/hexa4ce/character
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/hexa4ce/character"
    "os"
)

func main() {
    // Read character card PNG file
    data, err := os.ReadFile("character.png")
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
    fmt.Printf("Description: %s\n", char.Description())
    fmt.Printf("Avatar (base64): %s\n", char.Avatar())
}
```

## License

MIT License - see LICENSE file for details
