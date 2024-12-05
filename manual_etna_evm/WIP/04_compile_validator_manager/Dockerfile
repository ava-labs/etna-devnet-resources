FROM ubuntu:24.04

# Install git and other necessary dependencies
RUN apt-get update
RUN apt-get install -y git jq software-properties-common golang

# Set GOPATH and add Go binaries to PATH
ENV GOPATH /root/go
ENV PATH $GOPATH/bin:$PATH

RUN go install github.com/ava-labs/subnet-evm/cmd/abigen@v0.6.12

RUN apt-get install -y wget curl

RUN wget https://github.com/ethereum/solidity/releases/download/v0.8.25/solc-static-linux -O /usr/local/bin/solc
RUN chmod +x /usr/local/bin/solc

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
