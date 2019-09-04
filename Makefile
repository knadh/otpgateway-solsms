.PHONY: build
build:
	# Compile the plugin.
	go build -ldflags="-s -w" -buildmode=plugin -o solsms.prov solsms.go
