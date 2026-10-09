# Ubuntu WSL and ttysh

Run Linux shells at https://sh.sayad.dev through the existing `opencode` Cloudflare tunnel. Authentication is managed separately in Cloudflare.

## Setup status

Checked on 2026-10-09:

- Ubuntu and the Linux `sayad` account are configured.
- Node.js `24.21.0`, pnpm `12.10.1`, `@antfu/ni` `30.6.0`, the npm package `ttysh` `0.0.2`, and Ubuntu's Starship `1.22.1` package are installed.
- `ttysh.service` is installed, enabled, and running as `sayad`.
- Ubuntu and Windows both return HTTP 200 at `http://127.0.0.1:47474`.
- The DNS route for `sh.sayad.dev` points to the existing tunnel, and the updated ingress config passes validation.
- The Windows tunnel launcher was restarted, and `https://sh.sayad.dev/` returns HTTP 200 with the ttysh page.
- A hidden WSL keep-alive process is running so Ubuntu does not stop when the last terminal closes.

OpenCode is not installed in Ubuntu. The name `opencode` refers to the existing Windows tunnel. Its Windows launcher is `opencode-server`.

## Device setup

- Distro: `Ubuntu-26.04`, running on WSL 2.
- Linux account: `sayad`, with Bash as its shell.
- Windows dotfiles: `F:\.dotfiles`, available inside Ubuntu at `/mnt/f/.dotfiles`.
- Linux tools: mise, Node.js LTS, pnpm, `@antfu/ni`, `ttysh`, and the Ubuntu Starship package.
- Mise config: `~/.config/mise/config.toml`, linked to the main `../mise-config.toml` used by this dotfiles repository.
- Bash loads mise and Starship from `bashrc.sh`. Ubuntu installs Starship with `apt`, while Starship uses the shared `../shell/starship.toml` config.
- ttysh data: `/home/sayad/.ttysh`. Keep it in Linux, not on a Windows mount.
- ttysh's `shell.command` is `/bin/bash`, so new browser tabs explicitly open Bash.
- ttysh listener: `127.0.0.1:47474` inside Ubuntu.
- Tunnel: `opencode`, UUID `3e3fbb6c-25de-4bad-aa96-fa04d4eeba18`, running on Windows.

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

## Run ttysh

The `ttysh.service` file runs ttysh as `sayad`, not root. From Ubuntu:

```bash
sudo systemctl start ttysh
sudo systemctl status ttysh --no-pager
journalctl -u ttysh -n 50 --no-pager
```

Local address on Windows: http://127.0.0.1:47474. Remote address: https://sh.sayad.dev.

The service starts when Ubuntu boots. It does not boot Ubuntu when Windows starts, and systemd services do not keep WSL alive on their own. This was observed during setup: WSL shut down the service after the last session ended.

Start ttysh and a hidden WSL keep-alive process from PowerShell:

```powershell
& F:\.dotfiles\config\wsl\start.ps1
```

The script reuses an existing keep-alive process if one is running. It does not install a Windows startup task. Run it after signing in or after stopping WSL. The PC must be awake and the Windows tunnel must also be running.

Restart or stop the service from Ubuntu:

```bash
sudo systemctl restart ttysh
sudo systemctl stop ttysh
```

For a foreground session, stop the service first, then run:

```bash
cd ~
mise exec -- ttysh --host 127.0.0.1 --port 47474
```

## Tunnel routing

`../cloudflared.yml` is linked to `%USERPROFILE%\.cloudflared\config.yml` on Windows. It must contain both hostname routes, followed by the catch-all:

```yaml
tunnel: opencode
ingress:
  - hostname: oc.sayad.dev
    service: http://127.0.0.1:4747
  - hostname: sh.sayad.dev
    service: http://127.0.0.1:47474
  - service: http_status:404
```

Windows cloudflared reaches Ubuntu through WSL's localhost forwarding. No router port forwarding or public listener is needed.

Create the hostname's DNS route from Windows if it is missing:

```powershell
cloudflared tunnel route dns opencode sh.sayad.dev
```

Restart the existing tunnel connector after editing its ingress configuration. The current connector is started by the Windows `opencode-server` launcher, so restarting that launcher also restarts the Windows OpenCode server. Do not launch a second connector with different ingress rules for the same tunnel.

The launcher was restarted as a hidden Windows process during setup. Find its PID from PowerShell:

```powershell
Get-CimInstance Win32_Process | Where-Object Name -eq opencode-server.exe | Select-Object ProcessId
```

To reload it, stop that specific process tree with `taskkill /PID <PID> /T /F`, replacing `<PID>` with the displayed ID, then run `opencode-server` in a Windows terminal. This briefly interrupts `oc.sayad.dev` and `sh.sayad.dev`.

The launcher's password is saved in the Windows user environment. An old PowerShell session may not have inherited it. Load it without displaying or copying its value:

```powershell
$env:OPENCODE_SERVER_PASSWORD = [Environment]::GetEnvironmentVariable('OPENCODE_SERVER_PASSWORD', 'User')
opencode-server
```

The current hidden launcher's logs are in `%LOCALAPPDATA%\Temp\opencode\wsl-tunnel.stdout.log` and `wsl-tunnel.stderr.log`. No OpenCode installation in Ubuntu is needed.

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
ln -s /mnt/f/.dotfiles/config/mise-config.toml ~/.config/mise/config.toml
ln -s /mnt/f/.dotfiles/config/shell/starship.toml ~/.config/starship.toml
~/.local/bin/mise trust ~/.config/mise/config.toml
~/.local/bin/mise install
```

Install Starship separately through Ubuntu. This leaves the main mise config unchanged:

```bash
sudo apt update
sudo apt install starship
```

Add this line once to Ubuntu's `~/.bashrc`, then reopen the shell:

```bash
source /mnt/f/.dotfiles/config/wsl/bashrc.sh
```

Do not replace Ubuntu's `.bashrc` with `../shell/bashrc.sh`. That shared file currently expects Windows/macOS tools and does not activate mise on Linux.

Install and enable the service inside Ubuntu:

```bash
sudo install -m 644 /mnt/f/.dotfiles/config/wsl/ttysh.service /etc/systemd/system/ttysh.service
sudo systemctl daemon-reload
sudo systemctl enable --now ttysh
```

Update ttysh from Ubuntu, then restart it:

```bash
cd ~
mise upgrade npm:ttysh
sudo systemctl restart ttysh
```

Use `ni`, `nr`, and `nlx` for project package workflows. `npm:ttysh` is a mise package identifier; the installed command is `ttysh`.

## Troubleshooting

- Windows localhost fails: check `systemctl status ttysh` inside Ubuntu first, then `curl -I http://127.0.0.1:47474` in both environments.
- Localhost works but the hostname returns 404: the connector may still be using the old ingress configuration.
- The hostname returns 502: cloudflared cannot reach ttysh. Check that Ubuntu is running and ttysh is listening.
- The hostname returns 525: check the DNS route. The hostname may be reaching a different origin instead of this tunnel.
- Shells survive ttysh server crashes according to the package documentation. They do not survive Windows shutdown or `wsl --shutdown`.
- The installed npm package is `ttysh` `0.0.2`, but its binary reports `ttysh 0.0.0` with `--version`. Use `mise ls npm:ttysh` to check the installed package version.
- Avoid `wsl --shutdown` during normal work; it also stops Docker's WSL environment.

## References

- [ttysh usage and data location](https://github.com/NazmusSayad/ttysh#readme)
- [WSL networking and localhost forwarding](https://learn.microsoft.com/en-us/windows/wsl/networking)
- [WSL systemd behavior](https://learn.microsoft.com/en-us/windows/wsl/systemd)
- [mise npm backend](https://mise.jdx.dev/dev-tools/backends/npm.html)
