__dirname="$(dirname "$(realpath "$0")")"
# echo "source \"$__dirname/bash-config/__init__.sh\"" >~/.bashrc
# echo "Bash config linked"

echo "source \"$__dirname/fish-config/__init__.fish\"" >~/.config/fish/config.fish
echo "Fish config linked"
