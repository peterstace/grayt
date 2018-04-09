var activeUuid = '';

document.getElementById("new-resource-button").addEventListener("click", function() {
  let request = new XMLHttpRequest();
  request.open('POST', 'http://localhost:6060/renders');
  request.onload = updateResourceList;
  request.send();
});

function updateResourceList() {
  let request = new XMLHttpRequest();
  request.open('GET', 'http://localhost:6060/renders');
  request.onload = function () {

    // TODO: Remove event handlers.
    // TODO: Gray out inputs and buttons if nothing active.
    if (activeUuid != '') {

      document.getElementById("put-scene-name").addEventListener("click", function() {
        let request = new XMLHttpRequest();
        let url  = 'http://localhost:6060/renders/' + activeUuid + '/scene';
        request.open('PUT', url);
        request.onload = function() {
          alert('Status: ' + this.status + '\nResponse: ' + this.response);
        };
        request.send(document.getElementById('input-scene-name').value);
      });
      
      document.getElementById('put-running').addEventListener("click", function() {
        let request = new XMLHttpRequest();
        let url  = 'http://localhost:6060/renders/' + activeUuid + '/running';
        request.open('PUT', url);
        request.onload = function() {
          alert('Response: ' + this.response);
        };
        request.send('true');
      });

      document.getElementById('img-render').setAttribute('src', 'http://localhost:6060/renders/' + activeUuid + '/image');
      document.getElementById('img-render').addEventListener("click", function() {
        document.getElementById('img-render').setAttribute('src', 'http://localhost:6060/renders/' + activeUuid + '/image?random' + new Date().getTime());
      });
    };

    let data = JSON.parse(this.response);
    document.getElementById("resource-header").innerHTML = "Resources (" + data.length + ")";
    let resourceList = document.getElementById("resource-list");
    resourceList.innerHTML = "";
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
  request.send();
}

updateResourceList();
