
REVISION := $(shell git rev-parse --short HEAD)
VERSION = $(REVISION)

.PHONY: dependent
dependent:
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest

.PHONY: install-dep
install-dep:
	go mod tidy

.PHONY: proto
proto:
	protoc -I ./api/file/service/v1 file.proto  --go_out=./ --micro_out ./
	protoc -I ./api/ocr/service/v1 ocr.proto  --go_out=./ --micro_out ./
	protoc -I ./api/translation/service/v1 translation.proto  --go_out=./ --micro_out ./
	protoc -I ./api/paper/service/v1 paper.proto  --go_out=./ --micro_out ./
	protoc -I ./api/email/service/v1 email.proto  --go_out=./ --micro_out ./

.PHONY: wire
wire:
	wire ./app/file/service
	wire ./app/ocr/service
	wire ./app/translation/service
	wire ./app/paper/service
	wire ./app/email/service

.PHONY: put-config
put-config:
	cat ./config/file-config.json | etcdctl put /configs/file-service
	cat ./config/ocr-config.json | etcdctl put /configs/ocr-service
	cat ./config/paper-config.json | etcdctl put /configs/paper-service
	cat ./config/translation-config.json | etcdctl put /configs/translation-service
	cat ./config/email-config.json | etcdctl put /configs/email-service
	cat ./config/frontend-config.json | etcdctl put /configs/frontend

.PHONY: build-file
build-file: install-dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o ./build/file-service ./app/file/service

.PHONY: build-ocr
build-ocr: install-dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o ./build/ocr-service ./app/ocr/service

.PHONY: build-translation
build-translation: install-dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o ./build/translation-service ./app/translation/service

.PHONY: build-paper
build-paper: install-dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o ./build/paper-service ./app/paper/service

.PHONY: build-email
build-email: install-dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o ./build/email-service ./app/email/service

.PHONY: build-frontend
build-frontend: install-dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o ./build/frontend ./app/frontend/service

.PHONY: image-file
image-file: build-file
	docker build ./build -f ./build/Dockerfile --build-arg service_name='file-service' --platform linux/amd64 -t file-service:$(VERSION)

.PHONY: image-ocr
image-ocr: build-ocr
	docker build ./build -f ./build/Dockerfile-ocr --platform linux/amd64 -t ocr-service:$(VERSION)

.PHONY: image-translation
image-translation: build-translation
	docker build ./build -f ./build/Dockerfile --build-arg service_name='translation-service' --platform linux/amd64 -t translation-service:$(VERSION)

.PHONY: image-paper
image-paper: build-paper
	docker build ./build -f ./build/Dockerfile --build-arg service_name='paper-service' --platform linux/amd64 -t paper-service:$(VERSION)

.PHONY: image-email
image-email: build-email
	docker build ./build -f ./build/Dockerfile --build-arg service_name='email-service' --platform linux/amd64 -t email-service:$(VERSION)

.PHONY: image-frontend
image-frontend: build-frontend
	docker build ./build -f ./build/Dockerfile --build-arg service_name='frontend' --platform linux/amd64 -t frontend:$(VERSION)

.PHONY: images
images: image-file image-ocr image-translation image-paper image-email image-frontend
	echo 'TAG=$(VERSION)' > .env

.PHONY: clean
clean:
	rm -rf ./build/file-service
	rm -rf ./build/ocr-service
	rm -rf ./build/translation-service
	rm -rf ./build/paper-service
	rm -rf ./build/email-service
	rm -rf ./build/frontend
