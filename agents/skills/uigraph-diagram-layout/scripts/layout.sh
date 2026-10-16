# layout.sh — print every node's flow position, size and center, and every edge's path endpoints
playwright-cli -s=${UIG_SESSION:-uig} run-code "$(cat <<'JS'
async page => {
  const nodes = await page.$$eval('.react-flow__node', els => els.map(e => { const m = e.style.transform.match(/translate\(([-\d.]+)px, ([-\d.]+)px\)/); return e.dataset.id + ' x=' + (+m[1]).toFixed(1) + ' y=' + (+m[2]).toFixed(1) + ' ' + e.offsetWidth + 'x' + e.offsetHeight + ' cx=' + (+m[1] + e.offsetWidth / 2).toFixed(1) + ' cy=' + (+m[2] + e.offsetHeight / 2).toFixed(1); }));
  const edges = await page.$$eval('.react-flow__edge', els => els.map(e => { const n = e.querySelector('path.react-flow__edge-path').getAttribute('d').match(/-?[\d.]+/g).map(Number); return e.dataset.id + ' start ' + n[0].toFixed(1) + ',' + n[1].toFixed(1) + ' end ' + n[n.length - 2].toFixed(1) + ',' + n[n.length - 1].toFixed(1); }));
  return nodes.concat(edges).join('\n');
}
JS
)" 2>&1 | sed -n '/### Result/,/### Ran/p' | sed 's/\\n/\n/g'
