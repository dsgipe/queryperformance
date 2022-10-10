FROM golang:latest
WORKDIR /app
COPY . /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -o queryperformance .
ENTRYPOINT [ "/app/queryperformance" ]