build:
	go build -C cmd/link-forge/ -o ../../bin/link-forge

run: build
	CONFIG_PATH=.env ./bin/link-forge
