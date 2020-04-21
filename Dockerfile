FROM golang

RUN mkdir -p /app

WORKDIR /app

COPY . ./

RUN go build logserver.go

EXPOSE 11000

ENTRYPOINT [ "./logserver" ]