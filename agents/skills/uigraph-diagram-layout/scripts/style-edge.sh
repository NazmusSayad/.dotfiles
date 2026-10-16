# style-edge.sh EDGE_ID_PREFIX COLOR [LABEL] — e.g. "edge-server-db-" or "c4-rel-server-db-" (see layout.sh for ids)
edge=$1; color=$2
playwright-cli -s=${UIG_SESSION:-uig} run-code "$(cat <<JS
async page => {
  const e = page.locator('.react-flow__edge[data-id^="$edge"]');
  await e.locator('.react-flow__edge-interaction').dispatchEvent('click');
  await page.waitForTimeout(400);
  const field = (label, up = '..') => page.locator('label').getByText(label, { exact: true }).locator('xpath=' + up).locator('input:not([type=range]):not([type=hidden])').locator('visible=true').first();
  const set = async (label, value, up) => { const f = field(label, up); await f.scrollIntoViewIfNeeded(); await f.click(); await page.keyboard.press('Meta+A'); await page.keyboard.type(String(value), { delay: 30 }); await page.keyboard.press('Tab'); await page.waitForTimeout(350); };
  if ('$3' !== '') await set('Edge Label', '$3');
  await set('Stroke Color', '$color');
  await set('End Marker Color', '$color');
  await set('Stroke Width', 1.5, '../..');
  await page.waitForTimeout(300);
}
JS
)" 2>&1 | grep -E "^### Error" -A3 | head -5
