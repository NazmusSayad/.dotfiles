# Rendering

## Video

H.264, CRF 14, `yuv420p`, AAC at 320k, at the target resolution and frame rate.

## Verify

- Check the output's resolution, frame rate, frame count and audio stream with `ffprobe`.
- Confirm the frame rate is real on a section with motion by counting frames that changed in each second:

  ```bash
  ffmpeg -v error -i out.mp4 -f framemd5 - | grep -v "^#" | tr -d ' ' \
    | awk -F, -v tb=60 '{t=int($3/tb); n[t]++; if($6!=prev) u[t]++; prev=$6} END{for(k in n) print "sec", k, "frames", n[k], "changed", u[k]}' | sort -k2 -n
  ```

- Look at frames from the rendered file around every changed scene. Stills from the editor aren't enough.

## Delivery

- Keep the full-quality master. Make any other copy from it at the same resolution and frame rate, and check that text stays sharp.
- Check how the destination plays video before delivering; some places won't play a linked file inline.
