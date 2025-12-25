<script lang="ts">
	import type { model_BackupMetadata } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { RefreshCw } from '@lucide/svelte';
	import type { DateRange } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import PageHeader from '../molecules/PageHeader.svelte';
	import Pagination from '../molecules/Pagination.svelte';
	import BackupTable from '../organisms/BackupTable.svelte';
	import FilterBar from '../organisms/FilterBar.svelte';

	let { 
		backups = [], 
		search = $bindable(''),
		status = $bindable('all'),
		type = $bindable('all'),
		dateRange = $bindable(),
		pagination,
		onrefresh,
		onapply,
		onclear,
		onview,
		ondownload,
		ondelete,
		onpagechange,
		createBackupTrigger
	}: { 
		backups: model_BackupMetadata[], 
		search: string,
		status: string,
		type: string,
		dateRange?: DateRange,
		pagination?: { page: number, limit: number, total: number },
		onrefresh: () => void,
		onapply: () => void,
		onclear: () => void,
		onview: (id: string) => void,
		ondownload: (id: string) => void,
		ondelete: (id: string) => void,
		onpagechange?: (page: number) => void,
		createBackupTrigger: Snippet
	} = $props();
</script>

<PageHeader 
	title="Backups" 
	description="Manage and monitor your database backups."
>
	{#snippet actions()}
		{@render createBackupTrigger()}
		<Button variant="outline" onclick={onrefresh}>
			<RefreshCw class="mr-2 h-4 w-4" />
			Refresh
		</Button>
	{/snippet}
</PageHeader>

<div class="space-y-6">
	<div class="rounded-md border p-4 bg-card">
		<FilterBar 
			bind:search={search}
			bind:status={status}
			bind:type={type}
			bind:dateRange={dateRange}
			onapply={onapply}
			onclear={onclear}
		/>
	</div>

	<BackupTable 
		backups={backups}
		onview={onview}
		ondownload={ondownload}
		ondelete={ondelete}
	/>

	{#if pagination && onpagechange}
		<Pagination 
			page={pagination.page}
			limit={pagination.limit}
			total={pagination.total}
			onPageChange={onpagechange}
		/>
	{/if}
</div>
