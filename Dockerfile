FROM golang:latest
MAINTAINER Jason js910924@gmail.com
LABEL description="This is a accounting server" version="1.0" owner="Jason Chen"
RUN apt-get update
RUN apt-get install vim -y
RUN apt-get install net-tools -y
RUN apt-get install tmux -y
COPY ./ ./src/go-account
