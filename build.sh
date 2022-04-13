#!/bin/bash
mkdir -p bin
cd bin
go build ..
cd ..
cp config.ini bin/config.ini
