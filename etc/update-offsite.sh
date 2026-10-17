#!/bin/bash

if [ -z "$1" ]; then
	echo "Usage: $0 <version>"
	exit 1
fi

target=$1
cd ~

mac() { mise "$@"; }
win() { ssh -o BatchMode=yes -o ConnectTimeout=8 192.168.0.102 "mise $*"; }
wsl() { ssh -o BatchMode=yes -o ConnectTimeout=8 192.168.0.102 "wsl.exe -d Ubuntu-26.04 -u sayad --cd ~ -- /home/sayad/.local/bin/mise $*"; }

echo "Updating offsite to $target on mac, win and wsl..."
pending="mac win wsl"
while [ -n "$pending" ]; do
	remaining=""
	for host in $pending; do
		$host upgrade npm:offsite@$target >/dev/null 2>&1
		current=$($host current npm:offsite 2>/dev/null | tr -d '\r')
		if [ "$current" = "$target" ]; then
			echo "$(date +%T) $host: offsite $target"
		else
			remaining="$remaining $host"
		fi
	done
	pending=$remaining
	[ -n "$pending" ] && sleep 10
done
echo "$(date +%T) done"
