var request = new XMLHttpRequest();
request.open('GET', 'http://localhost:6060/renders', true);
request.onload = function () {
  var data = JSON.parse(this.response);
  console.log(data);
}
request.send();
