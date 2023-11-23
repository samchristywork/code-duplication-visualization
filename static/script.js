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

      let title = document.createElement('a');
      title.className = 'title';
      title.innerText = file;
      title.href = `/display-file?filename=${file}`;
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
