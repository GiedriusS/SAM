FROM golang:1.11

WORKDIR /app
ENV SRC_DIR=/go/src/github.com/GiedriusS/SAM
ADD . $SRC_DIR
RUN cd $SRC_DIR; go build -o SAM; cp SAM /app/
ENTRYPOINT ["/app/SAM"]