package main

import (
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"fortio.org/cli"
	"fortio.org/log"
)

func main() {
	maxlen := flag.Int("len", 79, "max line length") // 78 or 80 or 132 works but... somehow 79 is creating 2 new lines(!)
	cli.MinArgs = 1
	cli.MaxArgs = -1
	cli.ArgsHelp = "filenames..."
	cli.Main()
	for _, filename := range flag.Args() {
		process(filename, *maxlen)
		// Run gofumpt on the modified file
	}
}

// process modifies the file filename to split long comments at maxlen. making this line longer than 80 characters to test the lll fixer itself. fun no?
func process(filename string, maxlen int) {
	// Parse the Go file
	fset := token.NewFileSet()
	log.Infof("Processing file %q", filename)
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Errf("Error parsing %q: %v", filename, err)
		return
	}

	// Traverse and modify the AST
	ast.Inspect(node, func(n ast.Node) bool {
		// Split long comments
		if c, ok := n.(*ast.Comment); ok {
			if len(c.Text) > maxlen {
				log.LogVf("Splitting comment %q", c.Text)
				c.Text = strings.TrimSpace(c.Text[:maxlen-1]) + "\n// " + strings.TrimSpace(c.Text[maxlen-1:])
				log.LogVf("-> %q", c.Text)
			}
		}
		return true
	})

	// Generate the modified code
	newname := filename + ".lll"
	f, err := os.Create(newname)
	if err != nil {
		log.Errf("Error creating modified file %q: %v", newname, err)
	}
	defer f.Close()
	if err := format.Node(f, fset, node); err != nil {
		log.Errf("Error outputting modified file: %v", err)
	}
	log.Infof("Modified file written to %q", newname)
}
