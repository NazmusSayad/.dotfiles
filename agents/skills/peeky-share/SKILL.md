---
name: peeky-share
description: Use when the user must see something from this machine in a browser, for example a file, an image, an HTML page, a report, or a local server.
---

The user opens shared content at `https://<label>.sayad.cc`. The Cloudflare tunnel `peeky` sends all addresses under `sayad.cc` to one local port. Cloudflare Access protects all of these addresses. The user must log in with their email before they see the content. Do not add a different login.

## Use one port only

The tunnel sends all addresses to the port that you give with `--url`. If two tunnel processes run at the same time, each request goes to one of them at random. Thus the user can see the wrong content.

Before you start the tunnel, stop all tunnel processes:

```
pkill -f "cloudflared tunnel run.*peeky"
```

To share many items at the same time, put them all in one directory. Then serve that directory.

## Share content

1. Get the port.
   - If a local server shows the content, use the port of that server.
   - If the content is in files, copy the files into a directory. Then serve the directory on a free port:

     ```
     PORT=$(python3 -c 'import socket;s=socket.socket();s.bind(("127.0.0.1",0));print(s.getsockname()[1])')
     nohup python3 -m http.server "$PORT" --bind 127.0.0.1 --directory DIR > /tmp/peeky-http.log 2>&1 &
     ```

2. Start the tunnel in the background:

   ```
   nohup cloudflared tunnel run --url "http://localhost:$PORT" peeky > /tmp/peeky.log 2>&1 &
   ```

   Wait until `/tmp/peeky.log` shows `Registered tunnel connection`.

3. Select a short label that tells what the content is, for example `report` or `ui-demo`. Do not use a label that has its own DNS record. That record gets the requests, and the tunnel does not.

4. Test the address:

   ```
   curl -sI "https://<label>.sayad.cc/<path>"
   ```

   If the result sends you to `cloudflareaccess.com`, the address is correct. Any other result shows a problem with the tunnel or the label.

5. Give the user the full address. Include the path to the file.

## Stop sharing

Stop the tunnel and the file server when the user does not need the content. Also stop them before you share content on a different port.
