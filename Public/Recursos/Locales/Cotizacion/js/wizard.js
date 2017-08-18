$(document).ready(function () {
    //Initialize tooltips
    // $('.nav-tabs > li a[title]').tooltip();
    
    //Wizard
    $('a[data-toggle="tab"]').on('click', function (e) {

        var $target = $(e.target);
        console.log($target);
    
        if ( $target.parent().parent().parent().hasClass('disabled') || $target.parent().parent().hasClass('disabled') ) {
            console.log("No permitido");
        }else{
            $('div[class="tab-pane active"]').removeClass("active");
            $('li[role="presentation"]').removeClass("active");
            if( $target.hasClass('round-tab') ){
                $target.parent().parent().addClass("active");
                var id = $target.parent().attr("href");                
            }else if( $target.parent().hasClass('round-tab') ){
                $target.parent().parent().parent().addClass("active");                
                var id = $target.parent().parent().attr("href");
            }else{
                $target.parent().addClass("active");                
                var id = $target.attr("href");
            }
            $(id).addClass("active");            
            console.log("Permitido");
        }
    });

    $(".next-step").click(function (e) {
        var $active = $('.wizard .nav-tabs li.active');
        $active.removeClass('active');
        $active.next().removeClass('disabled').addClass('active');
        nextTab($active);

    });

    $(".prev-step").click(function (e) {
        var $active = $('.wizard .nav-tabs li.active');
        $active.removeClass('active');
        $active.prev().removeClass('disabled').addClass('active');
        prevTab($active);

    });
});

function nextTab(elem) {
    $(elem).next().find('a[data-toggle="tab"]').click();
}
function prevTab(elem) {
    $(elem).prev().find('a[data-toggle="tab"]').click();
}


//according menu

$(document).ready(function()
{
    //Add Inactive Class To All Accordion Headers
    $('.accordion-header').toggleClass('inactive-header');
	
	//Set The Accordion Content Width
	var contentwidth = $('.accordion-header').width();
	$('.accordion-content').css({});
	
	//Open The First Accordion Section When Page Loads
	$('.accordion-header').first().toggleClass('active-header').toggleClass('inactive-header');
	$('.accordion-content').first().slideDown().toggleClass('open-content');
	
	// The Accordion Effect
	$('.accordion-header').click(function () {
		if($(this).is('.inactive-header')) {
			$('.active-header').toggleClass('active-header').toggleClass('inactive-header').next().slideToggle().toggleClass('open-content');
			$(this).toggleClass('active-header').toggleClass('inactive-header');
			$(this).next().slideToggle().toggleClass('open-content');
		}
		
		else {
			$(this).toggleClass('active-header').toggleClass('inactive-header');
			$(this).next().slideToggle().toggleClass('open-content');
		}
	});
	
	return false;
});