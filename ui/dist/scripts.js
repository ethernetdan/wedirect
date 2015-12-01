'use strict';

var $input = $('.info-input');

function getDate() {
  var date = new Date();
  var month = date.getUTCMonth() + 1;
  var day = date.getUTCDate();
  var year = date.getUTCFullYear();

  return month + "/" + day + "/"  + year;
}

function submit() {
  if ($input.val()) {

    $('.info-submit').removeClass('failed');
    $.ajax({
      url: '/set',
      method: 'POST',
      data: {domain: $input.val()}
    })
    .done(function(data) {
      $('.history p').eq(0).find('span').eq(2).text('inactive');
      $('.history').prepend(
        '<p class="fade-in"><span class="history-date">'+getDate()+'</span> <span class="history-link">' + $input.val() +'</span> <span class="history-status">live</span></p>'
        );
      $input.val('');
    })
    .fail(function() {
      $('.info-submit').addClass('failed');
    });
  }
}

$(document).ready(function() {

  $input.on('input',function() {
    $('.info-submit').addClass('active');

    if ($(this).val() === '')
      $('.info-submit').removeClass('active');
  });

  $('.info-submit').click(submit);

  $input.keyup(function(event){
    if(event.keyCode == 13)
      submit()
  });
});
