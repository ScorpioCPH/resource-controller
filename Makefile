_output=_output

binaries=resource-controller

.PHONY: $(binaries)
$(binaries):
	go build -o $(_output)/$@ cmd/$@/main.go

.PHONY: clean
clean:
	rm -rf $(_output)
