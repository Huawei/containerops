Messenger.options = { extraClasses: 'messenger-fixed messenger-theme-future messenger-on-top'};

function notifyService(){
	return {
		notify : function(msg, type, showtime) {
			 Messenger().post({
		        "message": msg,
		        "type": type,
		        /* success, error, info*/
		        "showCloseButton": true,
		        "hideAfter" : showtime ? showtime : 3
		    });
		}	
	}
}
   
devops.factory('notifyService', [notifyService]);