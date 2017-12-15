$('document').ready(function(){
  var larghezza = $('#riga').width();
  var larghezza_immagine = 164
  var cont = 0;
  while (larghezza - larghezza_immagine > 0){
    cont++;
    larghezza-=larghezza_immagine;
  }
  var margin = Math.trunc(((larghezza-1)/(cont*2)).toFixed(0));
  $(".test").css({"margin-left": margin, "margin-right": margin});
});
