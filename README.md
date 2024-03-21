# lll-fixer
Fix lll (line length limit) lines too long linter errors in go files

Test on itself:
```
$ go run . lll_fixer.go
```

```diff
diff --git a/lll_fixer.go b/lll_fixer.go
index c5edf97..30452c3 100644
--- a/lll_fixer.go
+++ b/lll_fixer.go
@@ -43,7 +43,8 @@ func main() {
        }
 }

-// process modifies the file filename to split long comments at maxlen. making this line longer than 80 characters to test.
+// process modifies the file filename to split long comments at maxlen. making
+// this line longer than 80 characters to test.
 func process(fset *token.FileSet, filename string, maxlen int) string {
        log.Infof("Processing file %q", filename)
        // Parse the Go file
```
