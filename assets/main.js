var scenes = [];

function updateStatus() {
  let xhr = new XMLHttpRequest();
  xhr.open('GET', 'http://localhost:6060/renders');
  xhr.onload = function() {
    let obj = JSON.parse(this.response);
    let statusTxt = `
      <table>
        <tr>
          <td>scene name</td>
          <td>dimensions</td>
          <td>passes</td>
          <td>completed</td>
          <td>trace rate</td>
          <td colspan="3">workers</td>
          <td>image</td>
        </tr>
    `;
    for (let i = 0; i < obj.length; i++) {
      statusTxt += `
        <tr>
          <td>${obj[i].scene}</td>
          <td>${obj[i].px_wide}x${obj[i].px_high}</td>
          <td>${obj[i].passes}</td>
          <td>${obj[i].completed}</td>
          <td>${obj[i].trace_rate}</td>
          <td><button class="worker" uuid=${obj[i].uuid} workers=${obj[i].requested_workers-1}>-</button></td>
          <td>${obj[i].requested_workers} (${obj[i].actual_workers})</td>
          <td><button class="worker" uuid=${obj[i].uuid} workers=${obj[i].requested_workers+1}>+</button></td>
          <td>
            <a
              href="http://localhost:6060/renders/${obj[i].uuid}/image"
              target="_blank"
            >image</a>
          </td>
        </tr>`
    }
    statusTxt += '</table>';
    document.getElementById('status').innerHTML = statusTxt;

    let buttons = document.getElementsByClassName('worker');
    for (let i = 0; i < buttons.length; i++) {
      let btn = buttons[i];
      btn.addEventListener("click", function() {
        let xhr = new XMLHttpRequest();
        let uuid = btn.getAttribute('uuid');
        xhr.open('PUT', `http://localhost:6060/renders/${uuid}/workers`);
        xhr.onload = updateStatus;
        let workers = btn.getAttribute('workers');
        xhr.send(workers);
      });
    }
  };
  xhr.send();
}

window.setInterval(updateStatus, 250);

function populateSceneSelector() {
  let xhr = new XMLHttpRequest();
  xhr.open('GET', 'http://localhost:6060/scenes');
  xhr.onload = function () {
    if (this.status != 200) {
      alert(`${this.status}: ${this.response}`);
      return;
    }
    scenes = JSON.parse(this.response);
    let inner = '';
    for (let i = 0; i < scenes.length; i++) {
      inner += `<option value="${scenes[i].code}">${scenes[i].code}</option>`;
    }
    document.getElementById('scene-selection').innerHTML = inner;
    populateResolutionSelection();
  };
  xhr.send();
}

populateSceneSelector();

function populateAspectSelection() {
  let aspectDiv = document.getElementById('aspects');
  let ratios = [[1,1], [4,3], [16,9], [16,10], [2,1]];
  let inner = '';
  for (let i = 0; i < ratios.length; i++) {
    let rat = `${ratios[i][0]}:${ratios[i][1]}`
    inner += `<option value="${rat}">${rat}</option>`;
  }
  aspectDiv.innerHTML = inner;
}

populateAspectSelection();

function populateResolutionSelection() {
  const aspectSelect = document.getElementById('aspects');
  const aspects = aspectSelect.value.split(':');
  const aspectWide = aspects[0];
  const aspectHigh = aspects[1];
  const xWides = [640, 800, 1024, 1152, 1280, 1400, 1440, 1680, 1920, 2048, 2560, 2880, 3840, 4096, 5120];
  let inner = '';
  for (let i = 0; i < xWides.length; i++) {
    const pxWide = xWides[i];
    if ((aspectHigh * pxWide) % aspectWide === 0) {
      const pxHigh = aspectHigh * pxWide / aspectWide;
      const dim = `${pxWide}x${pxHigh}`
      inner += `<option value="${dim}">${dim}</option>`;
    }
  }
  document.getElementById('resolutions').innerHTML = inner;
}

document.getElementById('aspects').addEventListener("change", populateResolutionSelection);

function handleAddResource() {
  let uuid = '';
  let xhr = new XMLHttpRequest();
  xhr.open('POST', 'http://localhost:6060/renders', false);
  xhr.onload = function() {
    if (this.status != 200) {
      alert(`${this.status}: ${this.response}`);
      return;
    }
    let data = JSON.parse(this.response);
    uuid = data.uuid;
    updateStatus();
  }
  const dim = document.getElementById('resolutions').value.split("x");
  xhr.send(JSON.stringify({
    scene: document.getElementById('scene-selection').value,
    px_wide: Number(dim[0]),
    px_high: Number(dim[1]),
  }));
}

document.getElementById("add-resource").addEventListener("click", handleAddResource)
