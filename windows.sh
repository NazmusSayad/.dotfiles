#!/bin/bash

cat ./WindowsTerminal.json > ~/AppData/Local/Microsoft/Windows\ Terminal/settings.json
source ./config.sh

echo "Installing npm packges globally"
powershell -command "npm i -g npm nodemon ts-node live-server gitignore code-info netserv asem"