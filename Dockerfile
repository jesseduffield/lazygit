# run with:
# docker build -t lazygit .
# docker run -it lazygit:latest /bin/sh -l

FROM golang:alpine
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -o lazygit .

FROM alpine:latest
RUN apk add -U git xdg-utils
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY --from=0 /go/src/github.com/jesseduffield/lazygit /go/src/github.com/jesseduffield/lazygit
COPY --from=0 /go/src/github.com/jesseduffield/lazygit/lazygit /bin/
RUN echo "alias gg=lazygit" >> ~/.profile
