#!/bin/bash
# gossip_setup.sh
# adapted from https://www.mindbowser.com/deploying-go-application-on-aws-ec2-server/

sudo apt-get update
sudo apt-get -y upgrade

install_path='/usr/local/'
workspace="${HOME}/go/"
go_tar='go1.14.linux-amd64.tar.gz'

sudo curl https://dl.google.com/go/go1.14.linux-amd64.tar.gz -o ${go_tar}
cksum='08df79b46b0adf498ea9f320a0f23d6ec59e9003660b4c9c1ce8e5e2c6f823ca'
if ! echo "$cksum ${go_tar}" | sha256sum -c -; then
    echo "checksum failed" >&2
    exit 1
fi
sudo tar -xvf ${go_tar}
sudo mv go "${install_path}"

# append [install path / workspace path / bin directory of both] to bash profile
echo "\n\n# golang install" | sudo tee -a ~/.profile
echo "export GOROOT=${install_path}/go" | sudo tee -a ~/.profile
echo "export GOPATH=${workspace}" | sudo tee -a ~/.profile
echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' | sudo tee -a ~/.profile

if [ "$(go version)" = 'go version go1.14 linux/amd64' ]; then
    echo "'$(go version)' installed correctly"
else
    echo "golang failed to install, exiting"
    exit 1 >&2
fi

# golang is now installed and set up, let's move onto the next step

git clone 