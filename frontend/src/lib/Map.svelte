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

    let dateStart
    let dateEnd
   
    async function load() {
        map = new maplibregl.Map({
            container: 'map', // container id
            style: 'https://api.maptiler.com/maps/streets-v2/style.json?key=9gPJtf3luQ9Ud3c7AHJ6', // style URL
            // center: ["-104.9903", "39.7392"], // starting position [lng, lat]
            // zoom: 10 // starting zoom
        });
        map.addControl(new maplibregl.NavigationControl({}));

        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.has('start')) {
            dateStart = urlParams.get('start')
        }

        if (urlParams.has('end')) {
            dateEnd = urlParams.get('end')
        }

        updateLocations({
            detail: {
                start: dateStart ? new Date(dateStart) : '',
                end: dateEnd ? new Date(dateEnd): ''
            }
        })

    };

    async function updateLocations(e) {
        let url = 'http://192.168.0.100:8090/locations'
        //First remove current markers
        markers.forEach(marker => {
            marker.remove()
        });
        
        if (e.detail.start && e.detail.end) {
            let start: Date = e.detail.start;
            let end: Date = e.detail.end;

            if (start.toISOString() == end.toISOString()) {
                end.setHours(23)
                end.setMinutes(59)
                end.setSeconds(59)
            }

            console.log(start, end)
            
            url += `?start=${start.toISOString()}&end=${end.toISOString()}`
        }
        
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
    <DatePicker on:dateChange={updateLocations} dateStart={dateStart} dateEnd={dateEnd}/>
</div>