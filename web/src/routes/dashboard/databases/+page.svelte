<script lang="ts">
	import type { model_Database } from '$lib/api';
	import { DatabaseService } from '$lib/api';
	import DatabaseDialog from '$lib/components/database/DatabaseDialog.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import { Loader2, MoreVertical, Pencil, Play, Plus, Search, Trash2 } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	let databases: model_Database[] = [];
	let loading = false;
	let searchQuery = '';
	
	let showDialog = false;
	let selectedDatabase: model_Database | null = null;

	onMount(() => {
		fetchDatabases();
	});

	async function fetchDatabases() {
		try {
			loading = true;
			const response = await DatabaseService.getDatabases();
			databases = response.databases || [];
		} catch (error) {
			console.error('Failed to fetch databases:', error);
			toast.error('Failed to load databases');
		} finally {
			loading = false;
		}
	}

	$: filteredDatabases = databases.filter((db) => {
		if (!searchQuery) return true;
		const query = searchQuery.toLowerCase();
		return (
			db.name?.toLowerCase().includes(query) ||
			db.host?.toLowerCase().includes(query) ||
			db.type?.toLowerCase().includes(query) ||
            db.database?.toLowerCase().includes(query)
		);
	});

	function handleAdd() {
		selectedDatabase = null;
		showDialog = true;
	}

	function handleEdit(db: model_Database) {
		selectedDatabase = db;
		showDialog = true;
	}

	async function handleDelete(db: model_Database) {
        if (!db.id) return;
        if (!confirm('Are you sure you want to delete this database? This action cannot be undone.')) return;
        
		try {
			await DatabaseService.deleteDatabases(db.id);
			toast.success('Database deleted successfully');
			fetchDatabases();
		} catch (error) {
			console.error('Failed to delete database:', error);
			toast.error('Failed to delete database');
		}
	}

	async function handleTriggerBackup(db: model_Database) {
        if (!db.id) return;
		try {
			await DatabaseService.postDatabasesBackup(db.id);
			toast.success('Backup triggered successfully');
		} catch (error) {
			console.error('Failed to trigger backup:', error);
			toast.error('Failed to trigger backup');
		}
	}

	function onSaved() {
		fetchDatabases();
		showDialog = false;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h2 class="text-3xl font-bold tracking-tight">My Databases</h2>
			<p class="text-muted-foreground">Manage your database connections and scheduled backups.</p>
		</div>
		<Button onclick={handleAdd}>
			<Plus class="h-4 w-4" />
			Add Database
		</Button>
	</div>

	<div class="flex items-center space-x-2">
		<div class="relative flex-1 max-w-sm">
			<Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
			<Input
				type="search"
				placeholder="Search databases..."
				class="pl-8"
				bind:value={searchQuery}
			/>
		</div>
	</div>

	<div class="rounded-md border">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Name</Table.Head>
					<Table.Head>Type</Table.Head>
					<Table.Head>Connection</Table.Head>
					<Table.Head>Cron Schedule</Table.Head>
					<Table.Head>Active</Table.Head>
					<Table.Head class="text-right">Actions</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#if loading}
					<Table.Row>
						<Table.Cell colspan={6} class="text-center py-8">
							<div class="flex items-center justify-center">
								<Loader2 class="h-6 w-6 animate-spin text-muted-foreground" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if filteredDatabases.length === 0}
					<Table.Row>
						<Table.Cell colspan={6} class="text-center py-8 text-muted-foreground">
							No databases found. Add one to get started.
						</Table.Cell>
					</Table.Row>
				{:else}
					{#each filteredDatabases as db}
						<Table.Row>
							<Table.Cell class="font-medium">{db.name}</Table.Cell>
							<Table.Cell>
								<Badge variant="outline" class="capitalize">{db.type}</Badge>
							</Table.Cell>
							<Table.Cell>
								{#if db.connectionUri}
									<span class="text-xs text-muted-foreground truncate max-w-[200px] block" title={db.connectionUri}>
										URI
									</span>
								{:else}
									<div class="text-sm">
										<span class="font-medium">{db.host}</span>:{db.port}
                                        <br/>
                                        <span class="text-xs text-muted-foreground">{db.database}</span>
									</div>
								{/if}
							</Table.Cell>
							<Table.Cell>
                                {#if db.cronExpression}
                                    <code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm">
                                        {db.cronExpression}
                                    </code>
                                {:else}
                                    <span class="text-muted-foreground">-</span>
                                {/if}
                            </Table.Cell>
							<Table.Cell>
                                {#if db.isActive}
                                    <Badge>Active</Badge>
                                {:else}
                                    <Badge variant="secondary">Inactive</Badge>
                                {/if}
                            </Table.Cell>
							<Table.Cell class="text-right">
								<DropdownMenu.Root>
									<DropdownMenu.Trigger>
										{#snippet child({ props })}
											<Button variant="ghost" size="icon" {...props}>
												<MoreVertical class="h-4 w-4" />
											</Button>
										{/snippet}
									</DropdownMenu.Trigger>
									<DropdownMenu.Content align="end">
										<DropdownMenu.Item onclick={() => handleTriggerBackup(db)}>
											<Play class="mr-2 h-4 w-4" />
											Run Backup
										</DropdownMenu.Item>
										<DropdownMenu.Item onclick={() => handleEdit(db)}>
											<Pencil class="mr-2 h-4 w-4" />
											Edit
										</DropdownMenu.Item>
										<DropdownMenu.Separator />
										<DropdownMenu.Item class="text-destructive" onclick={() => handleDelete(db)}>
											<Trash2 class="mr-2 h-4 w-4" />
											Delete
										</DropdownMenu.Item>
									</DropdownMenu.Content>
								</DropdownMenu.Root>
							</Table.Cell>
						</Table.Row>
					{/each}
				{/if}
			</Table.Body>
		</Table.Root>
	</div>
</div>

<DatabaseDialog bind:open={showDialog} database={selectedDatabase} onsaved={onSaved} />
