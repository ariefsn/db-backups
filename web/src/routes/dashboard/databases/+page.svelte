<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { page } from '$app/state';
	import type { model_Database } from '$lib/api';
	import { DatabaseService } from '$lib/api';
	import DatabaseDialog from '$lib/components/atomic/organisms/DatabaseDialog.svelte';
	import DatabaseTemplate from '$lib/components/atomic/templates/DatabaseTemplate.svelte';
	import { toast } from 'svelte-sonner';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let searchQuery = $state('');
	let showDialog = $state(false);
	let selectedDatabase = $state<model_Database | null>(null);
	let loading = $state(false);

	const filteredDatabases = $derived(
		(data?.databases || []).filter((db: model_Database) =>
			db.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
			db.host?.toLowerCase().includes(searchQuery.toLowerCase()) ||
			db.database?.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	function handleAdd() {
		selectedDatabase = null;
		showDialog = true;
	}

	function handleEdit(database: model_Database) {
		selectedDatabase = database;
		showDialog = true;
	}

	async function handleTrigger(database: model_Database) {
		if (!database.id) return;
		try {
			loading = true;
			await DatabaseService.postDatabasesBackup(database.id);
			toast.success(`Backup triggered for ${database.name}`);
		} catch (error: any) {
			toast.error(error.message || 'Failed to trigger backup');
		} finally {
			loading = false;
		}
	}

	async function handleDelete(database: model_Database) {
		if (!database.id || !confirm(`Are you sure you want to delete ${database.name}?`)) return;
		try {
			await DatabaseService.deleteDatabases(database.id);
			toast.success('Database deleted successfully');
			invalidate('app:databases');
		} catch (error: any) {
			toast.error(error.message || 'Failed to delete database');
		}
	}

	function onSaved() {
		invalidate('app:databases');
	}

	function handlePageChange(newPage: number) {
		const params = new URLSearchParams(page.url.searchParams);
		params.set('page', newPage.toString());
		goto(`?${params.toString()}`);
	}
</script>

<svelte:head>
	<title>My Databases - DbBackup</title>
</svelte:head>

<DatabaseTemplate 
	databases={filteredDatabases}
	loading={loading}
	bind:searchQuery={searchQuery}
	pagination={{
		page: data.pagination?.page ?? 1,
		limit: data.pagination?.limit ?? 10,
		total: data.pagination?.total ?? 0
	}}
	onadd={handleAdd}
	ontrigger={handleTrigger}
	onedit={handleEdit}
	ondelete={handleDelete}
	onpagechange={handlePageChange}
/>

<DatabaseDialog bind:open={showDialog} database={selectedDatabase} onsaved={onSaved} />
