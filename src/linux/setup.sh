#!/bin/bash

# Full native Linux setup. Run as your user: sudo is used only for apt/snap.
# Package priority on Linux: mise (config/mise-config.toml) > snap > apt

set -e

if [ "$EUID" -eq 0 ]; then
	echo "Run 'make linux' without sudo; it asks for your password when needed." >&2
	exit 1
fi

cd "$(dirname "${BASH_SOURCE[0]}")/../.."
export PATH="$HOME/.local/bin:$PATH"

echo "Requesting sudo for apt and snap..."
sudo -v

echo "Installing apt packages..."
# No mise/snap source for a C toolchain; cargo:dotsh and cargo:shaka need it
sudo apt-get update
sudo apt-get install -y build-essential

echo "Installing snap packages..."
sudo snap install code --classic

echo "Installing mise..."
if [ ! -x "$HOME/.local/bin/mise" ]; then
	curl -fsSL https://mise.run | sh
fi

echo "Compiling scripts..."
mise exec go@latest -- go run ./src/compile-scripts/main.go

echo "Linking config files..."
./.local/bin.linux/config-init

echo "Installing mise packages..."
mise exec -- ./.local/bin.linux/psy

echo "Configuring shell..."
mise exec -- bash ./etc/shell-config.sh

echo "Installing FiraCode Nerd Font..."
font_dir="$HOME/.local/share/fonts/FiraCode"
if [ ! -d "$font_dir" ]; then
	tmp_dir="$(mktemp -d)"
	curl -fsSL -o "$tmp_dir/FiraCode.zip" https://github.com/ryanoasis/nerd-fonts/releases/latest/download/FiraCode.zip &&
		mkdir -p "$font_dir" &&
		unzip -q -o "$tmp_dir/FiraCode.zip" -d "$font_dir" &&
		fc-cache -f "$font_dir"
	rm -rf "$tmp_dir"
fi

echo "Configuring GNOME settings..."
gset() {
	gsettings set "$@" 2>/dev/null || echo "Skipped: gsettings set $*"
}

gset org.gnome.desktop.interface color-scheme prefer-dark
gset org.gnome.desktop.interface clock-show-weekday true
gset org.gnome.desktop.interface clock-format 12h
gset org.gnome.desktop.interface enable-hot-corners false

gset org.gnome.desktop.peripherals.keyboard repeat true
gset org.gnome.desktop.peripherals.keyboard repeat-interval 30
gset org.gnome.desktop.peripherals.keyboard delay 375
gset org.gnome.desktop.peripherals.touchpad tap-to-click true

gset org.gnome.desktop.session idle-delay 300
gset org.gnome.settings-daemon.plugins.power sleep-inactive-ac-type nothing
gset org.gnome.desktop.privacy report-technical-problems false
gset org.gnome.desktop.privacy send-software-usage-stats false

gset org.gnome.shell.extensions.dash-to-dock dock-fixed false
gset org.gnome.shell.extensions.dash-to-dock autohide true
gset org.gnome.shell.extensions.dash-to-dock intellihide false
gset org.gnome.shell.extensions.dash-to-dock show-trash false

gset org.gtk.gtk4.Settings.FileChooser show-hidden true
gset org.gtk.Settings.FileChooser show-hidden true
gset org.gnome.nautilus.preferences default-folder-viewer icon-view

echo
echo "Done. Open a new terminal."
