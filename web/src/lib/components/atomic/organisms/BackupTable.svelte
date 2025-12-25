<script lang="ts">
	import type { model_BackupMetadata } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Table from '$lib/components/ui/table';
	import { Database, Download, MoreVertical, Trash2 } from '@lucide/svelte';
	import DbTypeBadge from '../molecules/DbTypeBadge.svelte';
	import StatusBadge from '../molecules/StatusBadge.svelte';

	let { 
		backups = [], 
		onview, 
		ondownload, 
		ondelete 
	}: { 
		backups: model_BackupMetadata[], 
		onview: (id: string) => void,
		ondownload: (id: string) => void,
		ondelete: (id: string) => void
	} = $props();

	function formatBytes(bytes: number = 0) {
		if (bytes === 0) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	function formatDate(dateStr: string) {
		return new Date(dateStr).toLocaleString();
	}
</script>

<div class="rounded-md border p-4 bg-card">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Status</Table.Head>
				<Table.Head>Type</Table.Head>
				<Table.Head>Database / Host</Table.Head>
				<Table.Head>Size</Table.Head>
				<Table.Head>Created At</Table.Head>
				<Table.Head>Error</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#if backups.length === 0}
				<Table.Row>
					<Table.Cell colspan={7} class="h-24 text-center text-muted-foreground">
						No backups found.
					</Table.Cell>
				</Table.Row>
			{:else}
				{#each backups as backup (backup.id)}
					<Table.Row>
						<Table.Cell>
							<StatusBadge status={backup.status} />
						</Table.Cell>
						<Table.Cell>
							<DbTypeBadge type={backup.type} />
						</Table.Cell>
						<Table.Cell>
							<div class="flex flex-col">
								<span class="font-medium">{backup.database}</span>
								<span class="text-xs text-muted-foreground">{backup.host}</span>
							</div>
						</Table.Cell>
						<Table.Cell>{formatBytes(backup.fileSize || 0)}</Table.Cell>
						<Table.Cell>{backup.createdAt ? formatDate(backup.createdAt) : '-'}</Table.Cell>
						<Table.Cell>
							{#if backup.error}
								<span class="text-xs text-destructive" title={backup.error}>
									{backup.error.length > 50
										? backup.error.substring(0, 50) + '...'
										: backup.error}
								</span>
							{:else}
								<span class="text-muted-foreground">-</span>
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
									<DropdownMenu.Label>Actions</DropdownMenu.Label>
									<DropdownMenu.Separator />
									<DropdownMenu.Item onclick={() => backup.id && onview(backup.id)}>
										<Database class="mr-2 h-4 w-4" /> Details
									</DropdownMenu.Item>
									<DropdownMenu.Item
										onclick={() => backup.id && ondownload(backup.id)}
										disabled={backup.status !== 'completed'}
									>
										<Download class="mr-2 h-4 w-4" /> Download
									</DropdownMenu.Item>
									<DropdownMenu.Item
										class="text-destructive focus:text-destructive"
										onclick={() => backup.id && ondelete(backup.id)}
									>
										<Trash2 class="mr-2 h-4 w-4" /> Delete
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
