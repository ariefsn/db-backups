<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { BackupService, model_BackupType } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { Plus, RefreshCw } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	let open = $state(false);
	let loading = $state(false);
	let type = $state<model_BackupType>(model_BackupType.Postgres);
	let useConnectionString = $state(true);

	// Form state
	let connectionString = $state('');
	let host = $state('');
	let port = $state('');
	let username = $state('');
	let password = $state('');
	let database = $state('');
	let webhookUrl = $state('');

	// Polling state
	let pollingInterval: ReturnType<typeof setInterval> | null = null;

	const types = [
		{ value: model_BackupType.Postgres, label: 'PostgreSQL' },
		{ value: model_BackupType.MySQL, label: 'MySQL' },
		{ value: model_BackupType.Mongo, label: 'MongoDB' },
		{ value: model_BackupType.Redis, label: 'Redis' }
	];

	async function handleSubmit(e: Event) {
		e.preventDefault();
		loading = true;

		try {
			const payload: any = {
				type,
				webhookUrl: webhookUrl || undefined
			};

			if (useConnectionString) {
				payload.connectionUri = connectionString;
			} else {
				payload.host = host;
				payload.port = port;
				payload.username = username;
				payload.password = password;
				payload.database = database;
			}

			const response = await BackupService.postBackup(payload);

			toast.success('Backup triggered successfully');
			open = false;
			resetForm();
			invalidate('app:backups');

			// Start polling if we got a backup ID
			if (response.id) {
				startPolling(response.id);
			}
		} catch (error) {
			toast.error('Failed to trigger backup');
			console.error(error);
		} finally {
			loading = false;
		}
	}

	function startPolling(backupId: string) {
		// Clear any existing interval
		if (pollingInterval) {
			clearInterval(pollingInterval);
		}

		// Poll every 5 seconds
		pollingInterval = setInterval(async () => {
			try {
				const backup = await BackupService.getBackupById(backupId);

				// Update the backup list
				invalidate('app:backups');

				// Stop polling if backup reached final state
				if (backup.status === 'completed' || backup.status === 'failed') {
					if (pollingInterval) {
						clearInterval(pollingInterval);
						pollingInterval = null;
					}
				}
			} catch (error) {
				console.error('Failed to poll backup status:', error);
				// Stop polling on error
				if (pollingInterval) {
					clearInterval(pollingInterval);
					pollingInterval = null;
				}
			}
		}, 5000);
	}

	function resetForm() {
		connectionString = '';
		host = '';
		port = '';
		username = '';
		password = '';
		database = '';
		webhookUrl = '';
		type = model_BackupType.Postgres;
		useConnectionString = true;
	}

	// Cleanup on component destroy
	$effect(() => {
		return () => {
			if (pollingInterval) {
				clearInterval(pollingInterval);
			}
		};
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		<Button>
			<Plus class="h-4 w-4" />
			Create Backup
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[600px]">
		<Dialog.Header>
			<Dialog.Title>Create New Backup</Dialog.Title>
			<Dialog.Description>
				Enter the database connection details to trigger a new backup.
			</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="grid gap-4 py-4">
			<!-- Type Selection -->
			<div class="space-y-2">
				<Label>Database Type</Label>
				<Select.Root type="single" bind:value={type}>
					<Select.Trigger>
						{types.find((t) => t.value === type)?.label}
					</Select.Trigger>
					<Select.Content>
						{#each types as t}
							<Select.Item value={t.value}>{t.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>

			<!-- Connection Method Toggle -->
			<div class="flex items-center gap-4 rounded-lg border bg-muted/40 p-3">
				<Label class="flex-1 cursor-pointer">
					<input
						type="radio"
						name="connectionMethod"
						checked={useConnectionString}
						onchange={() => (useConnectionString = true)}
						class="mr-2"
					/>
					Connection String
				</Label>
				<Label class="flex-1 cursor-pointer">
					<input
						type="radio"
						name="connectionMethod"
						checked={!useConnectionString}
						onchange={() => (useConnectionString = false)}
						class="mr-2"
					/>
					Traditional
				</Label>
			</div>

			{#if useConnectionString}
				<!-- Connection String Input -->
				<div class="space-y-2">
					<Label>Connection String</Label>
					<Input
						bind:value={connectionString}
						placeholder="postgresql://user:password@localhost:5432/dbname"
						required
					/>
					<p class="text-xs text-muted-foreground">Example: postgresql://user:pass@host:port/db</p>
				</div>
			{:else}
				<!-- Traditional Connection Fields -->
				<div class="grid grid-cols-2 gap-4">
					<div class="space-y-2">
						<Label>Database Name</Label>
						<Input bind:value={database} placeholder="mydb" required />
					</div>
					<div class="space-y-2">
						<Label>Host</Label>
						<Input bind:value={host} placeholder="localhost" required />
					</div>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div class="space-y-2">
						<Label>Port</Label>
						<Input bind:value={port} placeholder="5432" required />
					</div>
					<div class="space-y-2">
						<Label>Username</Label>
						<Input bind:value={username} placeholder="user" />
					</div>
				</div>

				<div class="space-y-2">
					<Label>Password</Label>
					<Input type="password" bind:value={password} placeholder="password" />
				</div>
			{/if}

			<!-- Webhook URL (Optional) -->
			<div class="space-y-2">
				<Label>Webhook URL (Optional)</Label>
				<Input bind:value={webhookUrl} placeholder="https://api.example.com/webhook" />
			</div>

			<Dialog.Footer>
				<Button type="submit" disabled={loading}>
					{#if loading}
						<RefreshCw class="mr-2 h-4 w-4 animate-spin" /> Starting...
					{:else}
						Start Backup
					{/if}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
