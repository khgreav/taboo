phony: all build clean update

build: build-server build-frontend

build-server:
	go build -o output/taboo-server

build-server-rpi:
	GOOS=linux GOARCH=arm64 go build -o output/taboo-server

build-frontend:
	pnpm --dir frontend install
	pnpm --dir frontend build
	mkdir -p output/static
	mv frontend/dist/* output/static/

run:
	./output/taboo-server

clean:
	rm -rf output

update: update-go-deps update-js-deps

update-go-deps:
	go get -u
	go mod tidy

update-js-deps:
	pnpm --dir frontend update
