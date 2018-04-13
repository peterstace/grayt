var activeUuid = '';
var lastLoadPasses = {};

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
      for (field in obj) {
        statusTxt += '<tr><td align="right">' + field + '</td><td>' + obj[field] + '</td></tr>';
      }
      statusTxt += '</table>';
      statusSpan.innerHTML = statusTxt;

      if (!(obj.uuid in lastLoadPasses)) {
        lastLoadPasses[obj.uuid] = 1; // avoid division by 0
      }
      if (obj.passes / lastLoadPasses[obj.uuid] > 1.01) {
        lastLoadPasses[obj.uuid] = obj.passes;
        // Use a cache-breaker so we get the new image each time this changes.
        let imgRender = document.getElementById('img-render');
        let url = 'http://localhost:6060/renders/' + activeUuid + '/image?' + Date.now();
        imgRender.setAttribute('src', url);
      }
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
    let data = JSON.parse(this.response);
    let inner = '';
    for (let i = 0; i < data.length; i++) {
      inner += '<option value="' + data[i].code + '">' + data[i].code + '</option>';
    }
    document.getElementById('scene-selection').innerHTML = inner;
  };
  xhr.send();
}

populateSceneSelector();

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
  xhr.send();

  xhr = new XMLHttpRequest();
  xhr.open('PUT', 'http://localhost:6060/renders/' + activeUuid + '/scene', false);
  xhr.onload = function() {
    if (this.status != 200) {
      alert(this.status);
      return;
    }
  }
  xhr.send(document.getElementById('scene-selection').value);

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
        resourceList.innerHTML += '<li class="selected">' + data[i] + "</li>";
      } else {
        resourceList.innerHTML += "<li>" + data[i] + "</li>";
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
