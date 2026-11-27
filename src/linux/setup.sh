#!/bin/bash

# Native Linux system setup. Run as your user: sudo is used only where needed.
# Run after compile, config-init and psy, like `make macos`.
# Package priority on Linux: mise (config/mise-config.toml) > snap > apt

set -e

if [ "$EUID" -eq 0 ]; then
	echo "Run 'make linux' without sudo; it asks for your password when needed." >&2
	exit 1
fi

# Adds an item to a gsettings string list, like enabled-extensions, unless it is already there
gsettings_list_add() {
	local list
	list="$(gsettings get "$1" "$2")"
	case "$list" in
	*"'$3'"*) ;;
	"@as []") gsettings set "$1" "$2" "['$3']" ;;
	*) gsettings set "$1" "$2" "${list%]}, '$3']" ;;
	esac
}

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

echo "Starting OpenLogi agent..."
if systemctl --user cat openlogi-agent.service >/dev/null 2>&1; then
	systemctl --user enable --now openlogi-agent.service
	# GNOME extension that tells OpenLogi the focused app, for per-app buttons
	openlogi_ext=openlogi-frontmost@openlogi.dev
	openlogi_ext_dir="$HOME/.local/share/gnome-shell/extensions/$openlogi_ext"
	openlogi_ext_url="https://raw.githubusercontent.com/AprilNEA/OpenLogi/v$(dpkg-query -W -f='${Version}' openlogi)/crates/openlogi-hook/gnome-shell-extension/$openlogi_ext"
	mkdir -p "$openlogi_ext_dir"
	for file in metadata.json extension.js ringAnimation.js; do
		curl -fsSL -o "$openlogi_ext_dir/$file" "$openlogi_ext_url/$file"
	done
	# A new extension is only found after logging in again, so add it to the enabled list directly
	gnome-extensions enable "$openlogi_ext" 2>/dev/null || gsettings_list_add org.gnome.shell enabled-extensions "$openlogi_ext"
else
	echo "Skipped: OpenLogi not installed (run 'psy deb' first)"
fi

echo "Starting kanata keyboard remapper..."
# Keychron keys: Caps Lock (F20) as Alt+Tab, Print Screen key (F13) as Print (config/kanata)
if [ -f "$HOME/.config/systemd/user/kanata.service" ]; then
	systemctl --user daemon-reload
	systemctl --user enable --now kanata.service
else
	echo "Skipped: kanata service not linked (run config-init first)"
fi

echo "Starting Vicinae launcher..."
if command -v vicinae >/dev/null; then
	systemctl --user enable --now vicinae.service
	# GNOME extension that lets Vicinae see the clipboard (without it, Clipboard History says "Clipboard monitoring unavailable")
	vicinae_ext=vicinae@dagimg-dot
	vicinae_ext_url="$(curl -fsSL https://api.github.com/repos/vicinaehq/gnome-extension/releases/latest | jq -r '.assets[] | select(.name | endswith(".zip")) | .browser_download_url' | head -1)"
	vicinae_ext_zip="$(mktemp --suffix=.zip)"
	curl -fsSL -o "$vicinae_ext_zip" "$vicinae_ext_url"
	gnome-extensions install --force "$vicinae_ext_zip"
	rm -f "$vicinae_ext_zip"
	gnome-extensions enable "$vicinae_ext" 2>/dev/null || gsettings_list_add org.gnome.shell enabled-extensions "$vicinae_ext"
else
	echo "Skipped: Vicinae not installed"
fi

echo "Starting Slack manager timer..."
if [ -f "$HOME/.config/systemd/user/slack-manager.timer" ]; then
	systemctl --user daemon-reload
	systemctl --user enable --now slack-manager.timer
else
	echo "Skipped: slack-manager timer not linked (run config-init first)"
fi

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

# Adds a GNOME keyboard shortcut: id, name, command, key
gnome_shortcut() {
	local path=/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/$1/
	gset "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$path" name "$2"
	gset "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$path" command "$3"
	gset "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$path" binding "$4"
	gsettings_list_add org.gnome.settings-daemon.plugins.media-keys custom-keybindings "$path" 2>/dev/null || echo "Skipped: shortcut $2"
}

# Vicinae can't register global shortcuts on GNOME, so GNOME shortcuts open it. Super is the physical Ctrl key, since Ctrl and Win are swapped
# Only one keyboard layout, so Super+Space doesn't need to switch layouts
gset org.gnome.desktop.wm.keybindings switch-input-source "['XF86Keyboard']"
gnome_shortcut vicinae Vicinae "vicinae toggle" "<Super>space"
# Clipboard history on Super+V like Windows' Win+V; the notification list keeps Super+M
gset org.gnome.shell.keybindings toggle-message-tray "['<Super>m']"
gnome_shortcut vicinae-clipboard "Vicinae Clipboard History" "vicinae deeplink vicinae://launch/clipboard/history" "<Super>v"

echo "Configuring computer-use-linux..."
# Turns on the AT-SPI accessibility bridge and installs the GNOME window-control extension
mise exec -- computer-use-linux setup >/dev/null || echo "Skipped: computer-use-linux setup"
mise exec -- computer-use-linux setup-window-targeting >/dev/null || echo "Skipped: computer-use-linux setup-window-targeting"

echo
echo "Done. Log out and back in, then open a new terminal."
