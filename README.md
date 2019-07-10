# birdwatcher

birdwatcher is a small HTTP server meant to provide an API defined by
Barry O'Donovan's
[birds-eye](https://github.com/inex/birds-eye-design/) to
[the BIRD routing daemon](http://bird.network.cz/).

## Why

The [INEX implementation](https://github.com/inex/birdseye) of
birdseye runs PHP, which is not always desirable (and performant)
in a routeserver setting. By using Go, we are able to work with
regular binaries, which means deployment and maintenance might be
more convenient.

Our version also has a few more capabilities, as you will
discover when looking at [the modules section](https://github.com/alice-lg/birdwatcher/blob/master/etc/birdwatcher/birdwatcher.conf)
of the config.

## Installation

You will need to have go installed to build the package.
Please make sure your go version is `>= 1.9`.

Running `go get github.com/alice-lg/birdwatcher` will give you
a binary. You might need to cross-compile it for your
bird-running servive (`GOARCH` and `GOOS` are your friends).

We provide a Makefile for more advanced compilation/configuration.
Running `make linux` will create a Linux executable (by default for
`amd64`, but that is configurable by providing the `ARCH` argument
to the Makefile).


#### 2.0 Breaking Change

The per peer table configuration is no longer done in the birdwatcher,
but directly in alice.


### BIRD configuration

Birdwatcher parses the output of birdc and expects (for now)
the time format to be `iso long`. You need to configure

    timeformat base         iso long;
    timeformat log          iso long;
    timeformat protocol     iso long;
    timeformat route        iso long;

in your `/etc/bird[6].conf` for birdwatcher to work.

#### BIRD keep filtered routes
To also see the filtered routes in BIRD you need to make sure that you
have enabled the 'import keep filtered on' option for your BGP peers.

    protocol bgp 'peerX' {
        ...
        import keep filtered on;
        ...
    }

Now you should be able to do a 'show route filtered' in BIRD.

Do note that 'import keep filtered on' does NOT work for BIRD's pipe protocol
which is used when you have per peer tables, often used with Route Servers. If
your BIRD configuration has its import filters set on the BIRD pipe protocols
themselves then you will not be able to show the filtered routes.
However, you could move the import filters from the pipes to the BGP protocols
directly. For example:

    table master;
    table table_peer_X;

    protocol pipe pipe_peer_X {
        table master;
        peer table table_peer_X;
        mode transparent;
        import all;
        export where exportMagic();
    }

    protocol bgp 'peerX' {
        ...
        table table_peer_X;
        import where importFilter();
        import keep filtered on;
        export all;
        ...
    }

#### BIRD tagging filtered routes
If you want to make use of the filtered route reasons in the Birdseye then you need
to make sure that you are using BIRD 1.6.3 or up as you will need Large BGP Communities
(http://largebgpcommunities.net/).

You need to add a Large BGP Community just before you filter a route, for example:

    define yourASN = 12345
    define yourFilteredNumber = 65666
    define prefixTooLong = 1
    define pathTooLong = 2

    function importScrub() {
        ...
        if (net.len > 24) then {
            print "REJECTING: ",net.ip,"/",net.len," received from ",from,": Prefix is longer than 24: ",net.len,"!";
            bgp_large_community.add((YourASN,yourFilteredNumber,prefixTooLong));
            return false;
        }
        if (bgp_path.len > 64) then {
            print "REJECTING: ",net.ip,"/",net.len," received from ",from,": AS path length is ridiculously long: ",bgp_path.len,"!";
            bgp_large_community.add((yourASN,yourFilteredNumber,pathTooLong));
            return false;
        }
        ...
        return true;
    }

    function importFilter() {
        ...
        if !(importScrub()) then reject;
        ...
        accept;
    }

### Building an RPM

Building RPMs is supported through [fpm](https://github.com/jordansissel/fpm).
If you have `fpm` installed locally, you can run `make rpm`
to create a RPM in the folder `RPMS`. If you have a remote
build server with `fpm` installed, you can build and fetch
an RPM with `make remote_rpm BUILD_SERVER=<buildserver_url>`
(requires SSH access).

### Deployment

If you want to deploy `birdwatcher` on a system that uses
RPMs, you should be able to install it after following the
instructions on [building an RPM](#building-an-rpm).

We do not currently support other deployment methods.

## Configuration

An example config with sane defaults is provided in
[etc/birdwatcher/birdwatcher.conf](https://github.com/alice-lg/birdwatcher/blob/master/etc/birdwatcher/birdwatcher.conf).
You should be able to use it out of the box. If you need
to change it, it is well-commented and hopefully intuitive.
If you do not know how to configure it, please consider opening
[an issue](https://github.com/alice-lg/birdwatcher/issues/new).

## How

In the background `birdwatcher` runs the `birdc` client, sends
commands and parses the result. It also leverages simple caching
techniques to help reduce the load on the bird service.

## Who

Initially developed by Daniel and MC from [Netnod](https://www.netnod.se/) in
two days at the RIPE 73 IXP Tools Hackathon in Madrid, Spain.

Running bird and parsing the results was added by [Veit Heller](https://github.com/hellerve/) on behalf of [ecix](http://ecix.net/).

With major contributions from: Patrick Seeburger and Benedikt Rudolph on behalf of [DE-CIX](https://de-cix.net).

