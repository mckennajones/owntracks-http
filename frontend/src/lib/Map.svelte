<svelte:head>
    <!-- <script src="https://unpkg.com/maplibre-gl@2.4.0/dist/maplibre-gl.js"></script> -->
	<link href="https://unpkg.com/maplibre-gl@2.4.0/dist/maplibre-gl.css" rel="stylesheet" on:load={load}/>
</svelte:head>

<script lang="ts">
    import maplibregl from 'maplibre-gl';
    import { onDestroy } from 'svelte';
    import DatePicker from './DatePicker.svelte'

    let map: maplibregl.Map;
    var markers = [];
   
    async function load() {
        map = new maplibregl.Map({
            container: 'map', // container id
            style: 'https://api.maptiler.com/maps/streets-v2/style.json?key=9gPJtf3luQ9Ud3c7AHJ6', // style URL
            // center: ["-104.9903", "39.7392"], // starting position [lng, lat]
            // zoom: 10 // starting zoom
        });
        map.addControl(new maplibregl.NavigationControl({}));

        const locs =  await fetch("http://192.168.0.100:8090/locations");
        const data = await locs.json();

        let bounds = new maplibregl.LngLatBounds();
        data.forEach(point => {
            var marker = new maplibregl.Marker()
                        .setLngLat([point.Lon, point.Lat])
                        .addTo(map)
            
            markers.push(marker)
            bounds.extend([point.Lon, point.Lat])
        });
            
        map.fitBounds(bounds, {
            padding: 100
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
        console.log(url)
        const locs =  await fetch(url);
        const data = await locs.json();

        markers = []
        let bounds = new maplibregl.LngLatBounds();
        data.forEach(point => {
            var marker = new maplibregl.Marker()
                        .setLngLat([point.Lon, point.Lat])
                        .addTo(map)
            
            markers.push(marker)
            bounds.extend([point.Lon, point.Lat])
        });

        if (data.length != 0) {
            map.fitBounds(bounds, {
                padding: 100
            });
        }
    };

    onDestroy(() => {
		if (map) map.remove();
	});
</script>

<div id="map"></div>
<div id="date-picker">
    <DatePicker on:dateChange={updateLocations} />
</div>