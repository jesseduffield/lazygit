# run with:
# docker build -t lazygit .
# docker run -it lazygit:latest

FROM golang:alpine

RUN apk add -U git xdg-utils

ADD . /go/src/github.com/jesseduffield/lazygit

RUN go install github.com/jesseduffield/lazygit

WORKDIR /go/src/github.com/jesseduffield/lazygit