#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

# Just call the launch nodes step again
${SCRIPT_DIR}/../07_launch_nodes/launch.sh
