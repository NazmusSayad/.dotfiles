# Music

## Choosing

- Instrumental only, no vocals and no voice-over. Check for vocals before offering a track.
- Choose music for the feeling it should give the viewer, decided from the audience and the story: calm confidence, curiosity, a lift at the big claim. Write that feeling down before searching, and check every candidate against it.
- Avoid the generic stock tracks heard in every product video: the same upbeat ukulele, plinky corporate pop or bland background loop. If a track could sit under any ad, keep looking. It should have a character that fits this product.
- Aim for relaxed but modern: a smooth, steady groove that feels confident, not sleepy and not aggressive. Chillout-style house works well; hard-hitting, energetic tracks feel like an ad shouting at the viewer.
- A steady tempo that divides evenly into the frame rate makes cutting easy. At 120 BPM, a beat is exactly 30 frames at 60 fps.
- Look for this shape: the groove starts right away (a quiet intro wastes the opening), layers build gently to the fullest part, there's a calm stretch for detail scenes, and the full groove returns for the ending.
- You can't listen, so compare candidates by measurement: lower brightness, less high-frequency sizzle from hi-hats and cymbals, softer drum hits, and steady loudness. Smoother tracks score lower on all four.
- Match the motion to the music. Calmer music needs softer transitions and fewer hard whips.
- Shortlist two or three tracks and have the user listen before committing. Taste is their call.
- Use a track with a clear licence for the intended use, from a library such as Mixkit. Record the title, author, source and licence in `CREDITS.md`, including restrictions such as not registering the track with YouTube Content ID.

## Beat grid

- Analyse the track once into a JSON file: BPM, beat times, downbeats and sections.
- Expose cues by beat number (`cue.groups = beat(28)`) and place every shot on a cue.
- Scene changes go on beats; the biggest moments go on downbeats or section changes.

## Fitting length

- When the video needs more time, extend the track by repeating a section on bar boundaries, with short fades at the joins so they don't click. Fade out at the end.
- Normalise to about -14 LUFS integrated with a -1.5 dBTP true peak.
- If the extended track changes cue times, re-check every shot against its cue.
