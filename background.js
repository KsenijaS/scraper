function getUrl() {
    	chrome.tabs.query({'active': true, 'lastFocusedWindow': true}, function (tabs) {
		var url = tabs[0].url;
		alert(url);
        	chrome.storage.sync.get('email', function (result) {
                	var email = result.email;
			var obj = {"url": url, "email": email}
                	var myJSON = JSON.stringify(obj);
                	var req=new XMLHttpRequest();
                	req.open("POST","http://localhost:8080/urls",true);
                	req.setRequestHeader('Content-type', 'application/json');
                	req.send(myJSON);
        	});
    	});
}

chrome.browserAction.onClicked.addListener(getUrl)
