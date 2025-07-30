#!/usr/bin/env fish

set -g fish_term24bit 1

function fish_greeting
    clear
end

function fish_preexec --on-event fish_preexec
    set -g __cmd_start_time (date +%s%3N)
    set -g __last_cmd $argv
end

function fish_postexec --on-event fish_postexec
    set -l last_status $status

    switch "$__last_cmd"
        case clear
            return

        case exit
            exit 0
    end

    set -l end_time (date +%s%3N)
    set -l duration_ms (math "$end_time - $__cmd_start_time")
    set -l duration_sec (math "$duration_ms / 1000")

    set_color normal
    set -l msg_text "Executed: $__last_cmd"

    if test $last_status -ne 0
        set_color --dim red
        set msg_text "✘ ($last_status)"
    else
        set_color --dim
        set msg_text "✓"
    end

    if test $duration_ms -lt 1000
        echo "$msg_text $duration_ms"ms
    else
        echo "$msg_text $duration_sec"s
    end

    set_color normal
end


set __dirname (realpath (dirname (status --current-filename)))

source $__dirname/aliases.fish
source $__dirname/functions.fish

if status is-interactive
    source $__dirname/theme.fish
end
