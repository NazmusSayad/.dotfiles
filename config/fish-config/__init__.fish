#!/usr/bin/env fish

set -g fish_term24bit 1

function fish_greeting
    clear
end

function fish_preexec --on-event fish_preexec
    set -g __cmd_start_time (date +%s%N)
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

    set -l end_time (date +%s%N)
    set -l duration_ns (math "$end_time - $__cmd_start_time")
    set -l duration_text

    if test $duration_ns -lt 1000000
        set duration_text "$duration_ns"ns
    else if test $duration_ns -lt 1000000000
        set -l duration_ms (printf "%.2f" (math "$duration_ns / 1000000"))
        set duration_text "$duration_ms"ms
    else
        set -l duration_sec (printf "%.2f" (math "$duration_ns / 1000000000"))
        set duration_text "$duration_sec"s
    end

    set_color normal

    if test $last_status -ne 0
        set_color --dim red
        echo "✘ $duration_text ($last_status)"
    else
        set_color --dim
        echo "✓ $duration_text (0)"
    end

    set_color normal
end

set __dirname (realpath (dirname (status --current-filename)))

source $__dirname/aliases.fish
source $__dirname/functions.fish

if status is-interactive
    source $__dirname/theme.fish
end
