default:
	@echo Nothing is default

go-tangle: go-tangle.go go-tangle.md
	@go vet .
	@go build .

.PHONY: godoc
godoc: go-tangle.md

go-tangle.md: go-tangle.go
	@go doc -all . > $@


# run awk -F: '/^[a-z]+/ {print $1}' Makefile > TARGETS
# to generate TARGETS file
.PHONY: clean
clean:
	@rm -f `cat TARGETS` *~
