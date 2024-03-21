package main

import (
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"

	"fortio.org/cli"
	"fortio.org/log"
)

func main() {
	maxlen := flag.Int("len", 79, "max line length")
	funmpt := flag.Bool("fumpt", false, "run gofumpt on the modified file")
	cli.MinArgs = 1
	cli.MaxArgs = -1
	cli.ArgsHelp = "filenames..."
	cli.Main()
	fset := token.NewFileSet()
	for _, filename := range flag.Args() {
		newname := process(fset, filename, *maxlen)
		// swap .lll to .go and .go to .bak
		backup := filename + ".bak"
		if err := os.Rename(filename, backup); err != nil {
			log.Fatalf("Error renaming file %q to %q: %v", filename, backup, err)
		}
		log.Infof("Renamed file %q to %q", filename, backup)
		if err := os.Rename(newname, filename); err != nil {
			log.Fatalf("Error renaming file %q to %q: %v", newname, filename, err)
		}
		log.Infof("Renamed file %q to %q", newname, filename)
		// Run gofumpt on the modified file
		if *funmpt {
			cmd := exec.Command("gofumpt", "-w", filename)
			if err := cmd.Run(); err != nil {
				log.Errf("Error running gofumpt: %v", err)
				return
			}
			log.Infof("Ran gofumpt on the now modified file %q", filename)
		}
	}
}

// process modifies the file filename to split long comments at maxlen. making this line longer than 80 characters to test.
func process(fset *token.FileSet, filename string, maxlen int) string {
	log.Infof("Processing file %q", filename)
	// Parse the Go file
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Error parsing %q: %v", filename, err)
		return "error.lll" // unreachable
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
	return newname
}
