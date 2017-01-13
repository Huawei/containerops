function more(){
	return {
		show : function(next) {
			 $(window).scroll(function() {
			   if($(window).scrollTop() + $(window).height() == $(document).height()) {
			      next();
			   }
			});
		}	
	}
}
   
devops.factory('more', [more]);