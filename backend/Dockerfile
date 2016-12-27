FROM golang:1.7

# copy all files
WORKDIR /go/src/github.com/VirrageS/chirp/backend
COPY . /go/src/github.com/VirrageS/chirp/backend
COPY docker_config.yaml /go/src/github.com/VirrageS/chirp/backend/config.yaml

# get dependencies and compile binary
RUN make install
