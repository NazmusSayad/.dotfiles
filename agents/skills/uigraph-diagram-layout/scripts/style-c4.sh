# style-c4.sh NODE_ID COLOR TEXT_COLOR — C4 nodes only expose Color (border tint) and Text Color
id=$1; color=$2; text=$3
playwright-cli -s=${UIG_SESSION:-uig} run-code "$(cat <<JS
async page => {
  await page.mouse.click(30, 960);
  await page.locator('.react-flow__node[data-id="$id"]').click({ position: { x: 10, y: 10 } });
  await page.waitForTimeout(400);
  await page.getByRole('button', { name: 'Style' }).first().click();
  await page.waitForTimeout(400);
  const set = async (label, value) => { const f = page.getByText(label, { exact: true }).locator('xpath=../..').locator('input:not([type=range]):not([type=hidden])').locator('visible=true').first(); await f.click(); await page.keyboard.press('Meta+A'); await page.keyboard.type(value); await page.keyboard.press('Tab'); await page.waitForTimeout(350); };
  await set('Color', '$color');
  await set('Text Color', '$text');
}
JS
)" 2>&1 | grep -iE "error" | head -3
