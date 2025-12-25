<script lang="ts">
	import type { model_Database } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Plus } from '@lucide/svelte';
	import PageHeader from '../molecules/PageHeader.svelte';
	import Pagination from '../molecules/Pagination.svelte';
	import SearchBar from '../molecules/SearchBar.svelte';
	import DatabaseTable from '../organisms/DatabaseTable.svelte';

	let { 
		databases = [], 
		loading = false,
		searchQuery = $bindable(''),
		pagination,
		onadd,
		ontrigger,
		onedit,
		ondelete,
		onpagechange
	}: { 
		databases: model_Database[], 
		loading?: boolean,
		searchQuery: string,
		pagination?: { page: number, limit: number, total: number },
		onadd: () => void,
		ontrigger: (db: model_Database) => void,
		onedit: (db: model_Database) => void,
		ondelete: (db: model_Database) => void,
		onpagechange?: (page: number) => void
	} = $props();
</script>

<PageHeader 
	title="My Databases" 
	description="Manage your database connections and scheduled backups."
>
	{#snippet actions()}
		<Button onclick={onadd}>
			<Plus class="h-4 w-4" />
			Add Database
		</Button>
	{/snippet}
</PageHeader>

<div class="space-y-6">
	<div class="flex items-center space-x-2">
		<div class="max-w-sm flex-1">
			<SearchBar bind:value={searchQuery} placeholder="Search databases..." />
		</div>
	</div>

	<DatabaseTable 
		databases={databases}
		loading={loading}
		ontrigger={ontrigger}
		onedit={onedit}
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
