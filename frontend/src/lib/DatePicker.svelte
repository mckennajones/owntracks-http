<script lang="ts">
  	import { createEventDispatcher } from 'svelte';
    import { DateInput } from 'date-picker-svelte';

    const dispatch = createEventDispatcher();

    const defaultValue = -1
    let nowDate = new Date()

    let years = [2023, 2022, 2021]
    const months = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December']
    
    let date: Date;
    let day: any;
    let month: any;
    let year: any;

    function yearChange(e) {
        date = null;
        day = null;
        month = -1;

        dateChange();
    }

    function monthChange(e) {
        day = null;
        if (year == -1) {
            year = nowDate.getFullYear()
        }
        
        dateChange();
    }

    function dayChange(e) {
        year = date.getFullYear()
        month = date.getMonth()
        day = date.getDate()

        dateChange();
    }

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
    
</script>

<select name="year" bind:value={year} on:change={yearChange}>
    <option value={defaultValue} disabled selected>Year</option>
    {#each years as year}
		<option value={year}>
			{year}
		</option>
	{/each}
</select>

<select name="month" bind:value={month} on:change={monthChange}>
    <option value={defaultValue} disabled selected>Month</option>
    {#each months as month}
		<option value={months.indexOf(month)}>
			{month}
		</option>
	{/each}
</select>

<DateInput closeOnSelection placeholder="Day" bind:value={date} on:select={dayChange} format="dd" />

