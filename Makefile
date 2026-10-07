GDM_CLI := bin/gdm-cli
GDM_GEN := bin/gdm-gen

.PHONY: build gen check lint generate site site-serve

build:
	go build -o $(GDM_CLI) ./gdm/cmd/gdm-cli
	go build -o $(GDM_GEN) ./gdm/cmd/gdm-gen

generate:
	go generate ./...

gen: build
	$(GDM_GEN) --out manual

lint:
	golangci-lint run ./... --fix

check:
	go vet ./...
	go test ./...

site:
	mkdocs build -f website/mkdocs.yml

site-serve:
	mkdocs serve -f website/mkdocs.yml
