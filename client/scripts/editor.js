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

    $('.editable').text('');
    var editor = new MediumEditor('.editable', {
      placeholder: {
        text: ''
      }
    });

    var editorDiv = $('.editable');
    $('.editable').on('input', function() {
      var html = editorDiv.html();
      if (typeof conn !== 'undefined' && html !== '') {
        conn.send(html);
      }
      return false;
    });

    if (typeof conn !== 'undefined') {
      conn.onmessage = function(msg) {
        $('.editable').html($.parseHTML(msg.data));
      };
    }
  });
});
