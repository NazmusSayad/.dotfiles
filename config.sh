#!/bin/bash

cat ./.bashrc > ~/.bashrc
cat ./WindowsTerminal.json > ~/AppData/Local/Microsoft/Windows\ Terminal/settings.json

git config --global core.excludesfile ~/default.gitignore
echo "Default gitignore added"

cp ./default.gitignore ~/default.gitignore
git config --global core.excludesfile ~/default.gitignore
echo "Default gitignore added"

git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
echo "Git auth info added"

git config --global init.defaultBranch master
git config --global --add --bool push.autoSetupRemote true
echo "Git config added..."

npm config set ignore-scripts true
echo "Npm config added..."