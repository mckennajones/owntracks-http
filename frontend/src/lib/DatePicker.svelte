<script lang="ts">
  	import { createEventDispatcher } from 'svelte';
    import { DateInput } from 'date-picker-svelte';
    import { DatePicker, DatePickerInput } from "carbon-components-svelte";
    import "carbon-components-svelte/css/g90.css";


    const dispatch = createEventDispatcher();

    function datePickerChange(event: any) {
        console.log(event)

        const detail = event.detail
        
        if (detail.selectedDates.length != 2) {
            return
        }

        const from: Date = detail.selectedDates[0]
        const to: Date = detail.selectedDates[1]

        if (from.toISOString() == to.toISOString()) {
            console.log("Same")
            to.setHours(23)
            to.setMinutes(59)
            to.setSeconds(59)
        }

        console.log(from, to)
        dispatch('dateChange', {
			start: from.toISOString(),
            end: to.toISOString()
		});
    }
    
</script>

<DatePicker datePickerType="range" on:change={datePickerChange}>
    <DatePickerInput hideLabel placeholder="All" />
    <DatePickerInput hideLabel placeholder="time" />
</DatePicker>
