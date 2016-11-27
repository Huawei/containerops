import { Injectable } from '@angular/core';

declare var Messenger: any;
Messenger.options = { extraClasses: 'messenger-fixed messenger-theme-future messenger-on-top', theme: 'air' };

Injectable()
export class NotifyService {
	notify(msg: String, type: String, showtime: number = 3) {
		 Messenger().post({
	        "message": msg,
	        "type": type,
	        /* success, error, info*/
	        "showCloseButton": true,
	        "hideAfter" : showtime
	    });
	}
}
