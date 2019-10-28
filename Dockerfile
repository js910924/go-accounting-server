FROM golang:latest
MAINTAINER Jason jswind@myemail.com
LABEL description="This is a accounting server example" version="1.0" owner="Jason Chen"
RUN apt-get update
RUN apt-get install vim net-tools tmux
ADD go $HOME
COPY .vimrc ./
