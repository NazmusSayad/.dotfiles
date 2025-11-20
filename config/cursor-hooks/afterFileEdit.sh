#!/bin/sh

# Play a sound after file edit
if command -v paplay >/dev/null 2>&1; then
    paplay /usr/share/sounds/freedesktop/stereo/message.oga 2>/dev/null &
elif command -v afplay >/dev/null 2>&1; then
    afplay /System/Library/Sounds/Glass.aiff 2>/dev/null &
elif command -v powershell.exe >/dev/null 2>&1; then
    powershell.exe -c "(New-Object Media.SoundPlayer 'C:\Windows\Media\tada.wav').PlaySync();" 2>/dev/null &
elif command -v speaker-test >/dev/null 2>&1; then
    speaker-test -t sine -f 800 -l 1 >/dev/null 2>&1 &
else
    echo "File edited successfully!" | espeak-ng -s 150 2>/dev/null || echo -e "\a" || true
fi