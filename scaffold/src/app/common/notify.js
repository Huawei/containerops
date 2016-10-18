export function notify(msg,type){
	Messenger().post({
		"message": msg,
		"type": type,
		"showCloseButton": true
	});
}