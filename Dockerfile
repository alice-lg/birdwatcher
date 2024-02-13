
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
	coreutils \
	linux-headers \
	ncurses-static \
	readline-dev \
	readline-static

# Clone the latest version 2 tag of the bird repository
RUN git clone \
	--branch $(git ls-remote --tags https://gitlab.nic.cz/labs/bird | awk -F'/' '{print $3}' | grep '^v2\.' | grep -v '{}' | sort -V | tail -n 1) \
	https://gitlab.nic.cz/labs/bird.git
WORKDIR /src/bird
RUN autoreconf && \
	LDFLAGS="-static -static-libgcc" \
		./configure \
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

