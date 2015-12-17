// editor.js
$(function() {
  $(document).ready(function() {
    var conn;

    if (window['WebSocket']) {
      conn = new WebSocket('ws://localhost:8080/ws');
      console.log('New websocket!');
    } else {
      alert('Your browser does not support WebSockets.');
    }

    $('.editable').text('');
    var editor = new MediumEditor('.editable', {
      placeholder: {
        text: ''
      }
    });

    var editorDiv = $('.editable');
    $('.editable').on('input', function() {
      var text = editorDiv.text();
      if (typeof conn !== 'undefined' && text !== '') {
        console.log(text);
        conn.send(text);
      }
      return false;
    });

    if (typeof conn !== 'undefined') {
      conn.onmessage = function(msg) {
        $('.editable').text(msg.data);
      };
    }
  });
})
