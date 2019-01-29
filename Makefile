.PHONY: build
build:
	# Compile the plugin.
	go build -ldflags="-s -w" -buildmode=plugin -linkshared -o solsms.prov solsms.go
