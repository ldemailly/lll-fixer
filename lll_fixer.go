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
	maxlen := flag.Int("len", 79, "maximum line length")
	funmpt := flag.Bool("fumpt", false, "run gofumpt on the modified file")
	cli.MinArgs = 1
	cli.MaxArgs = -1
	cli.ArgsHelp = "filenames..."
	if false {
		// just to test literal string split
		cli.ArgsHelp = "this is a very a long string literal to test the split of long strings inside code"
	}
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

func lineLead(s string) string {
	i := strings.LastIndex(s, "\n")
	return s[i+1 : i+4] // will break with -len too low
}

/*
 * Multi line comment with one line longer than 80 characters to test the split of multi line comments.
 */
func splitAtWord(s string, maxlen int) string {
	if len(s) <= maxlen {
		return s
	}
	// find the last space before maxlen
	i := strings.LastIndex(s[:maxlen], " ")
	nospace := (i == -1)
	if nospace {
		// no space found, split at maxlen
		log.Warnf("No word/space found in first %d characters for %q", maxlen, s)
		i = maxlen
	}
	start := s[:i]
	lead := lineLead(start)
	var mid string
	switch {
	case strings.HasPrefix(lead, "/* "):
		mid = "\n * "
	case strings.HasPrefix(lead, " * "):
		mid = "\n * "
	case strings.HasPrefix(lead, "// "):
		mid = "\n// "
	case strings.HasPrefix(lead, "\""):
		mid = "\" +\n\t\"" // for string literals splitting
		if !nospace {
			mid += " "
		}
	default:
		log.Warnf("Unexpected lead %q", lead)
		mid = "\n "
	}
	log.Debugf("Start lead is %q", lead)
	return strings.TrimSpace(s[:i]) + mid + strings.TrimLeft(s[i:], " ")
}

// TODO process other nodes (and also maybe split leftmost position in line vs length of comment/literal
// which could be far to the right)

func processNode(n ast.Node, maxlen int) bool {
	if false {
		log.Debugf("Found node: %+v to shrink to %d", n, maxlen)
	}
	// process string literals
	if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		lit.Value = splitAtWord(lit.Value, maxlen)
	}
	// more nodes...
	return true
}

// process modifies the file filename to split long comments at maxlen. making this line longer than 80 characters to test.
func process(fset *token.FileSet, filename string, maxlen int) string {
	log.Infof("Processing file %q", filename)
	// Parse the Go file and this is an indented comment to test the split past column 80.
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Error parsing %q: %v", filename, err)
		return "error.lll" // unreachable
	}
	for _, cg := range node.Comments {
		for _, c := range cg.List {
			log.Debugf("Found comment     %q", c.Text)
			if len(c.Text) > maxlen {
				log.LogVf("Splitting comment %q", c.Text)
				c.Text = splitAtWord(c.Text, maxlen)
				log.LogVf("into ->           %q", c.Text)
			}
		}
	}
	// Traverse and modify the AST
	ast.Inspect(node, func(n ast.Node) bool {
		return processNode(n, maxlen)
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
