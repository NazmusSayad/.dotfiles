#!/bin/bash

# Native Linux system setup. Run as your user: sudo is used only where needed.
# Run after compile, config-init and psy, like `make macos`.
# Package priority on Linux: mise (config/mise-config.toml) > snap > apt

set -e

if [ "$EUID" -eq 0 ]; then
	echo "Run 'make linux' without sudo; it asks for your password when needed." >&2
	exit 1
fi

cd "$(dirname "${BASH_SOURCE[0]}")/../.."
export PATH="$HOME/.local/bin:$PATH"

echo "Installing mise..."
if [ ! -x "$HOME/.local/bin/mise" ]; then
	curl -fsSL https://mise.run | sh
fi

echo "Requesting sudo for apt, SSH server and tunnel startup..."
sudo -v

echo "Installing apt packages..."
# No mise/snap source: C toolchain (for cargo:dotsh/shaka) and SSH server
sudo apt-get update
sudo apt-get install -y build-essential openssh-server

echo "Configuring SSH server..."
if [ -f ./.dotsys/ssh/sshd-linux.conf ]; then
	sudo install -m 644 ./.dotsys/ssh/sshd-linux.conf /etc/ssh/sshd_config.d/10-dotfiles.conf
	sudo systemctl daemon-reload
	sudo systemctl enable --now ssh.socket
	sudo systemctl restart ssh.socket
fi

# ssh rejects a group-writable config, and Ubuntu's umask makes checkouts 664
[ -f ./.dotsys/ssh/ssh-linux ] && chmod 644 ./.dotsys/ssh/ssh-linux

echo "Installing SSH authorized keys..."
if [ -f ./.dotsys/ssh/authorized_keys ]; then
	mkdir -p -m 700 "$HOME/.ssh"
	install -m 600 ./.dotsys/ssh/authorized_keys "$HOME/.ssh/authorized_keys"
fi

echo "Configuring Cloudflare tunnel..."
cf() { mise exec -- cloudflared "$@"; }
tunnel_name=Linux
if [ ! -f "$HOME/.cloudflared/cert.pem" ]; then
	cf tunnel login
fi
if [ -z "$(cf tunnel list --name "$tunnel_name" -o json | jq -r '.[0].id // empty')" ]; then
	cf tunnel create "$tunnel_name"
fi
tunnel_id="$(cf tunnel list --name "$tunnel_name" -o json | jq -r '.[0].id')"
if [ ! -f "$HOME/.cloudflared/$tunnel_id.json" ]; then
	cf tunnel token --cred-file "$HOME/.cloudflared/$tunnel_id.json" "$tunnel_name"
fi
cf tunnel route dns "$tunnel_name" linux.sayad.cc || true
sudo loginctl enable-linger "$USER"
systemctl --user daemon-reload
systemctl --user enable --now cloudflared
systemctl --user restart cloudflared

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
gset org.gnome.desktop.input-sources xkb-options "['ctrl:swap_lwin_lctl', 'ctrl:swap_rwin_rctl']"
gset org.gnome.desktop.peripherals.touchpad tap-to-click true

gset org.gnome.desktop.session idle-delay 300
gset org.gnome.settings-daemon.plugins.power sleep-inactive-ac-type nothing
gset org.gnome.desktop.privacy report-technical-problems false
gset org.gnome.desktop.privacy send-software-usage-stats false
gset org.gnome.desktop.sound event-sounds false

gset org.gnome.shell.extensions.dash-to-dock dock-fixed false
gset org.gnome.shell.extensions.dash-to-dock autohide true
gset org.gnome.shell.extensions.dash-to-dock intellihide false
gset org.gnome.shell.extensions.dash-to-dock show-trash false

gset org.gtk.gtk4.Settings.FileChooser show-hidden true
gset org.gtk.Settings.FileChooser show-hidden true
gset org.gnome.nautilus.preferences default-folder-viewer icon-view

echo "Configuring computer-use-linux..."
# Turns on the AT-SPI accessibility bridge and installs the GNOME window-control extension
mise exec -- computer-use-linux setup >/dev/null || echo "Skipped: computer-use-linux setup"
mise exec -- computer-use-linux setup-window-targeting >/dev/null || echo "Skipped: computer-use-linux setup-window-targeting"

echo
echo "Done. Log out and back in, then open a new terminal."
