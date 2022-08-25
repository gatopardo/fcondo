$(function() {
    //$(document).foundation();
    // Hide any messages after a few seconds
    hideFlash();
});

function hideFlash(rnum)
{    
    if (!rnum) rnum = '0';
 
    _.delay(function() {
        $('.alert-box-fixed' + rnum).fadeOut(300, function() {
            $(this).css({"visibility":"hidden",display:'block'}).slideUp();
            
            var that = this;
            
            _.delay(function() { that.remove(); }, 400);
        });
    }, 4000);
}

function showFlash(obj)
{
    $('#flash-container').html();
    $(obj).each(function(i, v) {
        var rnum = _.random(0, 100000);
		var message = '<div id="flash-message" class="alert-box-fixed'
		+ rnum + ' alert-box-fixed alert alert-dismissible '+v.cssclass+'">'
		+ '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>'
		+ v.message + '</div>';
        $('#flash-container').prepend(message);
        hideFlash(rnum);
    });
}

function flashError(message) {
	var flash = [{Class: "alert-danger", Message: message}];
	showFlash(flash);
}

function flashSuccess(message) {
	var flash = [{Class: "alert-success", Message: message}];	
	showFlash(flash);
}

function flashNotice(message) {
	var flash = [{Class: "alert-info", Message: message}];
	showFlash(flash);
}

function flashWarning(message) {
	var flash = [{Class: "alert-warning", Message: message}];
	showFlash(flash);
}


$('.dropdown-menu a.dropdown-toggle').on('click', function(e) {
  if (!$(this).next().hasClass('show')) {
    $(this).parents('.dropdown-menu').first().find('.show').removeClass("show");
  }
  var $subMenu = $(this).next(".dropdown-menu");
  $subMenu.toggleClass('show');


  $(this).parents('li.nav-item.dropdown.show').on('hidden.bs.dropdown', function(e) {
    $('.dropdown-submenu .show').removeClass("show");
  });


  return false;
});

function getNameLang(){
     var xm;
     var form = document.querySelector("form")
     
     if (window.XMLHttpRequest){
        xm=new XMLHttpRequest();
     }else  {
        xm=new ActiveXObject('Microsoft.XMLHTTP');
     }

     xm.onreadystatechange = function() {
         if(xm.readyState == 4) {
              console.log(xm.response);
              alert(xm.response);
       }
     }

     xm.open("POST","/perscuenta",true);
     xm.send(new FormData(form));
   }

/*
function carousel() {
    var i;
     alert("hola 2");
    var x = document.getElementsByClassName("mySlides");
    for (i = 0; i < x.length; i++) {
       x[i].style.display = "none";  
    }
    myIndex++;
    if (myIndex > x.length) {myIndex = 1}    
    x[myIndex-1].style.display = "block";  
    setTimeout(carousel, 2000); // Change image every 2 seconds
}
*/


