# resize-boundary.sh BOUNDARY_ID top|bottom|left|right FLOW_DELTA — drag a C4 boundary edge (selects it via its "[SYSTEM]" label)
id=$1; ctl=$2; delta=$3
playwright-cli -s=${UIG_SESSION:-uig} run-code "$(cat <<JS
async page => {
  const n = page.locator('.react-flow__node[data-id="$id"]');
  await n.getByText('[SYSTEM]').first().click();
  await page.waitForTimeout(400);
  const s = parseFloat((await page.\$eval('.react-flow__viewport', el => el.style.transform)).match(/scale\(([\d.]+)\)/)[1]);
  const b = await n.locator('.react-flow__resize-control.line.$ctl').boundingBox();
  const vertical = '$ctl' === 'bottom' || '$ctl' === 'top';
  const gx = vertical ? b.x + b.width * 0.75 : b.x + b.width / 2, gy = vertical ? b.y + b.height / 2 : b.y + b.height * 0.25;
  const d = $delta * s;
  await page.mouse.move(gx, gy);
  await page.mouse.down();
  await page.mouse.move(gx + (vertical ? 0 : d / 2), gy + (vertical ? d / 2 : 0), { steps: 6 });
  await page.mouse.move(gx + (vertical ? 0 : d), gy + (vertical ? d : 0), { steps: 6 });
  await page.mouse.up();
  await page.waitForTimeout(400);
  await page.mouse.click(30, 960);
  return await n.evaluate(e => e.style.transform + ' ' + e.offsetWidth + 'x' + e.offsetHeight);
}
JS
)" 2>&1 | grep -A1 "### Result" | tail -1
