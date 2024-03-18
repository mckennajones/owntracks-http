<script lang="ts">
    import maplibregl from "maplibre-gl";
    import { onDestroy } from "svelte";
    import DatePicker from "./DatePicker.svelte";

    let map: maplibregl.Map;

    let dateStart;
    let dateEnd;

    function load() {
        map = new maplibregl.Map({
            container: "map", // container id
            // @ts-ignore
            style: window.MAP_STYLE_URL, // style URL
            bounds: [
                [-124.7844079, 49.3457868],
                [-66.9513812, 24.7433195],
            ],
        });
        map.addControl(new maplibregl.NavigationControl({}));

        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.has("start")) {
            dateStart = urlParams.get("start");
        }

        if (urlParams.has("end")) {
            dateEnd = urlParams.get("end");
        }

        map.on("load", function (e) {
            updateLocations({
                detail: {
                    start: dateStart ? new Date(dateStart) : "",
                    end: dateEnd ? new Date(dateEnd) : "",
                },
            });
        });
    }


    window.deletePoint = (id) => {
        if (window.confirm('Are you sure you want to delete this point?')) {
            console.log('Deleting point with id: ' + id)
            let url = window.API_BASE_URL + "/locations/" + id
            fetch(url, {
                method: 'DELETE'
            })
            // TODO surface API error somehow?
        }
    }

    async function updateLocations(e) {
        // @ts-ignore
        let url = window.API_BASE_URL + "/locations"

        if (e.detail.start && e.detail.end) {
            let start: Date = e.detail.start;
            let end: Date = e.detail.end;

            end.setHours(23);
            end.setMinutes(59);
            end.setSeconds(59);

            console.log(start, end);

            url += `?start=${start.toISOString()}&end=${end.toISOString()}`;
        }

        console.log(url);
        const locs = await fetch(url);
        const data = await locs.json();

        if (map.getLayer("heat") != undefined) {
            map.removeLayer("heat");
        }
        if (map.getLayer("points") != undefined) {
            map.removeLayer("points");
        }
        if (map.getSource("main") != undefined) {
            map.removeSource("main")
        }

        map.addSource("main", {
            type: "geojson",
            data: data,
        });

        map.on('mouseenter', 'points', (e) => {
            map.getCanvas().style.cursor = 'pointer';
        });

        map.on('mouseleave', 'points', () => {
            map.getCanvas().style.cursor = 'grab';
        });

        map.on('click', 'points', (e) => {
            console.log(e.features[0]);
            let properties = e.features[0].properties;
            const coordinates = e.features[0].geometry.coordinates.slice();
            
            let id = properties.id;
            let time = new Date(properties.time * 1000);
            
            let description = `
            <div class='popup'>
                <p>ID: ${id}</p>
                <p>Time: ${time.toLocaleString()}</p>

                <input type="button" value="Delete" onclick="deletePoint(${id})" />
            </div>
            `

            

            let popup = new maplibregl.Popup({closeOnClick: true})
                .setLngLat(coordinates)
                .setHTML(description)
                .addTo(map);

            popup.addClassName('popup')

        })

        map.addLayer(
            {
                'id': 'heat',
                'type': 'heatmap',
                'source': 'main',
                'maxzoom': 15,
                'paint': {
                    // Increase the heatmap weight based on frequency and property magnitude
                    'heatmap-weight': [
                        'interpolate',
                        ['linear'],
                        ['get', 'mag'],
                        0,
                        0,
                        6,
                        1
                    ],
                    // Increase the heatmap color weight weight by zoom level
                    // heatmap-intensity is a multiplier on top of heatmap-weight
                    'heatmap-intensity': [
                        'interpolate',
                        ['linear'],
                        ['zoom'],
                        0,
                        1,
                        9,
                        3
                    ],
                    // Color ramp for heatmap.  Domain is 0 (low) to 1 (high).
                    // Begin color ramp at 0-stop with a 0-transparancy color
                    // to create a blur-like effect.
                    'heatmap-color': [
                        'interpolate',
                        ['linear'],
                        ['heatmap-density'],
                        0,
                        'rgba(33,102,172,0)',
                        0.2,
                        'rgb(103,169,207)',
                        0.4,
                        'rgb(209,229,240)',
                        0.6,
                        'rgb(253,219,199)',
                        0.8,
                        'rgb(239,138,98)',
                        1,
                        'rgb(178,24,43)'
                    ],
                    // Adjust the heatmap radius by zoom level
                    'heatmap-radius': [
                        'interpolate',
                        ['linear'],
                        ['zoom'],
                        0,
                        2,
                        9,
                        20
                    ],
                    // Transition from heatmap to circle layer by zoom level
                    'heatmap-opacity': [
                        'interpolate',
                        ['linear'],
                        ['zoom'],
                        7,
                        1,
                        9,
                        0
                    ]
                }
            },
        );

        map.addLayer(
            {
                'id': 'points',
                'type': 'circle',
                'source': 'main',
                'minzoom': 8,
                'paint': {
                    'circle-stroke-color': 'white',
                    'circle-stroke-width': 1,
                    'circle-opacity': [
                        'interpolate',
                        ['linear'],
                        ['zoom'],
                        7,
                        0,
                        8,
                        1
                    ]
                }
            },
        );

        // first set bounds
        let bounds = new maplibregl.LngLatBounds();
        data.features.forEach((feature) => {
            bounds.extend([
                feature.geometry.coordinates[0],
                feature.geometry.coordinates[1],
            ]);
        });

        if (data.features.length != 0) {
            map.fitBounds(bounds, {
                padding: 100,
                duration: 1750,
                essential: true,
                animate: true,
                easing: function (t) {
                    return t * t * t;
                },
            });
        }
    }

    onDestroy(() => {
        if (map) map.remove();
    });
</script>

<svelte:head>
    <!-- <script src="https://unpkg.com/maplibre-gl@2.4.0/dist/maplibre-gl.js"></script> -->
    <link
        href="https://unpkg.com/maplibre-gl@2.4.0/dist/maplibre-gl.css"
        rel="stylesheet"
        on:load={load}
    />
</svelte:head>

<div id="map" />
<div id="date-picker">
    <DatePicker on:dateChange={updateLocations} {dateStart} {dateEnd} />
</div>
