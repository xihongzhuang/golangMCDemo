## install : Build and install webrtc-signaling-server
## make run to start all process needed to test grpc-web
EXE:=mcCodingExercise
DEMO_IMAGE := gcr.io/mc-coding-exercise:latest
SVR=mc-demo

.PHONY: install
install:
	go build -o ./$(EXE) .
	#go build -o $(GOBIN)/$(EXE) .

.PHONY: install
clear:
	rm ./$(EXE)

SVR_ID := $(shell docker ps -f name=$(SVR) -q)
.PHONY: run
run:
ifndef SVR_ID
	docker run -d --name=$(SVR) -e PORT=8080 -p 8080:8080 --expose 8080 $(DEMO_IMAGE)
endif

.PHONY: stop
stop:
ifdef SVR_ID
	docker container rm --force $(SVR_ID)
endif

.PHONY: build
build:
	cd ../../ && docker build -f Dockerfile -t $(DEMO_IMAGE) .
