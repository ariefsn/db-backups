<script lang="ts">
	import SearchBar from '../molecules/SearchBar.svelte';
	import * as Select from '$lib/components/ui/select';
	import DateRangePicker from '$lib/components/ui/date-range-picker/DateRangePicker.svelte';
	import { Button } from '$lib/components/ui/button';
	import { XCircle } from '@lucide/svelte';
	import type { DateRange } from 'bits-ui';

	let { 
		search = $bindable(''),
		status = $bindable('all'),
		type = $bindable('all'),
		dateRange = $bindable(),
		onapply,
		onclear
	}: { 
		search: string,
		status: string,
		type: string,
		dateRange?: DateRange,
		onapply: () => void,
		onclear: () => void
	} = $props();

	const statuses = [
		{ value: 'all', label: 'All Statuses' },
		{ value: 'pending', label: 'Pending' },
		{ value: 'generating', label: 'Generating' },
		{ value: 'completed', label: 'Completed' },
		{ value: 'failed', label: 'Failed' }
	];

	const types = [
		{ value: 'all', label: 'All Types' },
		{ value: 'postgre', label: 'PostgreSQL' },
		{ value: 'mysql', label: 'MySQL' },
		{ value: 'mongo', label: 'MongoDB' },
		{ value: 'redis', label: 'Redis' }
	];
</script>

<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
	<SearchBar 
		bind:value={search} 
		onkeydown={(e) => e.key === 'Enter' && onapply()} 
	/>

	<Select.Root type="single" bind:value={status}>
		<Select.Trigger>
			{statuses.find((s) => s.value === status)?.label || 'All Statuses'}
		</Select.Trigger>
		<Select.Content>
			{#each statuses as s}
				<Select.Item value={s.value}>{s.label}</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>

	<Select.Root type="single" bind:value={type}>
		<Select.Trigger>
			{types.find((t) => t.value === type)?.label || 'All Types'}
		</Select.Trigger>
		<Select.Content>
			{#each types as t}
				<Select.Item value={t.value}>{t.label}</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>

	<DateRangePicker bind:value={dateRange} className="w-full" />

	<div class="flex gap-2">
		<Button class="flex-1" onclick={onapply}>Apply</Button>
		<Button variant="outline" size="icon" onclick={onclear} title="Clear Filters">
			<XCircle class="h-4 w-4" />
		</Button>
	</div>
</div>
