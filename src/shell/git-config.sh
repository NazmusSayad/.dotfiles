#!/bin/bash

git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
git config --global init.defaultBranch main
git config --global --add safe.directory "*"
git config --global --add --bool push.autoSetupRemote true
git config --global core.eol lf
git config --global core.autocrlf false
git config --global core.pager cat
git config --system core.longpaths true
git config --global core.ignorecase false
git config --global core.editor "code --wait"

git config --global credential.helper manager
git config --global core.sshCommand "C:/Windows/System32/OpenSSH/ssh.exe"
