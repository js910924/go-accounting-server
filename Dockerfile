FROM ubuntu:latest
MAINTAINER Jason js910924@gmail.com
LABEL description="An accounting server" version="1.0" owner="Jason Chen"
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install net-tools -y
RUN apt-get install vim -y
RUN apt-get install curl -y
RUN apt-get install tmux -y
RUN apt-get install mysql-server -y
RUN apt-get install wget -y
RUN wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
RUN tar -xvf go1.13.4.linux-amd64.tar.gz
RUN mkdir ~/go ~/go/src ~/go/pkg ~/go/bin
ENV GOROOT /go
ENV GOPATH /root/go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH
COPY ./ ./root/go/src/account-server
RUN mv ./root/go/src/account-server/.vimrc ./root/.vimrc
RUN service mysql stop
RUN usermod -d /var/lib/mysql/ mysql
RUN ln -s /var/lib/mysql/mysql.sock /tmp/mysql.sock
RUN chown -R mysql:mysql /var/lib/mysql