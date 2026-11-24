#!/bin/bash

# Package priority on Linux: mise (config/mise-config.toml) > snap > apt

echo "Installing apt packages..."
# No mise/snap source for a C toolchain; cargo:dotsh and cargo:shaka need it
sudo apt-get update
sudo apt-get install -y build-essential

echo "Installing snap packages..."
sudo snap install code --classic

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
echo "Next: ./.local/bin.linux/config-init && ./.local/bin.linux/psy && make shell"
