.PHONY: all container clean run

all: container run

# Default settings: short name, basic web port and Central European timezone
PORT ?= 8080
TZ ?= CET
TAG ?= releasestatus

container: docker/Dockerfile docker/releasestatus docker/zoneinfo.tar.gz
	docker build -t $(TAG) ./docker/

# run interactively to catch Ctrl-C
run:
	docker run --rm -ite TZ=$(TZ) -e RS_PORT=$(PORT) -p $(PORT):$(PORT) $(TAG)

# builds for Docker (Linux 64-bit) with built in netgo to handle web requests
docker/releasestatus: main.go
	GOOS=linux GOARCH=amd64 go build -v -o docker/releasestatus -tags netgo -installsuffix netgo .

# required for timezone info in scratch container
docker/zoneinfo.tar.gz: /usr/share/zoneinfo
	tar cfz docker/zoneinfo.tar.gz /usr/share/zoneinfo

clean:
	rm -f releasestatus
	rm -f docker/releasestatus
	rm -f docker/zoneinfo.tar.gz
