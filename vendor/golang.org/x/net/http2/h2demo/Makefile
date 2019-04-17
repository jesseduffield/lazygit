# Copyright 2018 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

MUTABLE_VERSION ?= latest
VERSION ?= $(shell git rev-parse --short HEAD)
IMAGE_NAME=h2demo

update-deps:
	go install golang.org/x/build/cmd/gitlock
	gitlock --update=Dockerfile --ignore=golang.org/x/net --tags=h2demo golang.org/x/net/http2/h2demo

docker: Dockerfile
	go install golang.org/x/build/cmd/xb
	xb docker build -t golang/$(IMAGE_NAME) -f Dockerfile ../..

push-staging: docker
	go install golang.org/x/build/cmd/xb
	xb --staging docker tag golang/$(IMAGE_NAME) REPO/$(IMAGE_NAME):$(VERSION)
	xb --staging docker tag golang/$(IMAGE_NAME) REPO/$(IMAGE_NAME):latest
	xb --staging docker push REPO/$(IMAGE_NAME):$(VERSION)
	xb --staging docker push REPO/$(IMAGE_NAME):latest

push-prod: docker
	go install golang.org/x/build/cmd/xb
	xb --prod docker tag golang/$(IMAGE_NAME) REPO/$(IMAGE_NAME):$(VERSION)
	xb --prod docker tag golang/$(IMAGE_NAME) REPO/$(IMAGE_NAME):latest
	xb --prod docker push REPO/$(IMAGE_NAME):$(VERSION)
	xb --prod docker push REPO/$(IMAGE_NAME):latest

# TODO(bradfitz): add REPO subsitution in xb for the kubectl command.
deploy-prod: push-prod
	xb --prod kubectl set image deployment/h2demo-deployment h2demo=gcr.io/symbolic-datum-552/h2demo:$(VERSION)
deploy-staging: push-staging
	xb --staging kubectl set image deployment/h2demo-deployment h2demo=gcr.io/go-dashboard-dev/h2demo:$(VERSION)

FORCE:
