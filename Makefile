.PHONY: all container clean run

all: container run

# Default settings: short name and basic web port
PORT ?= 8080
TAG ?= releasestatus

container: docker/Dockerfile docker/releasestatus
	docker build -t $(TAG) ./docker/

# run interactively to catch Ctrl-C
run:
	docker run --rm -ite RS_PORT=$(PORT) -p $(PORT):$(PORT) $(TAG)

# builds for Docker (Linux 64-bit) with built in netgo to handle web requests
docker/releasestatus: main.go
	GOOS=linux GOARCH=amd64 go build -v -o docker/releasestatus -tags netgo -installsuffix netgo .

clean:
	rm -f releasestatus
	rm -f docker/releasestatus
