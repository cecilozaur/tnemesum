.DEFAULT_GOAL := build

build:
	go build -o app -ldflags '-w -s'

clean:
	go clean

run:
	docker build -t tnemesum .
	docker run -it --rm -p 8000:8000 --name tnemesum-api tnemesum