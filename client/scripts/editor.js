// editor.js
$(function() {
  $(document).ready(function() {
    var domain = window.tempVars.Domain;
    var shortCode = window.tempVars.ShortCode;
    var conn;

    if (window['WebSocket']) {
      conn = new WebSocket('ws://' + domain + '/' + shortCode + '/ws');
      console.log('Connected to websocket!');
    } else {
      alert('Your browser does not support WebSockets.');
    }

    $('#editor').text('');
    var quill = new Quill('#editor');

    $('#editor').on('input', function() {
      var html = quill.getHTML();
      if (typeof conn !== 'undefined' && html !== '') {
        conn.send(html);
      }
      return false;
    });

    if (typeof conn !== 'undefined') {
      conn.onmessage = function(msg) {
        quill.setHTML(msg.data);
      };
    }
  });
});
