GDM_CLI := bin/gdm-cli
GDM_GEN := bin/gdm-gen

SITE_ASSET_DIRS := manual/zh/assets manual/en/assets

.PHONY: build gen check lint lint-fix generate site site-assets site-serve

build:
	go build -o $(GDM_CLI) ./gdm/cmd/gdm-cli
	go build -o $(GDM_GEN) ./gdm/cmd/gdm-gen

generate:
	go generate ./...

gen: build
	$(GDM_GEN) --out manual

lint:
	go tool golangci-lint run ./...

lint-fix:
	go tool golangci-lint run ./... --fix

check:
	go vet ./...
	go test ./...

site-assets:
	@for dir in $(SITE_ASSET_DIRS); do rm -rf $$dir; mkdir -p $$dir; cp -R website/assets/. $$dir/; done

site: gen site-assets
	mkdocs build -f website/mkdocs.yml

site-serve: gen site-assets
	mkdocs serve -f website/mkdocs.yml
