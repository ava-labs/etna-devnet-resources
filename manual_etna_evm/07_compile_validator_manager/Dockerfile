FROM ubuntu:24.04

# Install git and other necessary dependencies
RUN apt-get update && \
    apt-get install -y git curl build-essential jq && \
    rm -rf /var/lib/apt/lists/*

# Download and run the install_foundry.sh script
ARG TELEPORTER_COMMIT
RUN curl -o install_foundry.sh https://raw.githubusercontent.com/ava-labs/teleporter/${TELEPORTER_COMMIT}/scripts/install_foundry.sh && \
    chmod +x install_foundry.sh && \
    ./install_foundry.sh && \
    rm install_foundry.sh

ENV PATH="/root/.foundry/bin/:${PATH}"

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
