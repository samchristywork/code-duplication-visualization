function makeLine(line) {
  let span = document.createElement('span');

  span.className = 'line';
  span.innerText = line.Line + '\n';
  span.style.color = `rgba(${line.Color.R}, ${line.Color.G}, ${line.Color.B}, ${line.Color.A})`;

  if (line.Tooltip) {
    span.addEventListener('mouseover', () => {
      tooltip.innerText = line.Tooltip;
      tooltip.style.display = 'block';
    });

    span.addEventListener('mouseout', () => {
      tooltip.style.display = 'none';
    });

    span.addEventListener('mousemove', (e) => {
      tooltip.style.left = `${e.clientX + 10}px`;
      tooltip.style.top = `${e.clientY + 10}px`;
    });
  }

  return span;
}
