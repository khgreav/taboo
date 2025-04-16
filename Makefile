phony: all build clean update

build:
	go build -o output/not-taboo-multiplayer

clean:
	rm -rf output

update:
	go get -u
	go mod tidy
