package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	// os.Setenv("FOO", "1")
	// os.Setenv("ATLAS_URI_PALABRAS_EXPRESS_API", "mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority")
	fmt.Println("FOO:", os.Getenv("FOO"))
	fmt.Println("BAR:", os.Getenv("BAR"))
	fmt.Println("ATLAS_URI_PALABRAS_EXPRESS_API:", os.Getenv("ATLAS_URI_PALABRAS_EXPRESS_API"))

	fmt.Println()
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0])
	}
}
