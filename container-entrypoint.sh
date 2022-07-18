#!/usr/bin/env bash

set -e

if [ -z $PREFIX ]; then
    /app/systemstat
else
    /app/systemstat --prefix=$PREFIX
fi