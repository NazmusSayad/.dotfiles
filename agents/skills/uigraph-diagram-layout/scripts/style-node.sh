# style-node.sh NODE_ID WIDTH HEIGHT FILL STROKE TEXT Solid|Dashed — style a regular (non-C4) node through the Style panel
id=$1; w=$2; h=$3; fill=$4; stroke=$5; text=$6; border=$7
playwright-cli -s=${UIG_SESSION:-uig} run-code "$(cat <<JS
async page => {
  await page.mouse.click(30, 960);
  await page.locator('.react-flow__node[data-id="$id"]').click({ position: { x: 10, y: 10 } });
  await page.waitForTimeout(400);
  const style = page.getByText('Style', { exact: true });
  if (await style.isVisible()) await style.click();
  await page.waitForTimeout(300);
  const set = async (label, value) => { const f = page.getByText(label, { exact: true }).locator('xpath=../..').locator('input:not([type=range]):not([type=hidden])').locator('visible=true').first(); await f.click(); await page.keyboard.press('Meta+A'); await page.keyboard.type(String(value), { delay: 30 }); await page.keyboard.press('Tab'); await page.waitForTimeout(350); };
  await set('Text Font Size', 15);
  await set('Text Color', '$text');
  await set('Background Color', '$fill');
  await set('Width (px)', $w);
  await set('Height (px)', $h);
  await page.getByRole('radio', { name: '$border' }).click(); await page.waitForTimeout(350);
  await set('Border Color', '$stroke');
  await set('Border Width', 1.5);
  await set('Border Radius', 10);
  await page.waitForTimeout(400);
}
JS
)" 2>&1 | grep -E "^### Error" -A3 | head -6
