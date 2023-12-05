#!/bin/bash

source ./config.sh

echo "Installing npm packges globally"
powershell -command "npm i -g npm nodemon ts-node live-server gitignore code-info netserv asem"