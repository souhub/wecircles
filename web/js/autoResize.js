$(function () {
  var $textarea = $(".autoScaleTextArea");
  var lineHeight = parseInt($textarea.css("lineHeight"));
  $textarea.on("input", function (e) {
    var lines = ($(this).val() + "\n").match(/\n/g).length;
    $(this).height(lineHeight * lines);
  });
});

$(function () {
  $("textarea.auto-resize").on("change keyup keydown paste cut", function () {
    if ($(this).outerHeight() > this.scrollHeight) {
      $(this).height(1);
    }
    while ($(this).outerHeight() < this.scrollHeight) {
      $(this).height($(this).height() + 1);
    }
  });
});
