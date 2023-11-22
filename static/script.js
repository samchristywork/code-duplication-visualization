let tooltip = document.getElementById('tooltip');

function makeLine(line) {
  let span = document.createElement('span');

  span.className = 'line';
  span.innerText = line.Line;
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

fetch(`/files`)
  .then(response => response.json())
  .then(data => {
    let files = data;

    let target = document.getElementById('target');
    for (let file of files) {
      let panel = document.createElement('div');
      panel.className = 'panel';
      target.appendChild(panel);

      let title = document.createElement('div');
      title.className = 'title';
      title.innerText = file;
      panel.appendChild(title);

      fetch(`/file?file=${file}`)
        .then(response => response.json())
        .then(data => {
          for (let line of data) {
            let span = makeLine(line);
            panel.appendChild(span);
          }
        });
    }
  });
