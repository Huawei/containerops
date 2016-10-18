export function resizeWidget(){
	var $widgets = $('.widget');

    $widgets.on("fullscreen.widgster", function(){
    	$('.content-wrap').css({
         	'-webkit-transform': 'none',
            '-ms-transform': 'none',
            'transform': 'none',
            'margin': 0,
            'z-index': 2
        });
        $(".treeview").css("max-height",window.screen.height * 2 / 3);
        $(".treeview").css("overflow","auto");
    }).on("restore.widgster closed.widgster", function(){
        $('.content-wrap').css({
        	'-webkit-transform': '',
            '-ms-transform': '',
            'transform': '',
            'margin': '',
            'z-index': ''
        });
        $(".treeview").css("max-height","");
        $(".treeview").css("overflow","");
    });

    $widgets.widgster();
}