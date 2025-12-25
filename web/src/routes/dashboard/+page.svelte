<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { page } from '$app/state';
	import { BackupService } from '$lib/api';
	import BackupDetailsDialog from '$lib/components/backup/BackupDetailsDialog.svelte';
	import CreateBackupDialog from '$lib/components/backup/CreateBackupDialog.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import DateRangePicker from '$lib/components/ui/date-range-picker/DateRangePicker.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import { CalendarDate } from '@internationalized/date';
	import {
		CheckCircle2,
		Clock,
		Database,
		Download,
		MoreVertical,
		RefreshCw,
		Search,
		ShieldAlert,
		Trash2,
		XCircle
	} from '@lucide/svelte';
	import type { DateRange } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	let { data } = $props();

	// Filters state
	let search = $state(page.url.searchParams.get('search') || '');
	let status = $state(page.url.searchParams.get('statuses') || 'all');
	let type = $state(page.url.searchParams.get('types') || 'all');
	
    // Date range filter
    let dateRange = $state<DateRange | undefined>(undefined);
    
    // Initialize date range from URL
    const startStr = page.url.searchParams.get('startDate');
    const endStr = page.url.searchParams.get('endDate');
    if (startStr) { // startStr is ISO string
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
         }
    }
    
	
	let isDeleting = $state(false);
	let deleteId = $state<string | null>(null);
    
    // Details Dialog
    let showDetails = $state(false);
    let selectedBackupId = $state<string | null>(null);

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

	function applyFilters() {
		const params = new URLSearchParams();
		if (search) params.set('search', search);
		if (status !== 'all') params.set('statuses', status);
		if (type !== 'all') params.set('types', type);
		
        if (dateRange?.start) {
            params.set('startDate', dateRange.start.toString());
        }
        if (dateRange?.end) {
            params.set('endDate', dateRange.end.toString());
        }
		
		goto(`?${params.toString()}`);
	}

	function clearFilters() {
		search = '';
		status = 'all';
		type = 'all';
		dateRange = undefined;
		goto('/dashboard');
	}

	function formatBytes(bytes: number) {
		if (bytes === 0) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	function formatDate(dateStr: string) {
		return new Date(dateStr).toLocaleString();
	}

	async function handleDelete() {
		if (!deleteId) return;
		isDeleting = true;
		try {
			await BackupService.deleteBackups(deleteId);
			toast.success('Backup deleted successfully');
			deleteId = null;
			invalidate('app:backups');
		} catch (error) {
			toast.error('Failed to delete backup');
			console.error(error);
		} finally {
			isDeleting = false;
		}
	}

	async function handleDownload(id: string) {
		try {
			const { url } = await BackupService.getBackupsDownload(id);
			window.open(url, '_blank');
		} catch (error) {
			toast.error('Failed to get download URL');
			console.error(error);
		}
	}
    
    function viewDetails(backupId: string) {
        selectedBackupId = backupId;
        showDetails = true;
    }
</script>

<svelte:head>
	<title>Backups - DbBackup</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Backups</h1>
			<p class="text-muted-foreground">Manage and monitor your database backups.</p>
		</div>
		<div class="flex gap-2">
            <CreateBackupDialog />
			<Button variant="outline" onclick={() => invalidate('app:backups')}>
				<RefreshCw class="mr-2 h-4 w-4" />
				Refresh
			</Button>
		</div>
	</div>

	<!-- Filters -->
	<Card.Root>
		<Card.Content class="p-4">
			<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
				<div class="relative">
					<Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
					<Input 
						type="search" 
						placeholder="Search..." 
						class="pl-8" 
						bind:value={search}
						onkeydown={(e: KeyboardEvent) => e.key === 'Enter' && applyFilters()}
					/>
				</div>
				
				<Select.Root type="single" bind:value={status}>
					<Select.Trigger>
						{statuses.find(s => s.value === status)?.label || 'All Statuses'}
					</Select.Trigger>
					<Select.Content>
						{#each statuses as s}
							<Select.Item value={s.value}>{s.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

				<Select.Root type="single" bind:value={type}>
					<Select.Trigger>
						{types.find(t => t.value === type)?.label || 'All Types'}
					</Select.Trigger>
					<Select.Content>
						{#each types as t}
							<Select.Item value={t.value}>{t.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

				<DateRangePicker bind:value={dateRange} className="w-full" />

				<div class="flex gap-2">
					<Button class="flex-1" onclick={applyFilters}>Apply</Button>
					<Button variant="outline" size="icon" onclick={clearFilters} title="Clear Filters">
						<XCircle class="h-4 w-4" />
					</Button>
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<!-- Table -->
	<Card.Root class="p-4">
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
				{#if data.backups.length === 0}
					<Table.Row>
						<Table.Cell colspan={7} class="text-center h-24 text-muted-foreground">
							No backups found.
						</Table.Cell>
					</Table.Row>
				{:else}
					{#each data.backups as backup (backup.id)}
						<Table.Row>
							<Table.Cell>
								{#if backup.status === 'completed'}
									<Badge class="bg-green-500/15 text-green-700 hover:bg-green-500/25 dark:bg-green-500/10 dark:text-green-400">
										<CheckCircle2 class="mr-1 h-3 w-3" /> Completed
									</Badge>
								{:else if backup.status === 'failed'}
									<Badge variant="destructive" class="bg-red-500/15 text-red-700 hover:bg-red-500/25 dark:bg-red-500/10 dark:text-red-400">
										<ShieldAlert class="mr-1 h-3 w-3" /> Failed
									</Badge>
								{:else if backup.status === 'generating'}
									<Badge variant="secondary" class="animate-pulse text-blue-500">
										<RefreshCw class="mr-1 h-3 w-3 animate-spin" /> Generating
									</Badge>
								{:else}
									<Badge variant="outline">
										<Clock class="mr-1 h-3 w-3" /> Pending
									</Badge>
								{/if}
							</Table.Cell>
							<Table.Cell>
								<div class="flex items-center gap-2 font-medium">
									<Database class="h-4 w-4 text-muted-foreground" />
									{backup.type?.toUpperCase() || '-'}
								</div>
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
										{backup.error.length > 50 ? backup.error.substring(0, 50) + '...' : backup.error}
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
                                        <DropdownMenu.Item onclick={() => backup.id && viewDetails(backup.id)}>
											<Database class="mr-2 h-4 w-4" /> Details
										</DropdownMenu.Item>
										<DropdownMenu.Item onclick={() => backup.id && handleDownload(backup.id)} disabled={backup.status !== 'completed'}>
											<Download class="mr-2 h-4 w-4" /> Download
										</DropdownMenu.Item>
										<DropdownMenu.Item class="text-destructive focus:text-destructive" onclick={() => deleteId = backup.id || null}>
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
	</Card.Root>
</div>

<!-- Delete Confirmation -->
<Dialog.Root open={!!deleteId} onOpenChange={(open) => !open && (deleteId = null)}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Are you sure?</Dialog.Title>
			<Dialog.Description>
				This action cannot be undone. This will permanently delete the backup from the database and storage.
			</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => deleteId = null} disabled={isDeleting}>Cancel</Button>
			<Button variant="destructive" onclick={handleDelete} disabled={isDeleting}>
				{#if isDeleting}
					<RefreshCw class="mr-2 h-4 w-4 animate-spin" /> Deleting...
				{:else}
					Delete
				{/if}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<BackupDetailsDialog bind:open={showDetails} backupId={selectedBackupId} />
