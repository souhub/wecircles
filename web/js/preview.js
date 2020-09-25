// function imgPreView(event) {
//   var file = event.target.files[0];
//   var reader = new FileReader();
//   var preview = document.getElementById("preview");
//   var previewImage = document.getElementById("previewImage");

//   if (previewImage != null) {
//     preview.removeChild(previewImage);
//   }
//   reader.onload = function (event) {
//     var img = document.createElement("img");
//     img.setAttribute("src", reader.result);
//     img.setAttribute("id", "previewImage");
//     preview.appendChild(img);
//   };

//   reader.readAsDataURL(file);
// }

function previewImage(obj) {
  var fileReader = new FileReader();
  fileReader.onload = function () {
    document.getElementById("preview").src = fileReader.result;
  };
  fileReader.readAsDataURL(obj.files[0]);
}
