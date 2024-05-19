# run with:
# docker build -t lazygit .
# docker run -it lazygit:latest /bin/sh

FROM golang:1.22 as build
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.15
RUN apk add --no-cache -U git xdg-utils
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY --from=build /go/src/github.com/jesseduffield/lazygit ./
COPY --from=build /go/src/github.com/jesseduffield/lazygit/lazygit /bin/
RUN echo "alias gg=lazygit" >> ~/.profile

ENTRYPOINT [ "lazygit" ]
