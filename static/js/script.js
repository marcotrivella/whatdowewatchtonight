$('document').ready(function(){
  var larghezza = $('#riga').outerWidth();
  var larghezza_immagine = $('#aa').outerWidth();
  //alert(larghezza + " + " + larghezza_immagine);
  var cont = 0;
  while (larghezza - larghezza_immagine > 0){
    cont++;
    larghezza-=larghezza_immagine;
  }
  //larghezza rimasta 156
  //var margin = Math.trunc(((larghezza-1)/(cont*2)).toFixed(0));
  var margin = (larghezza-1)/(cont*2);
  alert("larghezza residua: 155" + " margine per 12 " + margin*12);
  $(".test").css({"margin-left": margin, "margin-right": margin});
});
