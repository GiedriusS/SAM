FROM golang:alpine

WORKDIR /app
ADD bin/SAM /app/SAM
ENTRYPOINT ["/app/SAM"]