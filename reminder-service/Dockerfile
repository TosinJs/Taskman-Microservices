FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . . 
RUN go build -o /docker-rm-service
ENV PORT=8080
EXPOSE $PORT
CMD ["/docker-rm-service"]