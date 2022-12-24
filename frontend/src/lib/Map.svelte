<svelte:head>
    <!-- <script src="https://unpkg.com/maplibre-gl@2.4.0/dist/maplibre-gl.js"></script> -->
	<link href="https://unpkg.com/maplibre-gl@2.4.0/dist/maplibre-gl.css" rel="stylesheet" on:load={load}/>
</svelte:head>

<script lang="ts">
    import maplibregl from 'maplibre-gl';
    import { onDestroy } from 'svelte';
    import Day from './Day.svelte'

    let map: maplibregl.Map;
    var markers = [];
   
    async function load() {
        map = new maplibregl.Map({
            container: 'map', // container id
            style: 'https://api.maptiler.com/maps/streets-v2/style.json?key=9gPJtf3luQ9Ud3c7AHJ6', // style URL
            center: ["-104.9903", "39.7392"], // starting position [lng, lat]
            zoom: 10 // starting zoom
        });

        const locs =  await fetch("http://192.168.0.100:8090/locations");
        const data = await locs.json();

        data.forEach(point => {
            var marker = new maplibregl.Marker()
                        .setLngLat([point.Lon, point.Lat])
                        .addTo(map)
            
            markers.push(marker)
        });  
    };

    async function updateLocations(e) {
        //First remove current markers
        markers.forEach(marker => {
            marker.remove()
        });
        
        const start = e.detail.start;
        const end = e.detail.end;
        const url = `http://192.168.0.100:8090/locations?start=${start}&end=${end}`
        const locs =  await fetch(url);
        const data = await locs.json();

        markers = []
        data.forEach(point => {
            var marker = new maplibregl.Marker()
                        .setLngLat([point.Lon, point.Lat])
                        .addTo(map)
            
            markers.push(marker)
        });
    };

    onDestroy(() => {
		if (map) map.remove();
	});
</script>

<div id="map"></div>
<div id="date-picker">
    <Day on:dateChange={updateLocations} />
</div>