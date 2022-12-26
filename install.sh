#!/bin/sh

cp ./lib/.gitignore ~/.gitignore
cat ./scripts/main.sh > ~/.bashrc

git config --global core.excludesfile ~/.gitignore

mv ./config ./config.sh
source config.sh
mv ./config.sh ./config
