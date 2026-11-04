<script lang="ts">
    import maplibregl from "maplibre-gl";
    import "maplibre-gl/dist/maplibre-gl.css";
    import { onDestroy, onMount } from "svelte";
    import DatePicker from "./DatePicker.svelte";

    let map: maplibregl.Map;
    let container: HTMLDivElement;

    let dateStart;
    let dateEnd;

    let count: number | null = null;
    let loading = false;
    let error = "";

    const EMPTY = { type: "FeatureCollection", features: [] } as any;

    const POINT_COLOR = "#4f6bed";

    function esc(v: unknown) {
        return String(v ?? "").replace(/[&<>"']/g, (c) => `&#${c.charCodeAt(0)};`);
    }

    onMount(() => {
        const urlParams = new URLSearchParams(window.location.search);
        dateStart = urlParams.get("start");
        dateEnd = urlParams.get("end");

        map = new maplibregl.Map({
            container,
            // @ts-ignore
            style: window.MAP_STYLE_URL,
            bounds: [
                [-124.7844079, 49.3457868],
                [-66.9513812, 24.7433195],
            ],
            attributionControl: false,
        });
        // @ts-ignore
        if (import.meta.env.DEV) window.__map = map; // handy for debugging / e2e
        map.addControl(new maplibregl.NavigationControl({}), "bottom-right");
        map.addControl(
            new maplibregl.AttributionControl({ compact: true }),
            "bottom-left",
        );

        map.on("load", () => {
            setupLayers();
            setupInteractions();
            updateLocations({
                detail: {
                    start: dateStart ? new Date(dateStart) : "",
                    end: dateEnd ? new Date(dateEnd) : "",
                },
            });
        });
    });

    function setupLayers() {
        map.addSource("main", { type: "geojson", data: EMPTY });

        // faint route line; built from time-sorted points in updateLocations
        map.addSource("route", { type: "geojson", data: EMPTY });
        map.addLayer({
            id: "route",
            type: "line",
            source: "route",
            layout: { "line-join": "round", "line-cap": "round" },
            paint: {
                "line-color": "#9fb4ff",
                "line-opacity": ["interpolate", ["linear"], ["zoom"], 6, 0, 9, 0.35],
                "line-width": ["interpolate", ["linear"], ["zoom"], 8, 1, 16, 3],
            },
        });

        map.addLayer({
            id: "heat",
            type: "heatmap",
            source: "main",
            maxzoom: 15,
            paint: {
                "heatmap-weight": 0.25,
                "heatmap-intensity": ["interpolate", ["linear"], ["zoom"], 0, 0.4, 9, 1],
                "heatmap-color": [
                    "interpolate",
                    ["linear"],
                    ["heatmap-density"],
                    0, "rgba(140,170,255,0)",
                    0.2, "rgba(140,170,255,0.55)",
                    0.45, "rgb(90,130,245)",
                    0.7, "rgb(55,80,220)",
                    1, "rgb(25,30,130)",
                ],
                "heatmap-radius": ["interpolate", ["linear"], ["zoom"], 0, 3, 9, 22, 14, 40],
                "heatmap-opacity": ["interpolate", ["linear"], ["zoom"], 7, 1, 10, 0],
            },
        });

        // soft glow under the points
        map.addLayer({
            id: "points-glow",
            type: "circle",
            source: "main",
            minzoom: 8,
            paint: {
                "circle-color": POINT_COLOR,
                "circle-radius": ["interpolate", ["linear"], ["zoom"], 8, 6, 16, 16],
                "circle-blur": 1,
                "circle-opacity": ["interpolate", ["linear"], ["zoom"], 8, 0, 10, 0.35],
            },
        });

        map.addLayer({
            id: "points",
            type: "circle",
            source: "main",
            minzoom: 8,
            paint: {
                "circle-color": POINT_COLOR,
                "circle-radius": ["interpolate", ["linear"], ["zoom"], 8, 2.5, 16, 7],
                "circle-stroke-color": "rgba(10,12,20,0.85)",
                "circle-stroke-width": 1,
                "circle-opacity": ["interpolate", ["linear"], ["zoom"], 8, 0, 10, 1],
                "circle-stroke-opacity": ["interpolate", ["linear"], ["zoom"], 8, 0, 10, 1],
            },
        });

        map.addLayer({
            id: "points-hover",
            type: "circle",
            source: "main",
            minzoom: 8,
            filter: ["==", ["get", "id"], -1],
            paint: {
                "circle-color": "rgba(0,0,0,0)",
                "circle-radius": ["interpolate", ["linear"], ["zoom"], 8, 6, 16, 12],
                "circle-stroke-color": "#ffffff",
                "circle-stroke-width": 2,
            },
        });
    }

    function setupInteractions() {
        map.on("mouseenter", "points", () => {
            map.getCanvas().style.cursor = "pointer";
        });
        map.on("mousemove", "points", (e) => {
            map.setFilter("points-hover", ["==", ["get", "id"], e.features[0].properties.id]);
        });
        map.on("mouseleave", "points", () => {
            map.getCanvas().style.cursor = "";
            map.setFilter("points-hover", ["==", ["get", "id"], -1]);
        });

        map.on("click", "points", (e) => {
            const feature = e.features[0];
            const p = feature.properties;
            // @ts-ignore
            const coordinates = feature.geometry.coordinates.slice();
            const time = new Date(p.time * 1000);
            const fmt = (v, digits = 0) =>
                v === undefined || v === null || v === "null" ? "–" : Number(v).toFixed(digits);

            const html = `
            <div class="popup">
                <div class="popup-time">${esc(
                    time.toLocaleString(undefined, { dateStyle: "medium", timeStyle: "short" }),
                )}</div>
                <div class="popup-stats">
                    <div><span>${fmt(p.altitude)}</span> m</div>
                    ${p.battery !== undefined && p.battery !== "null" ? `<div><span>${esc(p.battery)}</span>% bat</div>` : ""}
                </div>
                <div class="popup-foot">
                    <small>#${esc(p.id)}</small>
                    <button type="button" class="popup-delete">Delete</button>
                </div>
                <div class="popup-error"></div>
            </div>`;

            const popup = new maplibregl.Popup({ closeOnClick: true, offset: 10 })
                .setLngLat(coordinates)
                .setHTML(html)
                .addTo(map);

            const el = popup.getElement();
            el.querySelector(".popup-delete").addEventListener("click", async (ev) => {
                const btn = ev.currentTarget as HTMLButtonElement;
                if (!window.confirm("Are you sure you want to delete this point?")) return;
                btn.disabled = true;
                try {
                    // @ts-ignore
                    const res = await fetch(`${window.API_BASE_URL}/locations/${p.id}`, {
                        method: "DELETE",
                    });
                    if (!res.ok) throw new Error(`HTTP ${res.status}`);
                    popup.remove();
                    await updateLocations(currentRange(), false);
                } catch (err) {
                    btn.disabled = false;
                    el.querySelector(".popup-error").textContent = `Delete failed: ${err.message}`;
                }
            });
        });
    }

    let lastDetail: any = { start: "", end: "" };
    function currentRange() {
        return { detail: lastDetail };
    }

    async function updateLocations(e, fit = true) {
        lastDetail = {
            start: e.detail.start ? new Date(e.detail.start) : "",
            end: e.detail.end ? new Date(e.detail.end) : "",
        };

        // @ts-ignore
        let url = window.API_BASE_URL + "/locations";

        if (lastDetail.start && lastDetail.end) {
            const start: Date = new Date(lastDetail.start);
            const end: Date = new Date(lastDetail.end);
            end.setHours(23, 59, 59);
            url += `?start=${start.toISOString()}&end=${end.toISOString()}`;
        }

        loading = true;
        error = "";
        let data;
        try {
            const res = await fetch(url);
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            data = await res.json();
        } catch (err) {
            error = `Couldn't load locations (${err.message})`;
            loading = false;
            return;
        }
        loading = false;

        const features = data.features ?? [];
        count = features.length;

        (map.getSource("main") as maplibregl.GeoJSONSource).setData(data);

        // split into trips wherever there is a long gap, so we don't draw
        // misleading straight lines between unrelated visits
        const MAX_GAP_S = 30 * 60;
        const sorted = [...features].sort((a, b) => a.properties.time - b.properties.time);
        const trips: number[][][] = [];
        let prev = null;
        for (const f of sorted) {
            if (!prev || f.properties.time - prev.properties.time > MAX_GAP_S) trips.push([]);
            trips[trips.length - 1].push(f.geometry.coordinates);
            prev = f;
        }
        (map.getSource("route") as maplibregl.GeoJSONSource).setData({
            type: "FeatureCollection",
            features: trips
                .filter((t) => t.length > 1)
                .map((coordinates) => ({
                    type: "Feature",
                    properties: {},
                    geometry: { type: "LineString", coordinates },
                })),
        } as any);

        if (fit && features.length != 0) {
            const bounds = new maplibregl.LngLatBounds();
            features.forEach((f) => bounds.extend([f.geometry.coordinates[0], f.geometry.coordinates[1]]));
            map.fitBounds(bounds, {
                padding: 100,
                maxZoom: 15,
                duration: 1750,
                essential: true,
                easing: (t) => t * t * t,
            });
        }
    }

    onDestroy(() => {
        if (map) map.remove();
    });
</script>

<div id="map" bind:this={container} />

<div id="date-picker" class="panel">
    <DatePicker on:dateChange={updateLocations} {dateStart} {dateEnd} />
    <div class="status">
        {#if error}
            <span class="status-error">{error}</span>
        {:else if loading}
            Loading…
        {:else if count !== null}
            {count.toLocaleString()} point{count === 1 ? "" : "s"}
        {/if}
    </div>
</div>
