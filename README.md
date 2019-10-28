# Accounting Server

###### tags: `Golang` `Docker` `Ubuntu18.04`

---

## Create Dockerfile
```shell
$ sudo touch Dockerfile
$ sudo vim Dockerfile
```	

```Dockerfile
FROM golang:latest
MAINTAINER Jason jswind@myemail.com
LABEL description="This is a accounting server example" version="1.0" owner="Jason Chen"
RUN apt-get update
RUN apt-get install vim net-tools tmux
RUN mkdir $HOME/go
COPY .vimrc ./
```
