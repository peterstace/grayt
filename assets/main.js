var activeUuid = '';

function populateSceneSelector() {
  let xhr = new XMLHttpRequest();
  xhr.open('GET', 'http://localhost:6060/scenes');
  xhr.onload = function () {
    let data = JSON.parse(this.response);
    let inner = '';
    for (let i = 0; i < data.length; i++) {
      inner += '<option value="' + data[i].code + '">' + data[i].code + '</option>';
    }
    document.getElementById('scene-selection').innerHTML = inner;
  }
  xhr.send();
}

populateSceneSelector();

function handleAddResource() {
  let xhr1 = new XMLHttpRequest();
  xhr1.open('POST', 'http://localhost:6060/renders');
  xhr1.onload = function() {
    let data = JSON.parse(this.response);
    activeUuid = data.uuid;
    updateResourceList();

    let xhr2 = new XMLHttpRequest();
    xhr2.open('PUT', 'http://localhost:6060/renders/' + activeUuid + '/scene');
    xhr2.onload = function() {
      // TODO
    }
    xhr2.send(document.getElementById('scene-selection').value);

    let xhr3 = new XMLHttpRequest();
    xhr3.open('PUT', 'http://localhost:6060/renders/' + activeUuid + '/running');
    xhr3.onload = function() {
      // TODO
    }
    xhr3.send('true');
  }
  xhr1.send();
}

document.getElementById("add-resource").addEventListener("click", handleAddResource)

function updateImage() {
  let imgRender = document.getElementById('img-render');
  imgRender.setAttribute('src', 'http://localhost:6060/renders/' + activeUuid + '/image');
  imgRender.addEventListener("click", function() {
    let url = 'http://localhost:6060/renders/' + activeUuid + '/image?random' + new Date().getTime();
    imgRender.setAttribute('src', url);
  });
}

document.getElementById('img-render').addEventListener("click", updateImage);

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
    if (activeUuid != '') {
      updateImage();
    }
  }
  xhr.send();
}

updateResourceList();
