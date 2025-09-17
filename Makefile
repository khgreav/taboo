phony: all build clean update

build:
	go build -o output/taboo-mp

run:
	./output/taboo-mp

clean:
	rm -rf output

update:
	go get -u
	go mod tidy
