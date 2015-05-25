#!/bin/bash

echo "this is broken!"
echo "warning: should ignore this" 1>&2
echo "big error!"
echo "warning: this can be ignored as well"
echo "broken sent to stderr" 1>&2

echo "before big sleep"
sleep 3
echo "after big sleep"

echo "$@";
