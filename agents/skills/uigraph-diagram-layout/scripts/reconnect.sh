# reconnect.sh EDGE_ID_PREFIX source|target NODE_ID HANDLE_ID — drag an edge endpoint onto a handle (label is kept)
# handle ids: source-top|bottom|left|right, target-top|bottom|left|right
edge=$1; end=$2; node=$3; handle=$4
playwright-cli -s=${UIG_SESSION:-uig} run-code "$(cat <<JS
async page => {
  const e = page.locator('.react-flow__edge[data-id^="$edge"]');
  await e.locator('.react-flow__edge-interaction').dispatchEvent('click');
  await page.waitForTimeout(300);
  const a = await e.locator('.react-flow__edgeupdater-$end').boundingBox();
  const h = await page.locator('.react-flow__node[data-id="$node"] .react-flow__handle[data-handleid="$handle"]').boundingBox();
  const ax = a.x + a.width / 2, ay = a.y + a.height / 2, hx = h.x + h.width / 2, hy = h.y + h.height / 2;
  await page.mouse.move(ax, ay);
  await page.mouse.down();
  await page.mouse.move((ax + hx) / 2, (ay + hy) / 2, { steps: 10 });
  await page.mouse.move(hx, hy, { steps: 10 });
  await page.waitForTimeout(200);
  await page.mouse.up();
  await page.waitForTimeout(600);
}
JS
)" 2>&1 | grep -E "^### Error" -A3 | head -5
