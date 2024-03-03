# format source
.PHONY: fmt
fmt:
	@gofmt -w *.go


# count code lines
.PHONY: cloc
cloc:
	@cloc *.go