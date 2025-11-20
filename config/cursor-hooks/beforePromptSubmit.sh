#!/bin/sh
set -e
# Read JSON payload from stdin and log it
payload=$(cat)
echo "$(date -u +"%Y-%m-%dT%H:%M:%SZ") beforeSubmitPrompt" >> /tmp/agent.log
echo "$payload" >> /tmp/agent.log
# Return an empty JSON object (valid response)
echo "{}"
