# run with:
# docker build -t lazytask .
# docker run -it lazytask:latest /bin/sh

FROM golang:1.21 as build
WORKDIR /go/src/github.com/lobes/lazytask/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.15
RUN apk add --no-cache -U git xdg-utils
WORKDIR /go/src/github.com/lobes/lazytask/
COPY --from=build /go/src/github.com/lobes/lazytask ./
COPY --from=build /go/src/github.com/lobes/lazytask/lazytask /bin/
RUN echo "alias lt=lazytask" >> ~/.profile

ENTRYPOINT [ "lazytask" ]
