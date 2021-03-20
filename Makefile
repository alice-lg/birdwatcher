#
# Birdseye Makefile
#

PROG=birdwatcher
ARCH=amd64

APP_VERSION=$(shell cat VERSION)
VERSION=$(APP_VERSION)_$(shell git rev-parse --short HEAD)

BUILD_SERVER=''

SYSTEM_INIT=systemd

DIST=DIST/
REMOTE_DIST=$(PROG)-$(DIST)

RPM=$(PROG)-$(VERSION)-1.x86_64.rpm

LOCAL_RPMS=RPMS

# OS Detection
UNAME=$(shell uname)
ifeq ($(UNAME), Darwin)
  TARGET=osx
endif
ifeq  ($(UNAME), FreeBSD)
  TARGET=freebsd
endif
ifeq  ($(UNAME), Linux)
  TARGET=linux
endif
ifneq ($(UNAME),$(filter $(UNAME),Darwin FreeBSD Linux))
  $(error error: Unkown OS )
endif

all: $(TARGET)
	@echo "Built $(VERSION) @ $(TARGET)"

osx:
	GO111MODULE=on GOARCH=$(ARCH) GOOS=darwin go build -o $(PROG)-osx-$(ARCH)

linux:
	GO111MODULE=on GOARCH=$(ARCH) GOOS=linux go build -o $(PROG)-linux-$(ARCH)

freebsd:
	GO111MODULE=on GOARCH=$(ARCH) GOOS=freebsd go build -o $(PROG)-freebsd-$(ARCH)


build_server:
ifeq ($(BUILD_SERVER), '')
	$(error BUILD_SERVER not configured)
endif

dist: clean linux

	mkdir -p $(DIST)opt/birdwatcher/birdwatcher/bin
	mkdir -p $(DIST)etc/birdwatcher

ifeq ($(SYSTEM_INIT), systemd)
	# Installing systemd services
	mkdir -p $(DIST)usr/lib/systemd/system/
	cp install/systemd/* $(DIST)usr/lib/systemd/system/.
else
	# Installing upstart configuration
	mkdir -p $(DIST)/etc/init/
	cp install/upstart/init/* $(DIST)etc/init/.
endif


	# Copy config and startup script
	cp etc/birdwatcher/* DIST/etc/birdwatcher/.
	rm -f DIST/etc/birdwatcher/*.local.*

	# Copy bin
	cp $(PROG)-linux-$(ARCH) DIST/opt/birdwatcher/birdwatcher/bin/.


release: linux

	mkdir -p ../birdseye-static/birdwatcher-builds/$(APP_VERSION)/
	cp birdwatcher-linux-amd64 ../birdseye-static/birdwatcher-builds/$(APP_VERSION)/
	rm -f ../birdseye-static/birdwatcher-builds/latest
	cd ../birdseye-static/birdwatcher-builds && ln -s $(APP_VERSION) latest


rpm: dist

	# Clear tmp failed build (if any)
	mkdir -p $(LOCAL_RPMS)

	# Create RPM from dist
	fpm -s dir -t rpm -n $(PROG) -v $(VERSION) -C $(DIST) \
		--config-files /etc/birdwatcher/birdwatcher.conf \
		opt/ etc/

	mv $(RPM) $(LOCAL_RPMS)


remote_rpm: build_server dist

	mkdir -p $(LOCAL_RPMS)

	# Copy distribution to build server
	ssh $(BUILD_SERVER) -- rm -rf $(REMOTE_DIST)
	scp -r $(DIST) $(BUILD_SERVER):$(REMOTE_DIST)
	ssh $(BUILD_SERVER) -- fpm -s dir -t rpm -n $(PROG) -v $(VERSION) -C $(REMOTE_DIST) \
		--config-files /etc/birdwatcher/birdwatcher.conf \
		opt/ etc/

	# Get rpm from server
	scp $(BUILD_SERVER):$(RPM) $(LOCAL_RPMS)/.


.PHONY: test clean
test:
	go test -v
	cd endpoints/ && go test -v
	cd bird/ && go test -v

clean:
	rm -f $(PROG)-osx-$(ARCH)
	rm -f $(PROG)-linux-$(ARCH)
	rm -rf $(DIST)

