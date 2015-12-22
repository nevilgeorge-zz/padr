// main.js
$(function() {
  $(document).ready(function() {
    $('.home-button').on('click', function(e) {
      $.ajax({
        url: '/session',
        method: 'POST'
      }).done(function(data) {
        window.location.replace('/' + data);
      })
    });
  });
});
