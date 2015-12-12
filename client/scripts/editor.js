// editor.js
$(function() {
  $(document).ready(function() {
    console.log('Ready!');
    $('.editable').text('');
    var editor = new MediumEditor('.editable', {
      placeholder: {
        text: ''
      }
    });
  });
})
