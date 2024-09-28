#!/usr/bin/env bash

__dirname="$(dirname "$(realpath "$0")")"
echo "source $__dirname/bash-config/__init__.sh" >~/.bashrc
echo "Bash config linked"

echo "source $__dirname/fish-config/__init__.fish" >~/.config/fish/config.fish
echo "Fish config linked"

cp ./lib/default.gitignore ~/default.gitignore -f
echo "Default gitignore added"

git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
git config --global core.autocrlf false
git config --global init.defaultBranch main
git config --global --add safe.directory '*'
git config --global core.excludesfile ~/default.gitignore
git config --global --add --bool push.autoSetupRemote true
echo "Git config added"

npm config set ignore-scripts true
echo "Npm config added"
