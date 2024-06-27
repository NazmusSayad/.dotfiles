#!/bin/bash

cp ./.bashrc ~/.bashrc -rf
echo "Bashrc added..."

cp ./default.gitignore ~/default.gitignore -rf
echo "Default gitignore added..."

git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
git config --global core.autocrlf false
git config --global init.defaultBranch main
git config --global --add safe.directory '*'
git config --global core.excludesfile ~/default.gitignore
git config --global --add --bool push.autoSetupRemote true
echo "Git config added..."

npm config set ignore-scripts true
echo "Npm config added..."

echo "Installing npm packges globally"
npm i -g npm nodemon ts-node ts-node-dev gitignore code-info netserv npm-run-all concurrently
