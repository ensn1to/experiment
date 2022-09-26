FROM golang:1.18-alpine

# 需要指定工作目录，不然会报错
WORKDIR /go/delivery

COPY . .


CMD ["go", "run", "main.go"]