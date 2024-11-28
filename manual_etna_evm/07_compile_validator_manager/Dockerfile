FROM golang:1.22.9-bookworm

# Install git and other necessary dependencies
RUN apt-get update
RUN apt-get install -y git jq

# Download and run the install_foundry.sh script
ARG TELEPORTER_COMMIT
RUN curl -o install_foundry.sh https://raw.githubusercontent.com/ava-labs/teleporter/${TELEPORTER_COMMIT}/scripts/install_foundry.sh && \
    chmod +x install_foundry.sh && \
    ./install_foundry.sh && \
    rm install_foundry.sh

RUN go install github.com/ava-labs/subnet-evm/cmd/abigen@v0.6.12

ENV PATH="/root/.foundry/bin/:${PATH}"

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
