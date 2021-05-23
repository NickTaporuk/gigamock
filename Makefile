# to be up to date with goreleaser
# you should follow this link https://goreleaser.com/quick-start/
release:
	rm -rf ./dist
	@goreleaser release