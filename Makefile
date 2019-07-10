ALL: bin/proxy bin/runtests

bin/proxy: $(shell find src/proxy -type f -name '*.go')
	go install proxy

bin/runtests: $(shell find src/runtests -type f -name '*.go')
	go install runtests

etc/memcache/bin/proxy: $(shell find src/proxy -type f -name '*.go')
	docker run -t --rm \
		-v $(shell pwd)/src:/go/src \
		-v $(shell pwd)/etc/memcache/bin:/go/bin \
		golang:1.12 \
		go install proxy