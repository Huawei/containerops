function more($location){
	return {
		show : function(next) {
			 $(window).scroll(function() {
			 	if($location.path() == "/component"){
			 		if($(window).scrollTop() + $(window).height() == $(document).height()) {
			 			next();
			 		}
			 	}
			});
		}	
	}
}
   
devops.factory('more', ['$location', more]);