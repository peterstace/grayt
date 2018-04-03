document.getElementById("new-resource-button").addEventListener("click", function() {
  var request = new XMLHttpRequest();
  request.open('POST', 'http://localhost:6060/renders');
  request.send();
  updateResourceList();
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
      resourceList.innerHTML += "<li>" + data[i] + "</li>";
    }
  }
  request.send();
}

updateResourceList();
