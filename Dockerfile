FROM golang:alpine

WORKDIR /app
ADD bin/sam /app/sam
ENTRYPOINT ["/app/SAM"]