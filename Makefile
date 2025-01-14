.PHONY: buf-generate clean

create-output-dirs:
	@mkdir -p autogenerated/go
	@mkdir -p autogenerated/python
	@mkdir -p autogenerated/rust/src
 
buf-generate: clean
	@buf generate
	@cp proto/setup.py autogenerated/python/
	@cp proto/go.* autogenerated/go/

clean:
	@rm -rf autogenerated/*

buf-lint:
	@buf lint

buf-breaking:
	@buf breaking --against '.git#branch=main'

