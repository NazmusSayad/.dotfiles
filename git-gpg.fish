#!/usr/bin/env fish

if not command -sq gpg
    echo "Error: GPG not installed"
    exit 1
end

if not command -sq git
    echo "Error: Git not installed"
    exit 1
end

set -l gpg_keys (gpg --list-secret-keys --keyid-format LONG 2>/dev/null | grep sec)
if test -z "$gpg_keys"
    set -l git_email (git config --get user.email)
    if test -z "$git_email"
        set git_email "user@example.com"
        git config --global user.email $git_email
    end
    
    set -l batch_file (mktemp)
    echo "Key-Type: RSA
Key-Length: 4096
Key-Usage: sign
Name-Real: Git User
Name-Email: $git_email
Expire-Date: 0
%no-protection
%commit" > $batch_file
    
    gpg --batch --generate-key $batch_file
    rm $batch_file
    
    if test $status -ne 0
        exit 1
    end
end

set -l git_email (git config --get user.email)
set -l gpg_key_id (gpg --list-secret-keys --keyid-format LONG | grep sec | head -1 | awk '{print $2}' | cut -d'/' -f2)
if test -z "$gpg_key_id"
    exit 1
end

git config --global user.signingkey $gpg_key_id
git config --global commit.gpgsign true
git config --global gpg.program (which gpg)

gpg --armor --export $gpg_key_id