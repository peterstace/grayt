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
          <td>px high</td>
          <td>px wide</td>
          <td>passes</td>
          <td>completed</td>
          <td colspan="3">workers</td>
          <td>image</td>
        </tr>
    `;
    for (let i = 0; i < obj.length; i++) {
      statusTxt += `
        <tr>
          <td>${obj[i].scene}</td>
          <td>${obj[i].px_high}</td>
          <td>${obj[i].px_wide}</td>
          <td>${obj[i].passes}</td>
          <td>${obj[i].completed}</td>
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
      alert(this.status);
      return;
    }
    scenes = JSON.parse(this.response);
    let inner = '';
    for (let i = 0; i < scenes.length; i++) {
      inner += `<option value="${scenes[i].code}">${scenes[i].code}</option>`;
    }
    document.getElementById('scene-selection').innerHTML = inner;
    populateResolutionCheckboxes();
  };
  xhr.send();
}

document.getElementById('scene-selection').addEventListener("change", populateResolutionCheckboxes);

populateSceneSelector();

function populateResolutionCheckboxes() {
  let resolutionsDiv = document.getElementById('resolutions');
  let selected = document.getElementById('scene-selection').value;
  let xWides = [640, 800, 1024, 1152, 1280, 1400, 1440, 1680, 1920, 2048, 2560, 2880, 3840, 4096, 5120];
  for (let i = 0; i < scenes.length; i++) {
    let scene = scenes[i];
    if (selected == scene.code) {
      let inner = '';
      for (let j = 0; j < xWides.length; j++) {
        let pxWide = xWides[j];
        if ((scene.aspect_high * pxWide) % scene.aspect_wide === 0) {
          let pxHigh = scene.aspect_high * pxWide / scene.aspect_wide;
          inner += `<option value="${pxWide}">${pxWide}x${pxHigh}</option>`;
        }
      }
      resolutionsDiv.innerHTML = inner;
      return;
    }
  }
  resolutionsDiv.innerHTML = ''; // Couldn't find the selected scene.
}

function handleAddResource() {
  let uuid = '';
  let xhr = new XMLHttpRequest();
  xhr.open('POST', 'http://localhost:6060/renders', false);
  xhr.onload = function() {
    if (this.status != 200) {
      alert(this.status);
      return;
    }
    let data = JSON.parse(this.response);
    uuid = data.uuid;
    updateStatus();
  }
  xhr.send(JSON.stringify({
    scene: document.getElementById('scene-selection').value,
    px_wide: Number(document.getElementById('resolutions').value),
  }));
}

document.getElementById("add-resource").addEventListener("click", handleAddResource)
