
#
# Birdwatcher - Your friendly alice looking glass data source
#

FROM golang:1.13 AS app

WORKDIR /src/birdwatcher
ADD vendor .
ADD go.mod .
ADD go.sum .
RUN go mod download

# Add sourcecode
ADD . .

# Build birdwatcher
RUN make

# Add birdwatcher to bird
FROM ehlers/bird2

COPY --from=app /src/birdwatcher/birdwatcher-linux-amd64 /usr/bin/birdwatcher
ADD etc/birdwatcher/birdwatcher.conf /etc/birdwatcher/birdwatcher.conf

EXPOSE 29184/tcp
EXPOSE 29186/tcp

ENTRYPOINT ["/usr/bin/birdwatcher", "-config", "/etc/birdwatcher/birdwatcher.conf"]

