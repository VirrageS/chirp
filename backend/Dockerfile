FROM golang:1.7

RUN adduser --disabled-password --gecos '' docker

WORKDIR /go/src/github.com/VirrageS/chirp/backend
COPY . /go/src/github.com/VirrageS/chirp/backend

RUN go get && go install
