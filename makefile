appname := QuickClip-0.7.1
applicationDir := QuickClip

sources := $(wildcard *.go)

build =  cd ./cmd && GOOS=$(1) GOARCH=$(2) go build -o ../bin/$(appname)$(3)
tar =  mkdir -p build && cd ../ && tar --exclude $(applicationDir)/static/js/src -cvzf ./$(appname)_$(1)_$(2).tar.gz $(applicationDir)/bin $(applicationDir)/config $(applicationDir)/static $(applicationDir)/templates && mv $(appname)_$(1)_$(2).tar.gz $(applicationDir)/build

.PHONY: all linux

all:	linux


run: web
	cd ./cmd && go run .


web:
	tsc -b
	minify-all-js static/js/out

build: web all linux clean

clean:
	rm -rf bin
	rm -rf log
	rm -rf static/js/out

cleanall: clean
	rm -rf build

##### LINUX BUILDS #####
linux: build/linux_arm.tar.gz build/linux_arm64.tar.gz build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: $(sources)
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)
	$(call tar,linux,arm)

build/linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

