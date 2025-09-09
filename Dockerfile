# run with:
# docker build -t lazygit .
# docker run -it lazygit:latest

FROM golang:1.24 AS build
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor

FROM alpine:3.19
RUN apk add --no-cache -U git xdg-utils
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY --from=build /go/src/github.com/jesseduffield/lazygit ./
COPY --from=build /go/src/github.com/jesseduffield/lazygit/lazygit /bin/
RUN echo "alias gg=lazygit" >> ~/.profile

ENTRYPOINT [ "lazygit" ]
