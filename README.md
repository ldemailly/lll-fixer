# lll-fixer
Fix lll (line length limit) lines too long linter errors in go files

## Installation

From source
```
go install github.com/ldemailly/lll-fixer@latest
```

Or see the numerous binaries in https://github.com/ldemailly/lll-fixer/releases

Or docker `fortio/lll-fixer:latest`

Or brew `brew install fortio/tap/lll-fixer`

(I manage the fortio org and usually put everything there but this one is a bit unrelated so for now it is hosted here in `ldemailly` yet uses fortio's org for docker and brew)

## Example

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
