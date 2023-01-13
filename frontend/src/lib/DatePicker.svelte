<script lang="ts">
  	import { createEventDispatcher } from 'svelte';
    import { DateInput } from 'date-picker-svelte';
    import { DatePicker, DatePickerInput } from "carbon-components-svelte";
    import "carbon-components-svelte/css/g90.css";

    export let dateStart;
    export let dateEnd;

    const dispatch = createEventDispatcher();

    function datePickerChange(event: any) {
        const detail = event.detail
        console.log(event)
        if (detail.selectedDates.length != 2) {
            return
        }

        const from: Date = detail.selectedDates[0]
        const to: Date = detail.selectedDates[1]

        let params = new URLSearchParams(window.location.search);
        
        params.set('start', detail.dateStr.from)
        params.set('end', detail.dateStr.to)

        // decode to avoid weird url encoding
        window.history.replaceState({}, '', decodeURIComponent(`${location.pathname}?${params}`));

        dispatch('dateChange', {
			start: from,
            end: to
		});
    }
    
</script>

<DatePicker datePickerType="range" on:change={datePickerChange} valueFrom={dateStart} valueTo={dateEnd}>
    <DatePickerInput hideLabel placeholder="All"  />
    <DatePickerInput hideLabel placeholder="time" />
</DatePicker>
