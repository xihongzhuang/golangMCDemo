# this dockerfile must be run from the root of the repo

FROM golang:alpine AS builder

RUN apk add --no-cache make

ADD . /src
WORKDIR /src/cmd/simple-rest-service/

RUN make install

FROM alpine

COPY --from=builder /src/cmd/simple-rest-service/mcCodingExercise /mcCodingExercise

CMD ["/mcCodingExercise"]
