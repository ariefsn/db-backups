<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import DateRangePicker from '$lib/components/ui/date-range-picker/DateRangePicker.svelte';
	import { CalendarDate } from '@internationalized/date';
	import { Activity, BarChart3, Database, XCircle } from '@lucide/svelte';
	import type { DateRange } from 'bits-ui';

	let { data } = $props();

	let dateRange = $state<DateRange | undefined>(undefined);

	// Initialize date range from URL
	const startStr = page.url.searchParams.get('startDate');
	const endStr = page.url.searchParams.get('endDate');
	if (startStr) {
		// startStr is ISO string
		// We need to parse ISO string to CalendarDate for the picker.
		// However, DateRangePicker expects CalendarDate.
		// Let's simplify and rely on the string parsing in DateRangePicker if passing strings,
		// BUT DateRangePicker accepts value: DateRange.
		// Let's manually reconstruct if needed or cleaner:
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

	function getPercentage(value: number, total: number | undefined) {
		if (!total || total === 0) return 0;
		return Math.round((value / total) * 100);
	}
</script>

<svelte:head>
	<title>Reports - DbBackup</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Reports</h1>
			<p class="text-muted-foreground">Detailed statistics about your backups.</p>
		</div>
		<div class="flex items-center gap-2">
			<DateRangePicker bind:value={dateRange} className="w-full" />

			<Button onclick={applyFilters}>Apply</Button>
			<Button variant="outline" size="icon" onclick={clearFilters} title="Clear Filters">
				<XCircle class="h-4 w-4" />
			</Button>
		</div>
	</div>

	<!-- Summary Cards -->
	<div class="grid gap-4 md:grid-cols-3">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Total Backups</Card.Title>
				<Database class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{data.stats.total}</div>
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Success Rate</Card.Title>
				<Activity class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{getPercentage(data.stats.byStatus?.completed || 0, data.stats.total)}%
				</div>
				<p class="text-xs text-muted-foreground">
					{data.stats.byStatus?.completed || 0} completed
				</p>
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Failed</Card.Title>
				<BarChart3 class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{data.stats.byStatus?.failed || 0}</div>
				<p class="text-xs text-muted-foreground">Needs attention</p>
			</Card.Content>
		</Card.Root>
	</div>

	<!-- Charts / Breakdown -->
	<div class="grid gap-4 md:grid-cols-2">
		<Card.Root>
			<Card.Header>
				<Card.Title>Backups by Type</Card.Title>
				<Card.Description>Distribution of backups across database types.</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-4">
				{#each Object.entries(data.stats.byType || {}) as [type, count]}
					<div class="space-y-1">
						<div class="flex items-center justify-between text-sm">
							<span class="font-medium capitalize">{type}</span>
							<span class="text-muted-foreground"
								>{count} ({getPercentage(count, data.stats.total)}%)</span
							>
						</div>
						<div class="h-2 w-full rounded-full bg-secondary">
							<div
								class="h-2 rounded-full bg-primary transition-all"
								style="width: {getPercentage(count, data.stats.total)}%"
							></div>
						</div>
					</div>
				{:else}
					<p class="text-muted-foreground text-sm">No data available.</p>
				{/each}
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header>
				<Card.Title>Backups by Status</Card.Title>
				<Card.Description>Current status of all backups.</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-4">
				{#each Object.entries(data.stats.byStatus || {}) as [status, count]}
					<div class="space-y-1">
						<div class="flex items-center justify-between text-sm">
							<span class="font-medium capitalize">{status}</span>
							<span class="text-muted-foreground"
								>{count} ({getPercentage(count, data.stats.total)}%)</span
							>
						</div>
						<div class="h-2 w-full rounded-full bg-secondary">
							<div
								class="h-2 rounded-full transition-all {status === 'completed'
									? 'bg-green-500'
									: status === 'failed'
										? 'bg-red-500'
										: status === 'generating'
											? 'bg-blue-500'
											: 'bg-gray-500'}"
								style="width: {getPercentage(count, data.stats.total)}%"
							></div>
						</div>
					</div>
				{:else}
					<p class="text-muted-foreground text-sm">No data available.</p>
				{/each}
			</Card.Content>
		</Card.Root>
	</div>
</div>
