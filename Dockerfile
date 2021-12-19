FROM ubuntu

ADD target/hook   /app/hook
ADD target/server /app/server
ADD init.sh       /app/init.sh
ADD conf.json     /app/conf.json
