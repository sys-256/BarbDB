# BarbDB

## Features

-   A solid key-value database engine in under 90 source lines of code.
-   Compatibility with strings including control characters.
-   A (for computers) easy to read file format.
-   Written in Golang, resulting in great performance.

## Usage

First, install the library with this command:

```sh
go get github.com/sys-256/BarbDB
```

Then you can use BarbDB by simply importing it, for example:

```go
package main

import (
	"fmt"

	"github.com/sys-256/BarbDB"
)

func main() {
	// Initialize the database
	db, openError := BarbDB.OpenDB("./main.barb")
	if openError != nil {
		fmt.Println(openError)
		return
	}

	// Set some values
	setError1 := db.Set("Foo", "This is BarbDB!")
	setError2 := db.Set("Bar", "Even \t, \n and \000 work!")
	if setError1 != nil || setError2 != nil {
		fmt.Println("Error setting values!\n", setError1, setError2)
	}

	// Try retrieving one of the set values
	result, getError := db.Get("Bar")
	if getError != nil {
		fmt.Println(getError)
	} else {
		fmt.Println(result)
	}

	// Delete a value
	deleteError := db.Delete("Foo")
	if deleteError != nil {
		fmt.Println(deleteError)
	}

	// Close the database
	closeError := db.Close()
	if closeError != nil {
		fmt.Println(closeError)
	}
}

```

## Contribution

Everyone is free to contribute, whether that is a bug report, feature request or feature implementation! If you want to implement a big feature, or if you're going to implement breaking changes, please create an issue to ask for approval first.

## License

See [the LICENSE file](./LICENSE).
