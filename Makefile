
manual-test:
	cp lll_fixer.go test.txt && go run . -loglevel debug test.txt ; colordiff -u lll_fixer.go test.txt


.PHONY: manual-test
