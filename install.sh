#!/usr/bin/env bash

# Initialize
__dirname="$(dirname "$(realpath "$0")")"
echo "CWD: $__dirname"
echo ""

# Setting up fish config
mkdir -p ~/.config/fish
echo "source \"$__dirname/fish-config/__init__.fish\"" >~/.config/fish/config.fish
echo "Fish config linked"
echo ""

git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"

git config --global core.eol lf
git config --global core.autocrlf false

git config --global core.pager cat
git config --global core.ignorecase false

git config --global init.defaultBranch main
git config --global core.excludesfile ~/default.gitignore

git config --global --add safe.directory '*'
git config --global --add --bool push.autoSetupRemote true

echo "Git config added"

# Setting up npm config
npm config set ignore-scripts true
echo "Npm config added"

echo ""
echo "Press any key to continue..."
read -n 1
