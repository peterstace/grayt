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
    var data = JSON.parse(this.response);
    document.getElementById("resource-header").innerHTML = "Resources (" + data.length + ")";
    var resourceList = document.getElementById("resource-list");
    resourceList.innerHTML = "";
    for (var i = 0; i < data.length; i++) {
      if (data[i] == activeUuid) {
        resourceList.innerHTML += '<li>Selected: ' + data[i] + "</li>"; // TODO: Use a CSS class and highlight.
      } else {
        resourceList.innerHTML += "<li>" + data[i] + "</li>";
      }
    }
    // TODO: Add an event to select the UUID.
  }
  request.send();
}

updateResourceList();
