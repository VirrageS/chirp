FROM golang:1.7

# copy all files
WORKDIR /go/src/github.com/VirrageS/chirp/backend
COPY . /go/src/github.com/VirrageS/chirp/backend

# get dependencies and compile binary
RUN go get -d && go install
