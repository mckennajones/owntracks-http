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
            style: "https://api.maptiler.com/maps/streets-v2/style.json?key=9gPJtf3luQ9Ud3c7AHJ6", // style URL
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

    async function updateLocations(e) {
        let url = "http://192.168.0.100:8090/locations";

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

        map.addLayer({
            id: "points",
            type: "circle",
            source: "main",
            paint: {
                "circle-radius": 6,
                "circle-color": "#B42222",
            },
            filter: ["==", "$type", "Point"],
        });

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
