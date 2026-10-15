---
name: peeky-share
description: Share a local server. Use when using Cloudflare Tunnel.
---

Every `*.sayad.cc` address goes through the tunnel to one local port. Cloudflare Access already asks the user to log in with their email, so do not add your own login.

## Share

Stop any old tunnel, then start the tunnel on the server's port as a background command of your session, not with `nohup` or `&`, so it ends when you stop. Two tunnel processes split requests at random, so the user can see the wrong content. Wait until the tunnel output shows `Registered tunnel connection`.

```
pkill -f "cloudflared tunnel run.*peeky"
cloudflared tunnel run --url "http://localhost:PORT" peeky
```

Pick a short label that names what you share, such as `report` or `ui-demo`. Avoid labels that have their own DNS record, because that record gets the requests instead of the tunnel.

Give the user the full address, including the path.

## Stop

Stop the tunnel when the user is done.
