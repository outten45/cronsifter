## cronsifter

A command line tool for filtering (sifting) the results of a command
that has been run typically under cron, but could be used elsewhere.

## Installing

Use the 'go get' command to download the package.

    go get github.com/outtenr/cronsifter

This will create the cronsifter binary in your $GOPATH/bin. Add it to
your PATH if it isn't already there.

    export PATH=$PATH:$GOPATH/bin

## Building

To build the linux binary:

    GOOS=linux go build -v .
