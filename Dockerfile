FROM golang:1.16

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE ${PORT}

ENTRYPOINT CompileDaemon --build="go build -o main_app" --command=./main_app