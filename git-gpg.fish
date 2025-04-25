#!/usr/bin/env fish

printf '\e[8;45;80t'

if not command -sq gpg
    set_color red
    echo "Error: GPG not installed"
    set_color normal
    exit 1
end

if not command -sq git
    set_color red
    echo "Error: Git not installed"
    set_color normal
    exit 1
end

set -l git_name (git config --get user.name)
set -l git_email (git config --get user.email)

if test -z "$git_email" -o -z "$git_name"
    echo "Error: Git user.email or user.name not configured"
    echo "Please run: git config --global user.name \"Your Name\""
    echo "Please run: git config --global user.email \"your@email.com\""
    exit 1
end

echo -n "Git user name  : "
set_color blue
echo "$git_name"
set_color normal

echo -n "Git user email : "
set_color blue
echo "$git_email"
set_color normal

set -l gpg_keys (gpg --list-secret-keys --keyid-format LONG 2>/dev/null | grep sec)
if test -z "$gpg_keys"
    set -l batch_file (mktemp)
    echo "Key-Type: RSA
Key-Length: 4096
Key-Usage: sign
Name-Real: $git_name
Name-Email: $git_email
Expire-Date: 0
%no-protection
%commit" >$batch_file

    gpg --batch --generate-key $batch_file
    rm $batch_file

    if test $status -ne 0
        exit 1
    end
end

set -l gpg_key_id (gpg --list-secret-keys --keyid-format LONG | grep sec | head -1 | awk '{print $2}' | cut -d'/' -f2)
if test -z "$gpg_key_id"
    exit 1
end

git config --global user.signingkey $gpg_key_id
git config --global commit.gpgsign true
git config --global gpg.program (which gpg)

echo -n "GPG key ID     : "
set_color blue
echo "$gpg_key_id"
set_color normal

set_color green
echo "GPG key generated and configured for Git."
set_color normal

echo ""
echo ""
set_color yellow
gpg --armor --export $gpg_key_id
set_color normal
echo ""
echo ""

echo -n "Press any key to exit..."
read -n 1 -p "" key
