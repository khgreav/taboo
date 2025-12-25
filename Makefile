phony: all update run clean build-backend build-frontend build-server-rpi \
		deps-update-backend deps-update-frontend \
		run-backend run-frontend

# all

all: build-backend build-frontend

update: update-go-deps update-js-deps

clean:
	rm -rf output/

# backend

build-backend:
	cd backend && go build -o ../output/taboo-server

build-server-rpi:
	GOOS=linux GOARCH=arm64 go build -o output/taboo-server

deps-update-backend:
	go get -u
	go mod tidy

run-backend:
	cd backend && go run .

# frontend

build-frontend:
	pnpm --dir frontend install
	pnpm --dir frontend build
	rm -rf output/static
	mkdir -p output/static
	mv frontend/dist/* output/static/

deps-update-frontend:
	pnpm --dir frontend update

run-frontend:
	pnpm --dir frontend dev
