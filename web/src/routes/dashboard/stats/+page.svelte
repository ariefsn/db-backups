<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import StatsTemplate from '$lib/components/atomic/templates/StatsTemplate.svelte';
	import { CalendarDate } from '@internationalized/date';
	import type { DateRange } from 'bits-ui';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let dateRange = $state<DateRange | undefined>(undefined);

	// Initialize date range from URL
	const startStr = page.url.searchParams.get('startDate');
	const endStr = page.url.searchParams.get('endDate');
	if (startStr) {
		const s = new Date(startStr);
		const e = endStr ? new Date(endStr) : undefined;
		dateRange = {
			start: new CalendarDate(s.getFullYear(), s.getMonth() + 1, s.getDate()),
			end: e ? new CalendarDate(e.getFullYear(), e.getMonth() + 1, e.getDate()) : undefined
		};
	}

	function applyFilters() {
		const params = new URLSearchParams();
		if (dateRange?.start) params.set('startDate', dateRange.start.toString());
		if (dateRange?.end) params.set('endDate', dateRange.end.toString());
		goto(`?${params.toString()}`);
	}

	function clearFilters() {
		dateRange = undefined;
		goto('/dashboard/stats');
	}
</script>

<svelte:head>
	<title>Reports - DbBackup</title>
</svelte:head>

<StatsTemplate 
	stats={data.stats}
	bind:dateRange={dateRange}
	onapply={applyFilters}
	onclear={clearFilters}
/>
