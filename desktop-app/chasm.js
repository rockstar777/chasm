window.$ = window.jQuery = require('./assets/js/jquery.min.js');
var exec = require('child_process').exec;
var electron = window.require('electron');
var remote = electron.remote;
var BrowserWindow = remote.BrowserWindow;
var PORT = 4567;
var socket = require('socket.io-client')('http://localhost:' + PORT);

// listeners in the GUI
$(document).ready(function() {
	$("#addDropbox").click(function() {
		addService("dropbox");
	});
	$("#addFolder").click(function() {
		addFolder();
	});
	$("#addDrive").click(function() {
		addService("drive");
	});
	$("#cleanChasm").click(function() {
		socket.emit("clean chasm");
	});
	$("#syncChasm").click(function() {
		socket.emit("sync chasm");
	});
});

// socket listeners
socket.on('connect', function(){
});

socket.on('disconnect', function(){});

socket.on('dropbox added', function(data) {
	console.log(data);
	if (data.Success) {
		alert(data.Message);
	} else {
		alert("Error:\n\n" + data.Message);
	}
});

socket.on('drive added', function(data) {
	console.log(data);
	if (data.Success) {
		alert(data.Message);
	} else {
		alert("Error:\n\n" + data.Message);
	}
});

socket.on('folder added', function(data) {
	console.log(data);
	if (data.Success) {
		alert(data.Message);
	} else {
		alert("Error:\n\n" + data.Message);
	}
});
 
socket.on('chasm cleaned', function() {
	alert("Chasm was successfully cleaned!");
});

socket.on('chasm synced', function() {
	alert("Chasm was successfully synced!");
});

socket.on('new event', function(data) {
	$("#event-window").append("<div class=\"event-text " + data.Color + "-text\">" + data.Message + "</div>");
});

// login handlers for the different stores
var addFolder = function() {
	remote.dialog.showOpenDialog({properties: ['createDirectory', 'openDirectory']}, function(paths) {
		socket.emit("add folder", paths[0]);
	});
}

var addService = function(service) {
	var authWindow = new BrowserWindow({ width: 800, height: 600, show: false, 'nodeIntegration': false });

	var driveURL = "https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=713278088797-agohh4u0l5vjscrmn7j0b79i54mtlein.apps.googleusercontent.com&redirect_uri=http://localhost:2000&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.appdata&state=state-token"
	var dropboxURL = "https://www.dropbox.com/1/oauth2/authorize?client_id=zpy424sdnluk9c1&response_type=code&redirect_uri=http://localhost:2000";

	if (service === "dropbox") {
		authWindow.loadURL(dropboxURL);
	} else if (service == "drive") {
		authWindow.loadURL(driveURL);
	}

	authWindow.show();

	function handleCallback (url) {
		var raw_code = /code=([^&]*)/.exec(url) || null;
		var code = (raw_code && raw_code.length > 1) ? raw_code[1] : null;
		var error = /\?error=(.+)$/.exec(url);
		console.log(code);

		if (code || error) {
		// Close the browser if code found or error
			authWindow.destroy();
		}

		// If there is a code, proceed to get token from github
		if (code) {
			if (service === "dropbox") {
				socket.emit("add dropbox", code);
			} else if (service == "drive") {
				socket.emit("add drive", code);
			}
		} else if (error) {
			alert('Oops! Something went wrong and we couldn\'t' +
		  		'log you in using Drive. Please try again.');
		}
	}

	authWindow.webContents.on('will-navigate', function (event, url) {
      handleCallback(url);
    });
}
