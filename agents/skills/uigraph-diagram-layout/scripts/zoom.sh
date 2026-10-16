# zoom.sh SCALE — ctrl+wheel to SCALE around canvas center, then pan content to the center
target=$1
playwright-cli -s=${UIG_SESSION:-uig} run-code "$(cat <<JS
async page => {
  const scale = async () => parseFloat((await page.\$eval('.react-flow__viewport', el => el.style.transform)).match(/scale\(([\d.]+)\)/)[1]);
  await page.mouse.move(800, 560);
  await page.keyboard.down('Control');
  for (let i = 0; i < 80 && (await scale()) > $target + 0.03; i++) { await page.mouse.wheel(0, 20); await page.waitForTimeout(50); }
  for (let i = 0; i < 80 && (await scale()) < $target - 0.03; i++) { await page.mouse.wheel(0, -20); await page.waitForTimeout(50); }
  await page.keyboard.up('Control');
  await page.waitForTimeout(200);
  const boxes = await page.\$\$eval('.react-flow__node', els => els.map(e => e.getBoundingClientRect()).map(r => [r.left, r.top, r.right, r.bottom]));
  const l = Math.min(...boxes.map(b => b[0])), t = Math.min(...boxes.map(b => b[1])), r = Math.max(...boxes.map(b => b[2])), b = Math.max(...boxes.map(b => b[3]));
  await page.mouse.move(800, 560);
  await page.mouse.wheel(-(820 - (l + r) / 2), -(575 - (t + b) / 2));
  await page.waitForTimeout(300);
  return await page.\$eval('.react-flow__viewport', el => el.style.transform);
}
JS
)" 2>&1 | grep -A1 "### Result" | tail -1
