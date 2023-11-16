fetch(`/files`)
  .then(response => response.json())
  .then(data => {
    let files = data;

    let target = document.getElementById('target');
    for (let file of files) {
      let panel = document.createElement('div');
      panel.className = 'panel';
      target.appendChild(panel);
      fetch(`/file?file=${file}`)
        .then(response => response.json())
        .then(data => {
          console.log(data);
          for (let line of data) {
            let span = document.createElement('span');
            span.className = 'line';
            span.innerText = line.Line;
            span.style.color = `rgba(${line.Color.R}, ${line.Color.G}, ${line.Color.B}, ${line.Color.A})`;
            panel.appendChild(span);
          }
        });
    }
  });
