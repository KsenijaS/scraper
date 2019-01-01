function load() {
    	var email = "";
    	chrome.storage.sync.get('email', function (result) {
        	email = result.email;
		console.log(result.email);
    		alert(email);
	});
}

document.addEventListener('DOMContentLoaded', load);
