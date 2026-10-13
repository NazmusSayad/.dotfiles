# Ubuntu WSL and sshtty

Run Linux shells at https://wsl.sayad.dev and native Windows shells at https://win.sayad.dev through the existing Windows `Desktop` Cloudflare tunnel. Authentication is managed separately in Cloudflare.

## Setup status

Checked on 2026-10-10:

- Ubuntu and the Linux `sayad` account are configured.
- Node.js `24.21.0`, pnpm `12.10.1`, `@antfu/ni` `30.6.0`, the npm package `sshtty`, and Ubuntu's Starship `1.22.1` package are installed.
- `sshtty.service` is installed, enabled, and running as `sayad`.
- Ubuntu and Windows both return HTTP 200 at `http://127.0.0.1:47474`.
- Both hostname DNS routes point to the existing `Desktop` tunnel, and the ingress config passes validation.
- The Windows cloudflared connector handles both hostnames. Unauthenticated requests redirect to the existing Cloudflare Access login.
- A hidden Windows-side WSL watchdog is running so Ubuntu does not stop when the last terminal closes and automatically relaunches after an unexpected termination.

OpenCode is not part of this setup. The `Desktop` tunnel runs as the native `Cloudflared` Windows service.

## Device setup

- Distro: `Ubuntu-26.04`, running on WSL 2.
- Linux account: `sayad`, with Bash as its login shell.
- Windows dotfiles: `F:\.dotfiles`, available inside Ubuntu at `/mnt/f/.dotfiles`.
- Linux tools: mise, Node.js LTS, pnpm, `@antfu/ni`, `sshtty`, and the Ubuntu Starship package.
- Mise config: `~/.config/mise/config.toml`, linked to the main `../mise-config.toml` used by this dotfiles repository.
- `config-init` detects WSL and creates the portable links from `../symlink.yml`, including OpenCode and CLI configuration.
- Running `config-init` on Windows also runs it inside `Ubuntu-26.04` as `sayad`, keeping both environments updated from one command.
- `~/.bashrc` links to the shared `../shell/bashrc.sh`. On Linux it activates mise, loads environment files, and initializes direnv, Shaka, zoxide, and Starship.
- Shaka uses the shared `../shell/alias.json` config. This provides `oc` for `opencode --standalone`.
- sshtty data: `/home/sayad/.sshtty`. Keep it in Linux, not on a Windows mount.
- sshtty's `shell.command` is `/usr/bin/fish`, so browser tabs open Fish without changing the WSL login shell.
- sshtty listener: `127.0.0.1:47474` inside Ubuntu.
- Native Windows sshtty listener: `127.0.0.1:47831`, with a PowerShell shell.
- sshtty settings: `../sshtty.json`, linked to `~/.sshtty/config.json` on both Windows and Ubuntu. Each side reads its own `shell` entry. Fish runs with `--features no-query-term` so tabs started without a connected browser do not wait for a terminal response.
- Tunnel: `Desktop`, UUID `b61784bd-cd62-431d-9718-720b95f14339`, running on Windows.

The machine already had WSL 2 for Docker Desktop. Leave the `docker-desktop` distro alone. Ubuntu is now the default distro.

## Open Ubuntu

Choose **Ubuntu (WSL)** in Windows Terminal, or run this in PowerShell:

```powershell
wsl -d Ubuntu-26.04 --cd ~
```

The Terminal config disables automatic WSL profiles, so Ubuntu has an explicit profile. The existing default Terminal profile is unchanged.

The `sayad` account was created with a locked password. Set its password privately from PowerShell:

```powershell
wsl -d Ubuntu-26.04 -u root -- passwd sayad
```

This enables password-based `sudo`. Do not put the password in this repository.

Use `~/projects` for Linux development. Windows files are available at `/mnt/c` and `/mnt/f`; Linux files are available in Explorer at `\\wsl.localhost\Ubuntu-26.04\home\sayad`.

## Run sshtty

The `sshtty.service` file runs sshtty as `sayad`, not root. From Ubuntu:

```bash
sudo systemctl start sshtty
sudo systemctl status sshtty --no-pager
journalctl -u sshtty -n 50 --no-pager
```

WSL local address: http://127.0.0.1:47474. Remote address: https://wsl.sayad.dev. Native Windows uses http://127.0.0.1:47831 and https://win.sayad.dev.

The service starts when Ubuntu boots. It does not boot Ubuntu when Windows starts, and systemd services do not keep WSL alive on their own. This was observed during setup: WSL shut down the service after the last session ended.

Start both sshtty servers and the existing WSL watchdog from PowerShell. The Windows tunnel starts separately as a service:

```powershell
& F:\.dotfiles\config\wsl\start.ps1
```

The script reuses an existing watchdog if one is running. The watchdog relaunches Ubuntu within a few seconds if the distro stops, and the enabled sshtty systemd service starts during the new boot. It does not install a Windows startup task. The PC must be awake. The tunnel service starts automatically with Windows.

Restart or stop the service from Ubuntu:

```bash
sudo systemctl restart sshtty
sudo systemctl stop sshtty
```

For a foreground session, stop the service first, then run:

```bash
cd ~
mise exec -- sshtty --host 127.0.0.1 --port 47474
```

## Tunnel routing

`../win-cloudflared.yml` is linked to `%USERPROFILE%\.cloudflared\config.yml` on Windows. The Windows service reads the dotfiles config directly. It must contain both hostname routes, followed by the catch-all:

```yaml
tunnel: Desktop
ingress:
  - hostname: wsl.sayad.dev
    service: http://127.0.0.1:47474
  - hostname: win.sayad.dev
    service: http://127.0.0.1:47831
  - service: http_status:404
```

Windows cloudflared reaches Ubuntu through WSL's localhost forwarding. No router port forwarding or public listener is needed.

Create the hostname's DNS route from Windows if it is missing:

```powershell
cloudflared tunnel route dns Desktop wsl.sayad.dev
cloudflared tunnel route dns Desktop win.sayad.dev
```

`start.ps1` starts both sshtty servers and the Windows-side WSL watchdog. It skips processes that are already running. The existing dotfiles startup launcher calls this script through `launch.jsonc`. It does not launch cloudflared. No cloudflared service runs inside WSL.

The native service is installed using `cloudflared service install` from an administrator terminal. Machine-specific paths and the resolved tunnel UUID belong in the Windows service command line, not the dotfiles config. Cloudflared's installer enables automatic startup and repeatedly restarts the service after failures with a 20-second delay, including non-crash failures. No retry wrapper is needed. The credentials stay outside the repository.

To reload the tunnel after changing `win-cloudflared.yml`, run:

```powershell
sudo powershell.exe -NoProfile -Command "Restart-Service Cloudflared"
```

The connector logs are in `%ProgramData%\cloudflared\logs\cloudflared.log`, with automatic size-based rotation. Service startup and failure events are also recorded in the Windows Application event log under `Cloudflared`.

## Recreate the Linux tool setup

Install the distro from PowerShell, then launch it and create the `sayad` account:

```powershell
wsl --install -d Ubuntu-26.04
wsl --set-default Ubuntu-26.04
```

Inside Ubuntu, install mise and link the repository's main mise config. These commands assume there is no existing mise config to overwrite:

```bash
mkdir -p ~/.cache ~/.config/mise
curl -fsSL https://mise.run -o ~/.cache/mise-install.sh
sh ~/.cache/mise-install.sh
ln -s /mnt/f/.dotfiles ~/.dotfiles
ln -s ~/.dotfiles/config/mise-config.toml ~/.config/mise/config.toml
~/.local/bin/mise trust ~/.config/mise/config.toml
~/.local/bin/mise install
cd ~/.dotfiles
~/.local/bin/mise exec -- go run ./src/scripts/config-init/main.go
```

Install the Linux shell and build dependencies through Ubuntu. This leaves the main mise config unchanged:

```bash
sudo apt update
sudo apt install build-essential direnv fish starship zoxide ripgrep btop ttyd shfmt zip unzip ffmpeg imagemagick git-lfs
```

Install and enable the service inside Ubuntu:

```bash
sudo install -m 644 /mnt/f/.dotfiles/config/wsl/sshtty.service /etc/systemd/system/sshtty.service
sudo systemctl daemon-reload
sudo systemctl enable --now sshtty
```

Update sshtty from Ubuntu, then restart it:

```bash
cd ~
mise upgrade npm:sshtty
sudo systemctl restart sshtty
```

Use `ni`, `nr`, and `nlx` for project package workflows. `npm:sshtty` is a mise package identifier; the installed command is `sshtty`.

## Troubleshooting

- Windows localhost fails: check `systemctl status sshtty` inside Ubuntu first, then `curl -I http://127.0.0.1:47474` in both environments.
- Localhost works but the hostname returns 404: the connector may still be using the old ingress configuration.
- The hostname returns 502: cloudflared cannot reach sshtty. Check that Ubuntu is running and sshtty is listening.
- The hostname returns 525: check the DNS route. The hostname may be reaching a different origin instead of this tunnel.
- Keep the Cloudflare Access applications for `wsl.sayad.dev` and `win.sayad.dev` separate to avoid cross-hostname login redirects.
- Shells survive sshtty server crashes according to the package documentation. They do not survive Windows shutdown or `wsl --shutdown`.
- Avoid `wsl --shutdown` during normal work; it also stops Docker's WSL environment.

## References

- [sshtty usage and data location](https://github.com/NazmusSayad/sshtty#readme)
- [WSL networking and localhost forwarding](https://learn.microsoft.com/en-us/windows/wsl/networking)
- [WSL systemd behavior](https://learn.microsoft.com/en-us/windows/wsl/systemd)
- [mise npm backend](https://mise.jdx.dev/dev-tools/backends/npm.html)
