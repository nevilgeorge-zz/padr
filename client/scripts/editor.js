// editor.js

// model
var Operation = function(start, count, chars, type) {
  this.start = start;
  this.count = count;
  this.chars = chars;
  this.type = type;
}

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
      var op, action, start, chars;
      if (delta.ops.length == 1) {
        action = delta.ops[0]
        if ('insert' in action) {
          op = new Operation(0, action.insert.length, action.insert, 'insert');
        } else if ('delete' in action) {
          op = new Operation(0, action.delete, '', 'delete');
        }
      } else {
        start = delta.ops[0].retain || 0;
        action = delta.ops[1];
        chars = action.insert;
        if ('insert' in action) {
          op = new Operation(start, action.insert.length, chars, 'insert');
        } else if ('delete' in action) {
          op = new Operation(start, action.delete, '', 'delete');
        }
      }
      console.log(op)

      if (source === 'user' && typeof conn !== 'undefined' && typeof op !== 'undefined') {
        conn.send(JSON.stringify(op));
      }
    });

    if (typeof conn !== 'undefined') {
      conn.onmessage = function(msg) {
        var range = quill.getSelection();
        quill.setText(msg.data);
        quill.setSelection(range);
      };
    }
  });
});
