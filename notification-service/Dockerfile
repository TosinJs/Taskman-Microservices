FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . . 
RUN go build -o /docker-n-service
ENV PORT=4000
EXPOSE $PORT
CMD ["/docker-n-service"]