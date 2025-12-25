<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { page } from '$app/state';
	import { BackupService } from '$lib/api';
	import CreateBackupDialog from '$lib/components/atomic/organisms/CreateBackupDialog.svelte';
	import DashboardTemplate from '$lib/components/atomic/templates/DashboardTemplate.svelte';
	import BackupDetailsDialog from '$lib/components/backup/BackupDetailsDialog.svelte';
	import { Button } from '$lib/components/ui/button';
	import { CalendarDate } from '@internationalized/date';
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
	if (startStr) {
		const s = new Date(startStr);
		const e = endStr ? new Date(endStr) : undefined;
		dateRange = {
			start: new CalendarDate(s.getFullYear(), s.getMonth() + 1, s.getDate()),
			end: e ? new CalendarDate(e.getFullYear(), e.getMonth() + 1, e.getDate()) : undefined
		};
	}

	let deleteId = $state<string | null>(null);
	let showDetails = $state(false);
	let selectedBackupId = $state<string | null>(null);

	function applyFilters() {
		const params = new URLSearchParams();
		if (search) params.set('search', search);
		if (status !== 'all') params.set('statuses', status);
		if (type !== 'all') params.set('types', type);
		if (dateRange?.start) params.set('startDate', dateRange.start.toString());
		if (dateRange?.end) params.set('endDate', dateRange.end.toString());
		goto(`?${params.toString()}`);
	}

	function clearFilters() {
		search = '';
		status = 'all';
		type = 'all';
		dateRange = undefined;
		goto('/dashboard');
	}

	async function handleDelete() {
		if (!deleteId) return;
		try {
			await BackupService.deleteBackups(deleteId);
			toast.success('Backup deleted successfully');
			deleteId = null;
			invalidate('app:backups');
		} catch (error) {
			toast.error('Failed to delete backup');
			console.error(error);
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
	function handlePageChange(newPage: number) {
		const params = new URLSearchParams(page.url.searchParams);
		params.set('page', newPage.toString());
		goto(`?${params.toString()}`);
	}
</script>

<svelte:head>
	<title>Backups - DbBackup</title>
</svelte:head>

<DashboardTemplate 
	backups={data.backups}
	bind:search={search}
	bind:status={status}
	bind:type={type}
	bind:dateRange={dateRange}
	pagination={{
		page: data.pagination?.page ?? 1,
		limit: data.pagination?.limit ?? 10,
		total: data.pagination?.total ?? 0
	}}
	onrefresh={() => invalidate('app:backups')}
	onapply={applyFilters}
	onclear={clearFilters}
	onview={viewDetails}
	ondownload={handleDownload}
	ondelete={(id) => (deleteId = id)}
	onpagechange={handlePageChange}
>
	{#snippet createBackupTrigger()}
		<CreateBackupDialog />
	{/snippet}
</DashboardTemplate>

<BackupDetailsDialog bind:open={showDetails} backupId={selectedBackupId} />

{#if deleteId}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
		<div class="w-full max-w-md bg-card p-6 rounded-lg shadow-lg border">
			<h2 class="text-xl font-bold mb-4">Confirm Delete</h2>
			<p class="text-muted-foreground mb-6">Are you sure you want to delete this backup? This action cannot be undone.</p>
			<div class="flex justify-end gap-2">
				<Button variant="outline" onclick={() => (deleteId = null)}>Cancel</Button>
				<Button variant="destructive" onclick={handleDelete}>Delete</Button>
			</div>
		</div>
	</div>
{/if}
