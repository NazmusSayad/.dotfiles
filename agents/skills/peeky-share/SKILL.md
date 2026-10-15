---
name: peeky-share
description: Share a local server. Use when using Cloudflare Tunnel.
---

## Cloudflare Tunnel

The tunnel `peeky` serves `*.sayad.cc` behind Cloudflare Access.

## Share

Run in the background of your session, so it stops when you stop:

```
pkill -f "cloudflared tunnel run.*peeky"
cloudflared tunnel run --url "http://localhost:PORT" peeky
```

## Stop

Stop the tunnel when the user is done.
