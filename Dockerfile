
#
# Birdwatcher - Your friendly alice looking glass data source
#

# Build birdwatcher
FROM golang:1.13 AS birdwatcher

WORKDIR /src/birdwatcher
ADD . .
RUN go mod download
RUN make linux_static

# Build bird
FROM alpine:latest AS bird
WORKDIR /src
RUN apk add --no-cache \
	gcc \
	make \
	musl-dev \
	autoconf \
	automake \
	flex \
	bison \
	git \
	curl \
	coreutils \
	linux-headers \
	ncurses-static \
	readline-dev \
	readline-static

# Fetch and build the latest version 2 of the BIRD release
RUN \
	birdRev=$(git ls-remote --tags https://gitlab.nic.cz/labs/bird | awk -F'/' '{print $3}' | grep '^v2\.' | grep -v '{}' | sort -V | tail -n 1 | sed -e 's/v//') && \
  curl -Lo - "https://bird.network.cz/download/bird-${birdRev}.tar.gz" | tar -C /src -zxvf - && \
  mv /src/bird-${birdRev} /src/bird && cd /src/bird && \
	LDFLAGS="-static -static-libgcc" ./configure \
	 		--prefix=/ \
	 		--exec-prefix=/usr \
	 		--runstatedir=/run/bird && \
	make -j

# Final stage
FROM scratch

COPY --from=bird /src/bird/birdcl /usr/bin/birdc
COPY --from=birdwatcher /src/birdwatcher/birdwatcher-linux-amd64 /usr/bin/birdwatcher
COPY --from=birdwatcher /src/birdwatcher/etc/birdwatcher/birdwatcher.conf /etc/birdwatcher/birdwatcher.conf

EXPOSE 29184/tcp
EXPOSE 29186/tcp

CMD ["/usr/bin/birdwatcher", "-config", "/etc/birdwatcher/birdwatcher.conf"]

