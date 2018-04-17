var activeUuid = '';
var scenes = [];

function updateStatus() {
  let statusSpan = document.getElementById('status');
  if (activeUuid == '') {
    statusSpan.innerHTML = '';
  } else {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', 'http://localhost:6060/renders/' + activeUuid);
    xhr.onload = function() {
      let statusTxt = '<table>';
      let obj = JSON.parse(this.response);
      obj.image = `<a href="http://localhost:6060/renders/${activeUuid}/image" target="_blank">link</a>`
      for (field in obj) {
        statusTxt += `<tr><td align="right">${field}</td><td>${obj[field]}</td></tr>`
      }
      statusTxt += '</table>';
      statusSpan.innerHTML = statusTxt;
    };
    xhr.send();
  }
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
  let xWides = [640, 800, 1024];
  for (let i = 0; i < scenes.length; i++) {
    let scene = scenes[i];
    if (selected == scene.code) {
      let inner = '';
      for (let j = 0; j < xWides.length; j++) {
        let pxWide = xWides[j];
        let pxHigh = scene.aspect_high * pxWide / scene.aspect_wide;
        inner += `<option value="${pxWide}">${pxWide}x${pxHigh}</option>`;
      }
      resolutionsDiv.innerHTML = inner;
      return;
    }
  }
  resolutionsDiv.innerHTML = ''; // Couldn't find the selected scene.
}

function handleAddResource() {
  let xhr = new XMLHttpRequest();
  xhr.open('POST', 'http://localhost:6060/renders', false);
  xhr.onload = function() {
    if (this.status != 200) {
      alert(this.status);
      return;
    }
    let data = JSON.parse(this.response);
    activeUuid = data.uuid;
    updateResourceList();
    updateStatus();
  }
  xhr.send(JSON.stringify({
    scene: document.getElementById('scene-selection').value,
    px_wide: Number(document.getElementById('resolutions').value),
  }));

  xhr = new XMLHttpRequest();
  xhr.open('PUT', 'http://localhost:6060/renders/' + activeUuid + '/running', false);
  xhr.onload = function() {
    if (this.status != 200) {
      alert(this.status);
      return;
    }
  }
  xhr.send('true');
}

document.getElementById("add-resource").addEventListener("click", handleAddResource)

function updateResourceList() {
  let xhr = new XMLHttpRequest();
  xhr.open('GET', 'http://localhost:6060/renders');
  xhr.onload = function () {
    let data = JSON.parse(this.response);
    document.getElementById('resources').innerHTML = "Resources (" + data.length + "):";
    let resourceList = document.getElementById('resource-list');
    resourceList.innerHTML = '';
    for (let i = 0; i < data.length; i++) {
      if (data[i] == activeUuid) {
        resourceList.innerHTML += `<li class="selected">${data[i]}</li>`;
      } else {
        resourceList.innerHTML += `<li>${data[i]}</li>`;
      }
    }
    let listItems = resourceList.childNodes;
    for (let i = 0; i < listItems.length; i++) {
      listItems[i].addEventListener("click", function() {
        activeUuid = this.innerHTML;
        updateResourceList();
      })
    }
  }
  xhr.send();
}

updateResourceList();
