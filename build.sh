#!/bin/bash
# Builds all the teslalogger-cli commands
mkdir -p bin
for f in `ls *.go`;
do
 echo "building $f"
 go build -o bin/ $f
done

