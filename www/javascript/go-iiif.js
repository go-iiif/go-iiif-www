window.addEventListener('load', function(e){

    var images_el = document.getElementById("images");
    
    images_el.onchange = function(e){

	var el = e.target;
	var id = el.value;

	display(id);
	return false;
    };
    
    var map = L.map('map', {
	center: [0, 0],
	crs: L.CRS.Simple,
	zoom: 1,
	minZoom: 1,
    });

    var iiif_layer;
    
    var display = function(id){
	
	var info = '/tiles/' + id + '/info.json';

	var opts = {
		'quality': 'color',
	};

	if (iiif_layer){
	    map.removeLayer(iiif_layer);
	}
	
	iiif_layer = L.tileLayer.iiif(info, opts);
	map.addLayer(iiif_layer);    

	var i = document.getElementById("image");
	i.onclick = function(){
		
		var b = document.getElementById("image");
		b.setAttribute("disabled", "disabled");

		b.innerHTML = '<img src="/images/party-parrot.gif" />';

		leafletImage(map, function(err, canvas) {
			
    			if (err){
    				console.log(err);
    				alert("Argh! There was a problem capturing your image");
    				return false;
    			}
			
			var dt = new Date();
			var iso = dt.toISOString();
			var iso = iso.split('T');
			var ymd = iso[0];
			ymd = ymd.replace("-", "", "g");
			
			var bounds = map.getPixelBounds();
			var zoom = map.getZoom();
			
			var pos = [
				bounds.min.x,
				bounds.min.y,
				bounds.max.x,
				bounds.max.y,
				zoom
			];
			
			pos = pos.join("-");
			
			var name = id + "-" + ymd + "-" + pos + ".png";
			
    			canvas.toBlob(function(blob) {
    				saveAs(blob, name);

				var b = document.getElementById("image");
				b.removeAttribute("disabled");
				b.innerText = '📷';
			});
			
    			// window.open(body);
		});

	};

    };

    function onload () {
	
	var data = JSON.parse(this.responseText);
	var count = data.length;
	console.log(data);

	var menu = document.getElementById("images");
	menu.innerHTML = "";
	
	for (var i=0; i < count; i++){
	    var image_name = data[i];

	    var option = document.createElement("option");
	    option.appendChild(document.createTextNode(image_name));

	    menu.appendChild(option);
	}

	display(data[0]);
    };

    var req = new XMLHttpRequest();
    
    req.addEventListener("load", onload);
    req.open("GET", "/images");
    req.send();
    
});
