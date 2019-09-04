# findTorrent


## Install

### Install Docker 

sudo apt install apt-transport-https ca-certificates curl software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"

sudo apt update

sudo apt install docker-ce

REF: https://www.digitalocean.com/community/tutorials/como-instalar-y-usar-docker-en-ubuntu-18-04-1-es


### Install Golang

sudo apt-get update

sudo apt-get -y upgrade

wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz #Check latest in https://golang.org/dl/

sudo tar -xvf go1.12.9.linux-amd64.tar.gz

sudo mv go /usr/local

mkdir ~/work

export GOROOT=/usr/local/go

export GOPATH=$HOME/work

export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

 
REF: https://tecadmin.net/install-go-on-ubuntu/



### Install bntoolkit

#### From github (if it's public)


#### From local file

cp -r BNToolkit /home/user/work/src/github.com/RaulCalvoLaorden/bntoolkit

cd /home/user/work/src/github.com/RaulCalvoLaorden/bntoolkit

go get .

OR

bash installLocal.sh
'''
#!/bin/bash

mkdir -p ~/work/src/github.com

cp -r * ~/work/src/github.com/

cd ~/work/src/github.com/RaulCalvoLaorden/bntoolkit/

ls

go get .

go install

bntoolkit

'''

- 
