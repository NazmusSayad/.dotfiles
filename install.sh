#!/bin/sh

cat ./scripts/.bashrc > ~/.bashrc

cp ./lib/.gitignore ~/.gitignore
git config --global core.excludesfile ~/.gitignore

mv ./config ./config.sh
source ./config.sh
mv ./config.sh ./config
