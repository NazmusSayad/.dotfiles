#!/bin/bash

source ./config.sh

echo "Installing npm packges globally"
powershell -command "npm i -g npm yarn nodemon ts-node live-server gitignore code-info asem eas-cli"