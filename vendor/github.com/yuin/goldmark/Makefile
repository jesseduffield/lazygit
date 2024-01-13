.PHONY: test fuzz

test:
	go test -coverprofile=profile.out -coverpkg=github.com/yuin/goldmark,github.com/yuin/goldmark/ast,github.com/yuin/goldmark/extension,github.com/yuin/goldmark/extension/ast,github.com/yuin/goldmark/parser,github.com/yuin/goldmark/renderer,github.com/yuin/goldmark/renderer/html,github.com/yuin/goldmark/text,github.com/yuin/goldmark/util ./...

cov: test
	go tool cover -html=profile.out

fuzz:
	cd ./fuzz && go test -fuzz=Fuzz
