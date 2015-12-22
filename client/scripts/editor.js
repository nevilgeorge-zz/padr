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

    quill.on('text-change', function(delta, source) {
      var html = quill.getHTML();
      if (source === 'user' && typeof conn !== 'undefined' && html !== '') {
        conn.send(html);
      }
      return false;
    });

    if (typeof conn !== 'undefined') {
      conn.onmessage = function(msg) {
        var range = quill.getSelection();
        quill.setHTML(msg.data);
        quill.setSelection(range);
      };
    }
  });
});
