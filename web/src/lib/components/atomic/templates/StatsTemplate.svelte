<script lang="ts">
	import PageHeader from '../molecules/PageHeader.svelte';
	import StatsCards from '../organisms/StatsCards.svelte';
	import BreakdownCard from '../organisms/BreakdownCard.svelte';
	import DateRangePicker from '$lib/components/ui/date-range-picker/DateRangePicker.svelte';
	import { Button } from '$lib/components/ui/button';
	import { XCircle } from '@lucide/svelte';
	import type { DateRange } from 'bits-ui';
	import type { database_BackupStats } from '$lib/api';

	let { 
		stats, 
		dateRange = $bindable(),
		onapply,
		onclear
	}: { 
		stats: database_BackupStats,
		dateRange?: DateRange,
		onapply: () => void,
		onclear: () => void
	} = $props();

	function getStatusColor(status: string) {
		switch (status) {
			case 'completed': return 'bg-green-500';
			case 'failed': return 'bg-red-500';
			case 'generating': return 'bg-blue-500';
			default: return 'bg-gray-500';
		}
	}
</script>

<PageHeader 
	title="Reports" 
	description="Detailed statistics about your backups."
>
	{#snippet actions()}
		<div class="flex items-center gap-2">
			<DateRangePicker bind:value={dateRange} className="w-full" />
			<Button onclick={onapply}>Apply</Button>
			<Button variant="outline" size="icon" onclick={onclear} title="Clear Filters">
				<XCircle class="h-4 w-4" />
			</Button>
		</div>
	{/snippet}
</PageHeader>

<div class="space-y-6">
	<StatsCards stats={stats} />

	<div class="grid gap-4 md:grid-cols-2">
		<BreakdownCard 
			title="Backups by Type"
			description="Distribution of backups across database types."
			data={stats.byType || {}}
			total={stats.total || 0}
		/>

		<BreakdownCard 
			title="Backups by Status"
			description="Current status of all backups."
			data={stats.byStatus || {}}
			total={stats.total || 0}
			getColor={getStatusColor}
		/>
	</div>
</div>
