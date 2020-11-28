.DEFAULT_GOAL := build_bmark

all: install_docker install_go build_bmark

# Builds bmark cli, located in the benchmark folder
build_bmark:
	cd ./benchmark && go build -o ../bin/bmark .
	./bin/bmark help

# Installs latest Go version on Ubuntu
install_go:
	sudo apt install golang
	go version

# Installs latest Docker version on Ubuntu 20.04
install_docker:
	sudo apt update
	sudo apt install apt-transport-https ca-certificates curl software-properties-common

	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
	sudo apt update

	apt-cache policy docker-ce
	sudo apt install docker-ce
	sudo systemctl status docker