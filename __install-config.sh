#!/bin/bash

echo "Configuring git..."
git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
git config --global init.defaultBranch main
git config --global --add safe.directory "*"
git config --global --add --bool push.autoSetupRemote true
git config --global core.eol lf
git config --global core.autocrlf false
git config --global core.pager cat
git config --global pull.rebase false
git config --system core.longpaths true
git config --global core.ignorecase false
git config --global core.editor "code --wait"
git config --global core.excludesfile ~/.gitignore

echo "Configuring fish..."
fish -c 'fish_config theme choose default'
fish -c 'set -U fish_greeting'
fish -c 'set -U fish_color_end normal'
fish -c 'set -U fish_color_quote green'
fish -c 'set -U fish_color_comment --dim'
fish -c 'set -U fish_color_command magenta'
