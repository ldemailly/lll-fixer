# lll-fixer
Fix lll (line length limit) lines too long linter errors in go files

Currently weird bug with `-len 79`: (2 new lines somehow)

```diff
21:02:56 main lll-fixer/$ colordiff -u lll_fixer.go lll_fixer.go.lll
--- lll_fixer.go	2024-03-20 21:02:46.473102471 -0700
+++ lll_fixer.go.lll	2024-03-20 21:02:56.382460319 -0700
@@ -24,7 +24,9 @@
 	}
 }

-// process modifies the file filename to split long comments at maxlen. making this line longer than 80 characters to test the lll fixer itself. fun no?
+// process modifies the file filename to split long comments at maxlen. making
+//
+//	this line longer than 80 characters to test the lll fixer itself. fun no?
 func process(filename string, maxlen int) {
 	// Parse the Go file
 	fset := token.NewFileSet()
```
