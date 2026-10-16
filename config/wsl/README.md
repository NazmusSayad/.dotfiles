# Ubuntu WSL and offsite

Run offsite, the Claude Code web UI, inside Ubuntu at https://wsl.sayad.cc through the existing Windows `Desktop` Cloudflare tunnel, and natively on Windows at https://win.sayad.cc. The sshtty browser terminal runs alongside it, at https://sshwsl.sayad.cc in Ubuntu and https://sshwin.sayad.cc on Windows. Any other local port opens at `https://<port>.sayad.cc`. Authentication is managed separately in Cloudflare.

## Setup status

Checked on 2026-10-10:

- Ubuntu and the Linux `sayad` account are configured.
- Node.js `24.21.0`, pnpm `12.10.1`, `@antfu/ni` `30.6.0`, and Ubuntu's Starship `1.22.1` package are installed.
- Both hostname DNS routes point to the existing `Desktop` tunnel, and the ingress config passes validation.
- The Windows cloudflared connector handles both hostnames. Unauthenticated requests redirect to the existing Cloudflare Access login.
- A hidden Windows-side WSL watchdog is running so Ubuntu does not stop when the last terminal closes and automatically relaunches after an unexpected termination.

OpenCode is not part of this setup. The `Desktop` tunnel runs as the native `Cloudflared` Windows service.

## Device setup

- Distro: `Ubuntu-26.04`, running on WSL 2.
- Linux account: `sayad`, with Bash as its login shell.
- Windows dotfiles: `F:\.dotfiles`, available inside Ubuntu at `/mnt/f/.dotfiles`.
- Linux tools: mise, Node.js LTS, pnpm, `@antfu/ni`, Claude Code, `offsite`, and the Ubuntu Starship package.
- Mise config: `~/.config/mise/config.toml`, linked to the main `../mise-config.toml` used by this dotfiles repository.
- `config-init` detects WSL and creates the portable links from `../symlink.yml`, including OpenCode and CLI configuration.
- Running `config-init` on Windows also runs it inside `Ubuntu-26.04` as `sayad`, keeping both environments updated from one command.
- `~/.bashrc` links to the shared `../shell/bashrc.sh`. On Linux it activates mise, loads environment files, and initializes direnv, Shaka, zoxide, and Starship.
- Shaka uses the shared `../shell/alias.json` config. This provides `oc` for `opencode --standalone`.
- offsite data: `/home/sayad/.offsite`. Claude Code data: `/home/sayad/.claude`. Keep both in Linux, not on a Windows mount.
- offsite listener: `127.0.0.1:10101` inside Ubuntu, the port the tunnel already sends `wsl.sayad.cc` to.
- sshtty listeners: `127.0.0.1:10102` inside Ubuntu for `sshwsl.sayad.cc`, and `127.0.0.1:20203` on Windows for `sshwin.sayad.cc`.
- offsite and sshtty have no login of their own. Cloudflare Access is the only thing protecting their hostnames, so keep it on and keep both on `127.0.0.1`.
- Windows offsite listener: `127.0.0.1:20202` for `win.sayad.cc`. It uses the Windows `C:\Users\sayad\.offsite` and `C:\Users\sayad\.claude`, separate from Ubuntu's.
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

## Run offsite and sshtty

The `offsite.service` and `sshtty.service` files run offsite and sshtty as `sayad`, not root. From Ubuntu:

```bash
sudo systemctl start offsite sshtty
sudo systemctl status offsite sshtty --no-pager
journalctl -u offsite -n 50 --no-pager
journalctl -u sshtty -n 50 --no-pager
```

offsite: http://127.0.0.1:10101 in WSL, https://wsl.sayad.cc remotely. sshtty: http://127.0.0.1:10102 in WSL, https://sshwsl.sayad.cc remotely.

The services start when Ubuntu boots. It does not boot Ubuntu when Windows starts, and systemd services do not keep WSL alive on their own. This was observed during setup: WSL shut down the service after the last session ended.

Start offsite and sshtty on both sides, and the existing WSL watchdog, from PowerShell. The Windows tunnel starts separately as a service:

```powershell
& F:\.dotfiles\config\wsl\start.ps1
```

The script reuses an existing watchdog if one is running. The watchdog relaunches Ubuntu within a few seconds if the distro stops, and the enabled offsite and sshtty systemd services start during the new boot. It does not install a Windows startup task. The PC must be awake. The tunnel service starts automatically with Windows.

Restart or stop a service from Ubuntu:

```bash
sudo systemctl restart offsite
sudo systemctl stop sshtty
```

For a foreground session, stop the service first, then run:

```bash
cd ~
mise exec -- offsite --host 127.0.0.1 --port 10101
mise exec -- sshtty --host 127.0.0.1 --port 10102
```

On Windows, offsite on port `20202` and sshtty on port `20203` are hidden processes, not services. offsite runs as `mise.exe`, `cmd.exe`, and a `node.exe` running `offsite.js`. To restart either one, stop its processes and run `start.ps1` again. To update the Windows offsite, run `mise upgrade npm:offsite` in PowerShell, then restart it.

## Tunnel routing

`../cloudflared-win.yml` is linked to `%USERPROFILE%\.cloudflared\config.yml` on Windows. The Windows service reads the dotfiles config directly. It must list the named hostname routes before the `*.sayad.cc` route to Caddy, followed by the catch-all:

```yaml
tunnel: Desktop
ingress:
  - hostname: wsl.sayad.cc
    service: http://127.0.0.1:10101
  - hostname: win.sayad.cc
    service: http://127.0.0.1:20202
  - hostname: sshwsl.sayad.cc
    service: http://127.0.0.1:10102
  - hostname: sshwin.sayad.cc
    service: http://127.0.0.1:20203
  - hostname: "*.sayad.cc"
    service: http://127.0.0.1:10100
  - service: http_status:404
```

Windows cloudflared reaches Ubuntu through WSL's localhost forwarding. No router port forwarding or public listener is needed.

Create the hostname's DNS route from Windows if it is missing. The Windows `cert.pem` must be logged in to the `sayad.cc` zone:

```powershell
cloudflared tunnel route dns Desktop wsl.sayad.cc
cloudflared tunnel route dns Desktop win.sayad.cc
cloudflared tunnel route dns Desktop *.sayad.cc
```

`sshwsl.sayad.cc` and `sshwin.sayad.cc` have no DNS routes of their own. The `*.sayad.cc` record already sends them to the tunnel.

### Any local port

`https://<port>.sayad.cc` opens `localhost:<port>` on Windows or Ubuntu. The tunnel sends `*.sayad.cc` to Caddy at `127.0.0.1:10100`, and `../caddy-win.caddyfile` forwards each request to the port in the first label. Caddy sends `Host: localhost:<port>`, so dev servers that only accept `localhost` work.

Caddy is installed with winget and runs as the auto-start `Caddy` Windows service under LocalSystem, restarting after failures with a 20-second delay. The service was created with:

```powershell
New-Service -Name Caddy -DisplayName Caddy -StartupType Automatic -BinaryPathName '"<caddy executable>" run --config "F:\.dotfiles\config\caddy-win.caddyfile" --adapter caddyfile'
sc.exe failure Caddy reset= 0 actions= restart/20000
sc.exe failureflag Caddy 1
```

To reload Caddy after changing `caddy-win.caddyfile`, run:

```powershell
sudo powershell.exe -NoProfile -Command "Restart-Service Caddy"
```

`start.ps1` starts offsite and sshtty in Ubuntu and on Windows, and the Windows-side WSL watchdog. It skips processes that are already running. The existing dotfiles startup launcher calls this script through `launch.jsonc`. It does not launch cloudflared. No cloudflared service runs inside WSL.

### Windows service setup and recovery

Installing the cloudflared executable does not install its Windows service. The native service was installed with:

```powershell
sudo cloudflared service install
```

Installation alone did not connect `Desktop`. The service ran as SYSTEM without the user's configuration. Its registered command line was then configured using values discovered from this machine:

- The executable path came from `Get-Command cloudflared`.
- The configuration path came from resolving `../cloudflared-win.yml` in this repository.
- The tunnel UUID came from `cloudflared tunnel list --output json --name Desktop`.
- The credentials file came from the current user's `.cloudflared` directory, using the resolved UUID. SYSTEM already had read access.
- The log directory came from `$env:ProgramData`, under `cloudflared\logs`.

The resulting service command has this form:

```text
"<cloudflared executable>" --config "<dotfiles config>" --log-directory "<log directory>" tunnel run --credentials-file "<credentials file>" <resolved tunnel UUID>
```

The service's `PathName` was updated through `Invoke-CimMethod` using the `Win32_Service.Change` method, then the service was restarted. Machine-specific paths and the resolved UUID are stored only in Windows service settings, not the dotfiles config. The credentials stay outside the repository. `start.ps1` no longer starts a separate connector.

Cloudflared's installer sets automatic startup, so the service starts after a Windows reboot without running `start.ps1` or logging in. It also configures Windows to restart the service after failures with a 20-second delay, indefinitely, including non-crash failures. If startup fails while offline, Windows keeps restarting it so it can connect when internet returns. No retry wrapper is needed. An intentional service stop does not trigger recovery.

The running service, four active Cloudflare connections, automatic startup, and recovery settings were verified after setup. A Windows reboot and an offline recovery cycle were not tested. To inspect the service and tunnel:

```powershell
Get-Service Cloudflared
sc.exe qc Cloudflared
sc.exe qfailure Cloudflared
sc.exe qfailureflag Cloudflared
cloudflared tunnel info Desktop
```

To reload the tunnel after changing `cloudflared-win.yml`, run:

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

Install and enable the services inside Ubuntu:

```bash
sudo install -m 644 /mnt/f/.dotfiles/config/wsl/offsite.service /etc/systemd/system/offsite.service
sudo install -m 644 /mnt/f/.dotfiles/config/wsl/sshtty.service /etc/systemd/system/sshtty.service
sudo systemctl daemon-reload
sudo systemctl enable --now offsite sshtty
```

Log in to Claude Code once inside Ubuntu with `claude`. offsite uses that login.

Update offsite and sshtty from Ubuntu, then restart them:

```bash
cd ~
mise upgrade npm:offsite github:nazmussayad/sshtty
sudo systemctl restart offsite sshtty
```

Use `ni`, `nr`, and `nlx` for project package workflows. `npm:offsite` is a mise package identifier; the installed command is `offsite`.

## Troubleshooting

- Windows localhost fails: check `systemctl status offsite sshtty` inside Ubuntu first, then `curl -I http://127.0.0.1:10101` and `curl -I http://127.0.0.1:10102` in both environments.
- Localhost works but the hostname returns 404: the connector may still be using the old ingress configuration.
- The hostname returns 502: cloudflared cannot reach offsite or sshtty. Check that Ubuntu is running and the server is listening on its port.
- The hostname returns 525: check the DNS route. The hostname may be reaching a different origin instead of this tunnel.
- Keep the Cloudflare Access applications for `wsl.sayad.cc` and `win.sayad.cc` separate to avoid cross-hostname login redirects.
- Restarting the service stops running Claude sessions and terminal tabs. Sessions are saved and can be resumed.
- Avoid `wsl --shutdown` during normal work; it also stops Docker's WSL environment.

## References

- [offsite usage and data location](https://github.com/NazmusSayad/offsite#readme)
- [WSL networking and localhost forwarding](https://learn.microsoft.com/en-us/windows/wsl/networking)
- [WSL systemd behavior](https://learn.microsoft.com/en-us/windows/wsl/systemd)
- [mise npm backend](https://mise.jdx.dev/dev-tools/backends/npm.html)
