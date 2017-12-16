$('document').ready(function(){
  var larghezza = $('#contenitore').width();
  var larghezza_immagine = 164
  var cont = 0;
  while (larghezza - larghezza_immagine > 0){
    cont++;
    larghezza-=larghezza_immagine;
  }
  var margin = Math.trunc(((larghezza-1)/(cont*2)).toFixed(0));
  $(".locandinaFilm").css({"margin-left": margin, "margin-right": margin});

  //<img class="img-fluid img-thumbnail locandinaFilm" src="https://image.tmdb.org/t/p/w154/kqjL17yufvn9OVLyXYpvtyrFfak.jpg">
});
