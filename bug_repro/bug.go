package main

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

/*
 * multi line comment
 */
func main() {
	code := `package main
/*
 * multi line comment
 */
func main() {
}
`
	fmt.Println("---input---")
	fmt.Println(code)
	fmt.Println("---processing---")
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "bug.go", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, cg := range node.Comments {
		for _, c := range cg.List {
			fmt.Printf("Found comment     %q\n", c.Text)
			if len(c.Text) > 10 {
				fmt.Printf("Splitting comment %q\n", c.Text)
				c.Text = c.Text[:10] + "\n * " + c.Text[10:]
				fmt.Printf("into ->           %q\n", c.Text)
			}
		}
	}
	fmt.Println("---result---")
	if err := format.Node(os.Stdout, fset, node); err != nil {
		panic(err)
	}
}
