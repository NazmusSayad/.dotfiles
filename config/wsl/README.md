# Ubuntu WSL and ttysh

Run Linux shells at https://sh.sayad.dev through the existing `Desktop` Cloudflare tunnel. Authentication is managed separately in Cloudflare.

## Setup status

Checked on 2026-10-09:

- Ubuntu and the Linux `sayad` account are configured.
- Node.js `24.21.0`, pnpm `12.10.1`, `@antfu/ni` `30.6.0`, the npm package `ttysh` `0.0.2`, and Ubuntu's Starship `1.22.1` package are installed.
- `ttysh.service` is installed, enabled, and running as `sayad`.
- Ubuntu and Windows both return HTTP 200 at `http://127.0.0.1:47474`.
- The DNS route for `sh.sayad.dev` points to the existing tunnel, and the updated ingress config passes validation.
- The Windows cloudflared connector is running. `https://sh.sayad.dev/` returns HTTP 200 through Cloudflare Access, which then forwards authenticated requests to ttysh.
- A hidden WSL keep-alive process is running so Ubuntu does not stop when the last terminal closes.

OpenCode is not part of this setup. The `Desktop` tunnel runs as an independent cloudflared process on Windows.

## Device setup

- Distro: `Ubuntu-26.04`, running on WSL 2.
- Linux account: `sayad`, with Bash as its login shell.
- Windows dotfiles: `F:\.dotfiles`, available inside Ubuntu at `/mnt/f/.dotfiles`.
- Linux tools: mise, Node.js LTS, pnpm, `@antfu/ni`, `ttysh`, and the Ubuntu Starship package.
- Mise config: `~/.config/mise/config.toml`, linked to the main `../mise-config.toml` used by this dotfiles repository.
- `config-init` detects WSL and creates the portable links from `../symlink.yml`, including OpenCode and CLI configuration.
- Running `config-init` on Windows also runs it inside `Ubuntu-26.04` as `sayad`, keeping both environments updated from one command.
- `~/.bashrc` links to the shared `../shell/bashrc.sh`. On Linux it activates mise, loads environment files, and initializes direnv, Shaka, zoxide, and Starship.
- Shaka uses the shared `../shell/alias.json` config. This provides `oc` for `opencode --standalone`.
- ttysh data: `/home/sayad/.ttysh`. Keep it in Linux, not on a Windows mount.
- ttysh's `shell.command` is `/usr/bin/fish`, so browser tabs open Fish without changing the WSL login shell.
- ttysh listener: `127.0.0.1:47474` inside Ubuntu.
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
tunnel: Desktop
ingress:
  - hostname: sh.sayad.dev
    service: http://127.0.0.1:47474
  - service: http_status:404
```

Windows cloudflared reaches Ubuntu through WSL's localhost forwarding. No router port forwarding or public listener is needed.

Create the hostname's DNS route from Windows if it is missing:

```powershell
cloudflared tunnel route dns Desktop sh.sayad.dev
```

`start.ps1` starts ttysh, the WSL keep-alive process, and `cloudflared tunnel run`. It skips processes that are already running. Windows startup runs this script through `launch.jsonc`.

To reload the tunnel after changing `cloudflared.yml`, stop its cloudflared process and run:

```powershell
& F:\.dotfiles\config\wsl\start.ps1
```

The connector logs are in `%LOCALAPPDATA%\Temp\desktop-tunnel\stdout.log` and `stderr.log`.

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
sudo apt install build-essential direnv fish starship zoxide
```

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
- If signing into `sh.sayad.dev` visits another hostname, edit the Cloudflare Access application so its only public hostname is `sh.sayad.dev`. Cloudflare preemptively visits every domain in small multi-domain applications to issue authorization cookies.
- Shells survive ttysh server crashes according to the package documentation. They do not survive Windows shutdown or `wsl --shutdown`.
- The installed npm package is `ttysh` `0.0.2`, but its binary reports `ttysh 0.0.0` with `--version`. Use `mise ls npm:ttysh` to check the installed package version.
- Avoid `wsl --shutdown` during normal work; it also stops Docker's WSL environment.

## References

- [ttysh usage and data location](https://github.com/NazmusSayad/ttysh#readme)
- [WSL networking and localhost forwarding](https://learn.microsoft.com/en-us/windows/wsl/networking)
- [WSL systemd behavior](https://learn.microsoft.com/en-us/windows/wsl/systemd)
- [mise npm backend](https://mise.jdx.dev/dev-tools/backends/npm.html)
