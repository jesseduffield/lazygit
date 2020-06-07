# run with:
# docker build -t lazygit .
# docker run -it lazygit:latest /bin/sh

FROM golang:1.14-alpine3.11
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.11
RUN apk add -U git xdg-utils
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY --from=0 /go/src/github.com/jesseduffield/lazygit /go/src/github.com/jesseduffield/lazygit
COPY --from=0 /go/src/github.com/jesseduffield/lazygit/lazygit /bin/
RUN echo "alias gg=lazygit" >> ~/.profile

ENTRYPOINT [ "lazygit" ]
