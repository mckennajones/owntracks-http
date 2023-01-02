<script lang="ts">
  	import { createEventDispatcher } from 'svelte';
    import { DateInput } from 'date-picker-svelte';
    import { DatePicker, DatePickerInput } from "carbon-components-svelte";
    import "carbon-components-svelte/css/g90.css";


    const dispatch = createEventDispatcher();

    let day: any;
    let month: any;
    let year: any;

    function dateChange() {
        console.log(year, month, day)
        let start: Date;
        let end: Date;
        
        if (year && month >= 0 && day) {
            console.log("all set")
            start = new Date(year, month, day, 0)
            end = new Date(year, month, day, 23, 59, 59)
        } else if (year && month >= 0) {
            console.log("year and month set")
            // get month length by setting date to 0
            let tempDate = new Date(year, month + 1, 0)
            start = new Date(year, month, 0)
            end = new Date(year, month, tempDate.getDate())
        } else if (year) {
            console.log("year set")
            start = new Date(year, 0, 1)
            end = new Date(year, 11, 31)
        }

        console.log(start, end)

        dispatch('dateChange', {
			start: start.toISOString(),
            end: end.toISOString()
		});
    }

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
    <DatePickerInput hideLabel placeholder="mm/dd/yyyy" />
    <DatePickerInput hideLabel placeholder="mm/dd/yyyy" />
</DatePicker>
