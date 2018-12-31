// Saves options to chrome.storage
function save_options() {
  	var email = document.getElementById('email').value;
  	
	chrome.storage.sync.set({email: email}, function() {
		var status = document.getElementById('status');
		status.textContent = 'Options saved.';
  	});

	var obj = {"email": email}
        var myJSON = JSON.stringify(obj);
        var req=new XMLHttpRequest();
        req.open("POST","http://localhost:8080/users",true);
        req.setRequestHeader('Content-type', 'application/json');
        req.send(myJSON);
}

document.getElementById('save').addEventListener('click', save_options);
