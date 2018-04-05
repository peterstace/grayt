var activeUuid = '';

document.getElementById("new-resource-button").addEventListener("click", function() {
  var request = new XMLHttpRequest();
  request.open('POST', 'http://localhost:6060/renders');
  request.onload = updateResourceList;
  request.send();
});

function updateResourceList() {
  var request = new XMLHttpRequest();
  request.open('GET', 'http://localhost:6060/renders');
  request.onload = function () {

    var mainWin = document.getElementById("main-window");
    mainWin.innerHTML = '';
    if (activeUuid != '') {
      mainWin.innerHTML += '<input id="input-scene-name" type="text"/>';
      mainWin.innerHTML += '<button id="post-scene-name">POST Scene Name</button>';
      document.getElementById("post-scene-name").addEventListener("click", function() {
        // TODO: Send a post
        //var addr = 'localhost:6060/renders/' + activeUuid + '/scene'
        //alert(document.getElementById("input-scene-name").value);
        //alert(addr);
      });
    }

    var data = JSON.parse(this.response);
    document.getElementById("resource-header").innerHTML = "Resources (" + data.length + ")";
    var resourceList = document.getElementById("resource-list");
    resourceList.innerHTML = "";
    for (var i = 0; i < data.length; i++) {
      if (data[i] == activeUuid) {
        resourceList.innerHTML += '<li class="selected">' + data[i] + "</li>";
      } else {
        resourceList.innerHTML += "<li>" + data[i] + "</li>";
      }
    }
    var listItems = resourceList.childNodes;
    for (var i = 0; i < listItems.length; i++) {
      listItems[i].addEventListener("click", function() {
        activeUuid = this.innerHTML;
        updateResourceList();
      })
    }
  }
  request.send();
}

updateResourceList();
