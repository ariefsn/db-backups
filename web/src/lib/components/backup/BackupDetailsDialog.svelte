<script lang="ts">
	import type { model_BackupMetadata } from '$lib/api';
	import { BackupService } from '$lib/api';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Separator } from '$lib/components/ui/separator';
	import { AlertCircle, CheckCircle2, Clock, Database, RefreshCw, ShieldAlert } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	let { open = $bindable(false), backupId } = $props<{ 
		open: boolean; 
		backupId: string | null 
	}>();

	let backup = $state<model_BackupMetadata | null>(null);
	let loading = $state(false);
	let error = $state<string | null>(null);

	// Fetch backup details when dialog opens
	$effect(() => {
		if (open && backupId) {
			fetchBackupDetails();
		} else if (!open) {
			// Reset state when dialog closes
			backup = null;
			error = null;
		}
	});

	async function fetchBackupDetails() {
		if (!backupId) return;
		
		loading = true;
		error = null;
		
		try {
			backup = await BackupService.getBackupById(backupId);
		} catch (err) {
			error = 'Failed to load backup details';
			toast.error('Failed to load backup details');
			console.error(err);
		} finally {
			loading = false;
		}
	}

	function formatBytes(bytes: number) {
		if (bytes === 0) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	function formatDate(dateStr?: string) {
		if (!dateStr) return '-';
		return new Date(dateStr).toLocaleString();
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[600px]">
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<div class="text-center space-y-4">
					<RefreshCw class="h-8 w-8 animate-spin mx-auto text-muted-foreground" />
					<p class="text-sm text-muted-foreground">Loading backup details...</p>
				</div>
			</div>
		{:else if error}
			<div class="flex items-center justify-center py-12">
				<div class="text-center space-y-4">
					<AlertCircle class="h-8 w-8 mx-auto text-destructive" />
					<p class="text-sm text-destructive">{error}</p>
					<Button variant="outline" onclick={() => open = false}>Close</Button>
				</div>
			</div>
		{:else if backup}
			<Dialog.Header>
				<Dialog.Title class="flex items-center gap-2">
					<Database class="h-5 w-5" />
					Backup Details
				</Dialog.Title>
				<Dialog.Description>
					ID: <span class="font-mono text-xs">{backup.id}</span>
				</Dialog.Description>
			</Dialog.Header>

			<div class="grid gap-6 py-4">
				<div class="flex items-center justify-between rounded-lg border p-4 bg-muted/40">
					<div class="space-y-1">
						<p class="text-sm font-medium text-muted-foreground">Status</p>
						<div class="flex items-center">
							{#if backup.status === 'completed'}
								<Badge class="bg-green-500/15 text-green-700 dark:bg-green-500/10 dark:text-green-400">
									<CheckCircle2 class="mr-1 h-3 w-3" /> Completed
								</Badge>
							{:else if backup.status === 'failed'}
								<Badge variant="destructive" class="bg-red-500/15 text-red-700 dark:bg-red-500/10 dark:text-red-400">
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
						</div>
					</div>
					<div class="space-y-1 text-right">
						<p class="text-sm font-medium text-muted-foreground">Size</p>
						<p class="text-lg font-semibold">{formatBytes(backup.fileSize || 0)}</p>
					</div>
				</div>

				<div class="space-y-4">
					<h4 class="font-medium leading-none">Database Information</h4>
					<div class="grid grid-cols-2 gap-4 text-sm">
						<div>
							<p class="text-muted-foreground mb-1">Type</p>
							<p class="font-medium uppercase">{backup.type}</p>
						</div>
						<div>
							<p class="text-muted-foreground mb-1">Database Name</p>
							<p class="font-medium">{backup.database}</p>
						</div>
						<div>
							<p class="text-muted-foreground mb-1">Host</p>
							<p class="font-medium">{backup.host || '-'}</p>
						</div>
						<div>
							<p class="text-muted-foreground mb-1">Created At</p>
							<p class="font-medium">{formatDate(backup.createdAt)}</p>
						</div>
					</div>
				</div>

				<Separator />

				<div class="space-y-4">
					<h4 class="font-medium leading-none">Storage Information</h4>
					<div class="grid gap-2 text-sm">
						<div>
							<p class="text-muted-foreground mb-1">Object Key</p>
							<code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm">
								{backup.objectKey}
							</code>
						</div>
						{#if backup.filePath}
							<div>
								<p class="text-muted-foreground mb-1">Local Path</p>
								<code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm">
									{backup.filePath}
								</code>
							</div>
						{/if}
					</div>
				</div>
                
                {#if backup.error}
                    <Separator />
                    <div class="rounded-md bg-destructive/10 p-4">
                        <div class="flex items-center gap-2 text-destructive mb-2">
                            <ShieldAlert class="h-4 w-4" />
                            <h4 class="font-medium">Error Message</h4>
                        </div>
                        <p class="text-sm text-destructive/90 font-mono whitespace-pre-wrap">{backup.error}</p>
                    </div>
                {/if}
			</div>

			<Dialog.Footer>
				<Button variant="outline" onclick={() => open = false}>Close</Button>
			</Dialog.Footer>
		{/if}
	</Dialog.Content>
</Dialog.Root>
