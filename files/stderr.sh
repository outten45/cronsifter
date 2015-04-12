#!/bin/bash

echo "this is broken!"
echo "warning: should ignore this" 1>&2
echo "big error!"

echo "$@";
