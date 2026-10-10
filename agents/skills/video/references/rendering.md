# Rendering

## Video

- H.264, CRF 14, `yuv420p`, 60 fps, AAC at 320k.
- Never render while footage is being screen-recorded: it drops the capture below 60 fps. Take turns, with a lock if several agents share the machine.

## Verify

- `ffprobe` the output: 1920x1080, `r_frame_rate` and `avg_frame_rate` both 60/1, the expected frame count, an audio stream.
- Check true 60 fps on a section with motion by counting frames that changed in each second:

  ```bash
  ffmpeg -v error -i out.mp4 -f framemd5 - | grep -v "^#" | tr -d ' ' \
    | awk -F, -v tb=60 '{t=int($3/tb); n[t]++; if($6!=prev) u[t]++; prev=$6} END{for(k in n) print "sec", k, "frames", n[k], "changed", u[k]}' | sort -k2 -n
  ```

- Extract and look at frames from the rendered file around every changed scene. Stills from the editor aren't enough.

## Delivery

- Keep the full-quality master. For a repo, commit a smaller copy (CRF ~24 kept 1080p60 sharp at about 14 MB per minute).
- In a README, use absolute GitHub URLs if the README is also published to npm. GitHub doesn't play repo-stored videos inline; for an inline player, the user drags the file into the README editor on github.com.
