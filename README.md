## cronsifter

A command line tool for filtering (sifting) the results of a command
that has been run typically under cron, but could be used elsewhere.

## Installing

Use the 'go get' command to download the package.

    go get github.com/outtenr/cronsifter

To install:

    go install github.com/outten45/cronsifter/cmd/...


This will create the `cronsifter` and `srunner` binaries in your $GOPATH/bin.
Add it to your PATH if it isn't already there.

    export PATH=$PATH:$GOPATH/bin

## Building

To build the linux binary:

    GOOS=linux GOARCH=amd64 go install github.com/outten45/cronsifter/cmd/...

Look in the `$GOPATH/bin/linux_amd64/` directory for the `cronsifter` and
`srunner`.

## Binaries

### cronsifter

Here is the general usage of the `cronsifter` command.

    Usage of cronsifter:
      -count int
            rotate log file count for stdout and stderr files (default 20)
      -dir string
            directory to write standard out and err log files to (default ".")
      -excludes string
            the file containing the regular expressions to remove from stdout/stderr
      -size int
            file size in megabytes to rotate stdout and stderr files (default 20)


The cronsifter command can be run 2 different ways. The command can either call
a job command to be run.  The other option is to pipe the output of the job
command to cronsifter to then filter the results.

**Calling the job command from cronsifter**

    > cronsifter -excludes=excludes.txt command_to_run.sh

**Piping the stdout/stderr to cronsifter**

    > command_to_run.sh | cronsifter -excludes=excludes.txt

### srunner

A simple program to run the given command and keep it running. It assumes the
command is meant to run in the foreground so if the command stops, `srunner`
will restart it. `srunner` will default to waiting 10 seconds before starting
the command again.  You can change this by setting the SRUNNER_SPAWN_DELAY
environment variable which is parsed using
[ParseDuration](https://golang.org/pkg/time/#ParseDuration).

    > SRUNNER_SPAWN_DELAY="3m" srunner gohttp

This would start the gohttp command and if it failed, `srunner` would wait 3
minutes before respawning the process.

Thanks to [goforever](https://github.com/gwoo/goforever) as a reference on how
to keep a subprocess running.
