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

echo "Adding $USER to the docker group..."
if getent group docker >/dev/null; then
	sudo usermod -aG docker "$USER"
else
	echo "Skipped: docker group not found (run 'psy apt' first)"
fi

echo "Mounting Windows D: drive at /mnt/d..."
d_drive="UUID=01DC2CA9456E1650"
sudo mkdir -p /mnt/d
if ! grep -q "^$d_drive " /etc/fstab; then
	echo "$d_drive /mnt/d ntfs3 uid=$(id -u),gid=$(id -g),nofail 0 0" | sudo tee -a /etc/fstab >/dev/null
	sudo systemctl daemon-reload
fi
mountpoint -q /mnt/d || sudo mount /mnt/d || echo "Skipped: mount /mnt/d (turn off Windows Fast Startup if the drive is locked)"

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

echo "Configuring Android Studio and emulator..."
# kvm group: hardware-accelerated emulator; inotify: Android Studio file watching
sudo usermod -aG kvm "$USER"
echo 'fs.inotify.max_user_watches = 524288' | sudo tee /etc/sysctl.d/60-android-studio.conf >/dev/null
sudo sysctl -p /etc/sysctl.d/60-android-studio.conf >/dev/null
# React Native DevTools (Electron) needs user namespaces for its sandbox
sudo tee /etc/apparmor.d/react-native-devtools >/dev/null <<'EOF'
abi <abi/5.0>,
include <tunables/global>

profile react-native-devtools "@{HOME}/.cache/dotslash/**/React Native DevTools" flags=(unconfined) {
  userns,
}
EOF
sudo apparmor_parser -r /etc/apparmor.d/react-native-devtools

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

gset org.gnome.desktop.session idle-delay 900
gset org.gnome.desktop.screensaver lock-enabled false
gset org.gnome.desktop.screensaver ubuntu-lock-on-suspend false
gset org.gnome.desktop.lockdown disable-lock-screen true
gset org.gnome.settings-daemon.plugins.power sleep-inactive-ac-type nothing
gset org.gnome.settings-daemon.plugins.power sleep-inactive-battery-type nothing
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

gset io.github.celluloid-player.Celluloid always-use-floating-controls false
gset io.github.celluloid-player.Celluloid autofit-enable false
gset io.github.celluloid-player.Celluloid draggable-video-area-enable true
gset io.github.celluloid-player.Celluloid mpv-config-enable true
gset io.github.celluloid-player.Celluloid mpv-config-file "file://$HOME/.config/mpv/mpv.conf"

echo "Configuring computer-use-linux..."
# Turns on the AT-SPI accessibility bridge and installs the GNOME window-control extension
mise exec -- computer-use-linux setup >/dev/null || echo "Skipped: computer-use-linux setup"
mise exec -- computer-use-linux setup-window-targeting >/dev/null || echo "Skipped: computer-use-linux setup-window-targeting"

echo
echo "Done. Log out and back in, then open a new terminal."
