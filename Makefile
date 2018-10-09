CONTAINER=aa-server
CI_JOB_ID ?= 0-dev
VERSION=0.1.$(CI_JOB_ID)

build:
	docker build -t $(CONTAINER):$(VERSION) .

run: build
	docker stop $(CONTAINER) || true
	mkdir -p data
	chmod g+w data
	docker run --rm -d --name $(CONTAINER) -p 9095:9095 -v $(shell pwd)/data:/home/app/data $(CONTAINER):$(VERSION) /home/app/aa-server --behind-proxy --data-dir /home/app/data