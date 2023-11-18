let tooltip = document.getElementById('tooltip');

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
          console.log(data);
          for (let line of data) {
            let span = document.createElement('span');
            span.className = 'line';
            span.innerText = line.Line;
            span.style.color = `rgba(${line.Color.R}, ${line.Color.G}, ${line.Color.B}, ${line.Color.A})`;

            let t = line.Tooltip;
            if (t) {
              span.addEventListener('mouseover', () => {
                tooltip.innerText = t;
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

            panel.appendChild(span);
          }
        });
    }
  });
