#!/usr/bin/env fish

set -g fish_term24bit 1

function fish_greeting
    echo -n "⏰ "

    set_color yellow

    echo -n (date "+%A, %d %B %Y")
    echo -n " - "
    echo (date "+%I:%M %p")

    set_color normal
end

function fish_preexec --on-event fish_preexec
    set -g __cmd_start_time (date +%s%3N)
    set -g __last_cmd $argv
end

function fish_postexec --on-event fish_postexec
    if test "$__last_cmd" = exit
        exit 0
        return
    end

    set -l last_status $status
    set -l end_time (date +%s%3N)
    set -l duration_ms (math "$end_time - $__cmd_start_time")
    set -l duration_sec (math "$duration_ms / 1000")

    if test "$__last_cmd" = clear
        return
    end

    set_color normal
    if test $last_status -ne 0
        set_color --dim red
        echo "✘ Code: $last_status"
    else
        set_color --dim
        if test $duration_ms -lt 1000
            echo "✓ Took: $duration_ms"ms
        else
            echo "✓ Took: $duration_sec"s
        end

    end
    set_color normal
end


set __dirname (realpath (dirname (status --current-filename)))

source $__dirname/aliases.fish
source $__dirname/functions.fish

if status is-interactive
    source $__dirname/theme.fish
end
