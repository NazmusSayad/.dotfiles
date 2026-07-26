#!/bin/bash

echo "Setting up fish shell colors..."
fish -c 'set -U fish_greeting'
fish -c 'set -U fish_color_end normal'
fish -c 'set -U fish_color_quote green'
fish -c 'set -U fish_color_comment --dim'
fish -c 'set -U fish_color_command magenta'
