<script lang="ts">
  	import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    let nowDate = new Date()

    let days = []
    let years = []
    const monthLengths = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]
    const months = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December']

    for (let x=2022; x <= nowDate.getFullYear(); x+=1){
        years.push(x)
    }
    
    let day;
    let month;
    let year;

    function isLeapYear(year: number): boolean {
        return (year % 4 === 0 && year % 100 !== 0) || year % 400 === 0
    }

    function monthChange(e) {
        // TODO
        const feb = isLeapYear(year) ? 29 : 28
        
        days = []
        for (let x = 1; x <= monthLengths[months.indexOf(month)]; x += 1) {
            days.push(x)
        }

        dateChange();

        console.log(day, month, year)
    }

    function dayChange(e) {
        dateChange();
        console.log(day, month, year)
    }

    function yearChange(e) {
        dateChange();
        console.log(day, month, year)
    }

    function dateChange() {
        let start: Date;
        let end: Date;
        if (year && month && day) {
            start = new Date(year, months.indexOf(month), day, 0)
            end = new Date(year, months.indexOf(month), day, 23, 59, 59)
        } else if (year && month) {
            start = new Date(year, months.indexOf(month), 0)
            end = new Date(year, months.indexOf(month), monthLengths[months.indexOf(month)])
        } else if (year) {
            start = new Date(year, 0, 1)
            end = new Date(year, 12, 31)
        }

        dispatch('dateChange', {
			start: start.toISOString(),
            end: end.toISOString()
		});
    }
    
</script>

<select name="year" bind:value={year} on:change={yearChange}>
    <option value="" disabled selected hidden>Year</option>
    {#each years as year}
		<option value={year}>
			{year}
		</option>
	{/each}
</select>
<select name="month" bind:value={month} on:change={monthChange}>
    <option value="" disabled selected hidden>Month</option>
    {#each months as month}
		<option value={month}>
			{month}
		</option>
	{/each}
</select>
<select name="day" bind:value={day} on:change={dayChange}>
    <option value="" disabled selected hidden>Day</option>
    {#each days as day}
		<option value={day}>
			{day}
		</option>
	{/each}
</select>

