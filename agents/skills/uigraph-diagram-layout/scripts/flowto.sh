# flowto.sh NODE_ID X Y — drag a node until its top-left is at flow coords (X, Y); prints remaining error per attempt
id=$1; fx=$2; fy=$3
playwright-cli -s=${UIG_SESSION:-uig} run-code "$(cat <<JS
async page => {
  const n = page.locator('.react-flow__node[data-id="$id"]');
  const vp = async () => { const m = (await page.\$eval('.react-flow__viewport', el => el.style.transform)).match(/scale\(([\d.]+)\)/); return +m[1]; };
  const pos = async () => { const m = (await n.evaluate(el => el.style.transform)).match(/translate\(([-\d.]+)px, ([-\d.]+)px\)/); return { x: +m[1], y: +m[2] }; };
  const log = [];
  for (let i = 0; i < 4; i++) {
    const s = await vp(), p = await pos();
    log.push(($fx - p.x).toFixed(1) + ',' + ($fy - p.y).toFixed(1));
    if (Math.abs($fx - p.x) < 0.8 && Math.abs($fy - p.y) < 0.8) break;
    const b = await n.boundingBox();
    const gx = b.x + 8, gy = b.y + b.height / 2;
    await page.mouse.click(30, 960);
    await page.mouse.move(gx, gy);
    await page.mouse.down();
    await page.mouse.move(gx + 20, gy + 20, { steps: 5 });
    await page.mouse.move(gx + ($fx - p.x) * s + 4, gy + ($fy - p.y) * s + 4, { steps: 8 });
    await page.mouse.up();
    await page.waitForTimeout(250);
  }
  return log.join(' -> ');
}
JS
)" 2>&1 | grep -A1 "### Result" | tail -1
