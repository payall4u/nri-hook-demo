#!/bin/bash

mkdir    /host/etc/nri
mkdir -p /host/opt/nri/bin

cp /app/hook        /host/opt/nri/bin
cp /app/conf.json   /host/etc/nri