# build args
SERVICE_NAME		?=	ssp-service
VERSION			?=	latest
REGISTRY		?=	test
IMG			?=	$(REGISTRY)/$(SERVICE_NAME):$(VERSION)

# helpful things
CID 			?= 	$(shell docker ps --no-trunc -aqf name=$(SERVICE_NAME))
BUILD_TIME 		?= 	$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT_SHA 		?= 	$(shell git rev-parse HEAD)

DEFAULT_HOST		?= 0.0.0.0

# http args
HTTP_HOST		?= $(DEFAULT_HOST)
HTTP_PORT 		?= 9000
HTTP_TIMEOUT		?= 10s

# db args
DRIVER 			?= postgres
DB_HOST 		?= $(DEFAULT_HOST)
DB_PORT 		?= 5433
DB_USER 		?= bd0d8bc3d93622df2ddd645617f125fd7521b1ac9f24404b4f9b7e41e8c3cc92494d4b84
DB_PASSWORD		?= 75f90aeacaff10303672934cb8bd3e93dc072429c3eb999d2d8ca12f8ee1b173141359
DB_NAME 		?= test
DB_SCHEMA 		?= ssp
DB_MODE 		?= disable
DSN			?= $(DB_SCHEMA)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_MODE)

# docker args
D_DB_HOST 		?= 194.168.0.4
D_DB_PORT 		?= 5432

# secrets
SEC_SECRET		?= test

all: create run

.PHONY: all

# create binary file with race flag and specified name
create:
ifneq ("$(wildcard $(SERVICE_NAME))","")
	@rm -rf $(SERVICE_NAME) && go build -race -o $(SERVICE_NAME)
else
	@go build -race -o $(SERVICE_NAME)
endif

# run binary file
run:
	@./$(SERVICE_NAME) -host $(HTTP_HOST) -port $(HTTP_PORT) \
		-timeout $(HTTP_TIMEOUT) -db-driver $(DRIVER) \
		-db-host $(DB_HOST) -db-port $(DB_PORT) -db-user $(DB_USER) \
		-db-password $(DB_PASSWORD) -is-encrypted \
		-db-name $(DB_NAME) -db-schema $(DB_SCHEMA) \
		-db-mode $(DB_MODE)

test:
	@go test -count=1 -v -failfast ./...

local: test local-build local-start

# build docker container
local-build:
	@clear
	@echo "Building an image"
	@docker build -t $(IMG) --build-arg REGISTRY=$(REGISTRY) .
	@docker tag $(REGISTRY)/$(SERVICE_NAME) test-ssp

local-start:
	@clear
	@echo "Running container"
	@docker run --rm --name $(SERVICE_NAME) --network test2 --ip 194.168.0.3 -p 9001:9000/tcp \
		 -e SEC_SECRET=$(SEC_SECRET) $(IMG) \
		-host $(HTTP_HOST) -port $(HTTP_PORT) -timeout $(HTTP_TIMEOUT) \
		-db-driver $(DRIVER) -db-host $(D_DB_HOST) -db-port $(D_DB_PORT) -db-user $(DB_USER) \
		-db-password $(DB_PASSWORD) -is-encrypted -db-name $(DB_NAME) -db-schema $(DB_SCHEMA) -db-mode $(DB_MODE) \

stop:
	@docker stop $(CID)

# build docker image
image:
	@docker build --no-cache -t $(IMG) \
           --build-arg BUILD_TIME=$(BUILD_TIME) \
           --build-arg COMMIT_SHA=$(COMMIT_SHA) \
           --build-arg VERSION=$(VERSION)       \
           --build-arg REGISTRY=$(REGISTRY)  \
           .

# push docker image
push: image
	docker push $(IMG)
