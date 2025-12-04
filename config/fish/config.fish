#!/usr/bin/env fish

if status is-interactive
    starship init fish | source
end

function fish_greeting
    if not set -q TERM_PROGRAM
        fastfetch
    end
end

function ___pre-exec___ --on-event fish_preexec
    set -g __cmd_start_time (date +%s%N)
end

function ___post-exec___ --on-event fish_postexec
    set -l last_status $status

    switch "$argv"
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

source ~/.dotfiles/config/fish/aliases.fish
source ~/.dotfiles/config/fish/functions.fish
