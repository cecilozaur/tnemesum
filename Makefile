.DEFAULT_GOAL := build

build:
	go build -o app -ldflags '-w -s'

clean:
	go clean

run:
	docker build -t tnemesum .
	docker run -it --rm --name tnemesum-api tnemesum