FROM golang
RUN git clone https://github.com/aeleniumfor/Practice-Gin_ToDoApp.git
WORKDIR /Practice-Gin_ToDoApp
RUN go mod init && \
    go run main.go
