
var email = "";
var password = "";
var songTitle = "";

/*
// function to int 
function init() {
	// goal of this function is to go to backend, trigger an endpoint with an
	// http request and get a response
}
// html populate list with get
function updateSongTitles(p) {
	songTitle = p;
	document.getElementByID("SongList")
// send info in a param in the url
// serve get parameters in golang
}

var json={'a': 'First', 'b': 'Second', 'c': 'Third'};
function makeUL(json) {
    // Create the list element:
    var list = document.createElement('ul');

    for(var i = 0; i < Object.keys(json).length; i++) {
        // Create the list item:
        var item = document.createElement('li');

        // Set its contents:
        item.appendChild(document.createTextNode(Object.values(json)[i]));

        // Add it to the list:
        list.appendChild(item);
    }

    // Finally, return the constructed list:
    return list;
}

// Add the contents of json to #foo:
document.getElementById('foo').appendChild(makeUL(json));
*/
// iter through songs
// for each song update corresponding list element
// for loop with updating (".innerText)" inside loop)
// issues inc how to return a list from go function to javascript
// also how do you iterate through list in js
//$.getJSON('/homepage', function(songs) {

// July 29 currently trying to figure out how to parse json file 
//either from html/javascript to load onto html “/homepage” but 
//as there are issues loading locally-held files from javascript 
//(from security issues that could cause) I may have to do it 
//from a server XMLHttpRequest()
/*
function getJSON() {     	 
	//for(let i = 1; i <= 3; i++) {
		var xobj = new XMLHttpRequest();
		//xobj.overrideMimeType("application/json");
		//xobj.open('GET', "../backend/artist.json", true);
		//xobj.responseType = 'json';
		xobj.onreadystatechange = function () {
          		if (xobj.readyState == 4 && xobj.status == "200") {
            			document.getElementById("Song1").innerHTML = "Nas";	
            			document.getElementById("Song2").innerHTML = this.ResponseText;	
            			document.getElementById("Song3").innerHTML = this.ResponseText;	
				//callback(xobj.responseText);
          		}
			xobj.open('GET', "backend/artist.json", true);
			xobj.send();
		};
		//var Song = JSON.Parse(
		//var strSong = "Song" + string(i)
		//var Song = document.getElementById(strSong).innerText;
		//document.getElementById(strSong).innerHTML = Song;
		//$(strSong).html(songs)
	//}
        
        //$(".mypanel").html(text);
    });
    */
