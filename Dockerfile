# run with:
# docker build -t lazygit .
# docker run -it lazygit:latest

FROM golang:alpine
WORKDIR /go/src/github.com/jesseduffield/lazygit/
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lazygit .

FROM alpine:latest
RUN apk add -U git xdg-utils
WORKDIR /root/
COPY --from=0 /go/src/github.com/jesseduffield/lazygit/lazygit .
CMD ["./lazygit"]
