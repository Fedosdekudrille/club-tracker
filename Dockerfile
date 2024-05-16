FROM golang:1.22.3-alpine3.19

WORKDIR /app

COPY . .

ENV FILE_PATH=assets/test_file.txt
#CMD ["go", "run", "cmd/tracker/main.go", "assets/test_file.txt"]

