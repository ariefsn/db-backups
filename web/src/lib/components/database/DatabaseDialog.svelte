<script lang="ts">
	import type { model_Database } from '$lib/api';
	import { DatabaseService, model_BackupType } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { Loader2 } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	let { 
		open = $bindable(false), 
		database = null, 
		onsaved 
	}: { 
		open: boolean, 
		database: model_Database | null, 
		onsaved?: () => void 
	} = $props();

	let loading = $state(false);
	let useConnectionString = $state(false);

	// Form fields
	let name = $state('');
	let type = $state('');
	let host = $state('');
	let port = $state('');
	let username = $state('');
	let password = $state('');
	let dbName = $state('');
	let connectionUri = $state('');
	let cronExpression = $state('');
	let webhookUrl = $state('');
	let isActive = $state(true);

	const dbTypes = [
		{ value: model_BackupType.Postgres, label: 'PostgreSQL' },
		{ value: model_BackupType.MySQL, label: 'MySQL' },
		{ value: model_BackupType.Mongo, label: 'MongoDB' },
		{ value: model_BackupType.Redis, label: 'Redis' }
	];

	// Initialize form when database changes or dialog opens
	$effect(() => {
		if (open) {
			if (database) {
				name = database.name || '';
				type = database.type || '';
				host = database.host || '';
				port = database.port || '';
				username = database.username || '';
				password = database.password || '';
				dbName = database.database || '';
				connectionUri = database.connectionUri || '';
				cronExpression = database.cronExpression || '';
				webhookUrl = database.webhookUrl || '';
				isActive = database.isActive !== undefined ? database.isActive : true;
				useConnectionString = !!connectionUri;
			} else {
				resetForm();
			}
		}
	});

	function resetForm() {
		name = '';
		type = '';
		host = '';
		port = '';
		username = '';
		password = '';
		dbName = '';
		connectionUri = '';
		cronExpression = '';
		webhookUrl = '';
		isActive = true;
		useConnectionString = false;
	}

	async function handleSubmit() {
		try {
			loading = true;

			// Validation
			if (!name) {
				toast.error('Name is required');
				return;
			}
			if (!type) {
				toast.error('Database type is required');
				return;
			}
			if (!useConnectionString && !host) {
				toast.error('Host is required');
				return;
			}
			if (useConnectionString && !connectionUri) {
				toast.error('Connection URI is required');
				return;
			}

			const payload = {
				name,
				type: type as any,
				host: useConnectionString ? '' : host,
				port: useConnectionString ? '' : port,
				username: useConnectionString ? '' : username,
				password: useConnectionString ? '' : password,
				database: useConnectionString ? '' : dbName,
				connectionUri: useConnectionString ? connectionUri : '',
				cronExpression,
				webhookUrl,
				isActive
			};

			if (database && database.id) {
				await DatabaseService.putDatabases(database.id, payload);
				toast.success('Database updated successfully');
			} else {
				await DatabaseService.postDatabases(payload);
				toast.success('Database created successfully');
			}

			open = false;
			onsaved?.();
		} catch (error: any) {
			console.error('Failed to save database:', error);
			toast.error(error.message || 'Failed to save database');
		} finally {
			loading = false;
		}
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[600px]">
		<Dialog.Header>
			<Dialog.Title>{database ? 'Edit Database' : 'Add Database'}</Dialog.Title>
			<Dialog.Description>
				Configure your database connection and backup schedule.
			</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label>Name</Label>
					<Input bind:value={name} placeholder="Production DB" />
				</div>
				<div class="space-y-2">
					<Label>Type</Label>
					<Select.Root type="single" bind:value={type}>
						<Select.Trigger>
							{dbTypes.find(t => t.value === type)?.label || 'Select Database Type'}
						</Select.Trigger>
						<Select.Content>
							{#each dbTypes as t}
								<Select.Item value={t.value}>{t.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
			</div>

			<div class="flex items-center space-x-2">
				<Label>Connection Method</Label>
				<div class="flex items-center space-x-2 rounded-md border p-1">
					<Label
						class="cursor-pointer rounded-sm px-2 py-1 text-sm {useConnectionString
							? 'bg-primary text-primary-foreground'
							: 'text-muted-foreground'}"
					>
						<input
							type="radio"
							name="connMethod"
							class="hidden"
							bind:group={useConnectionString}
							value={true}
						/>
						URI
					</Label>
					<Label
						class="cursor-pointer rounded-sm px-2 py-1 text-sm {!useConnectionString
							? 'bg-primary text-primary-foreground'
							: 'text-muted-foreground'}"
					>
						<input
							type="radio"
							name="connMethod"
							class="hidden"
							bind:group={useConnectionString}
							value={false}
						/>
						Manual
					</Label>
				</div>
			</div>

			{#if useConnectionString}
				<div class="space-y-2">
					<Label>Connection URI</Label>
					<Input bind:value={connectionUri} placeholder="postgresql://user:pass@host:port/db" />
					<p class="text-xs text-muted-foreground">Enter the full connection string.</p>
				</div>
			{:else}
				<div class="grid grid-cols-2 gap-4">
					<div class="space-y-2">
						<Label>Host</Label>
						<Input bind:value={host} placeholder="localhost" />
					</div>
					<div class="space-y-2">
						<Label>Port</Label>
						<Input bind:value={port} placeholder="5432" />
					</div>
					<div class="space-y-2">
						<Label>Username</Label>
						<Input bind:value={username} placeholder="postgres" />
					</div>
					<div class="space-y-2">
						<Label>Password</Label>
						<Input type="password" bind:value={password} placeholder="••••••••" />
					</div>
					<div class="space-y-2">
						<Label>Database Name</Label>
						<Input bind:value={dbName} placeholder="postgres" />
					</div>
	
				</div>
			{/if}

			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label>Cron Schedule</Label>
					<Input bind:value={cronExpression} placeholder="0 0 * * *" />
					<p class="text-xs text-muted-foreground">Standard cron expression (e.g., 0 0 * * *)</p>
				</div>
				<div class="space-y-2">
					<Label>Webhook URL</Label>
					<Input bind:value={webhookUrl} placeholder="https://api.example.com/webhook" />
				</div>
			</div>

			<div class="flex items-center space-x-2">
				<input type="checkbox" id="isActive" bind:checked={isActive} class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary" />
				<Label for="isActive">Active (Enable automated backups)</Label>
			</div>
		</div>

		<Dialog.Footer>
			<Button variant="outline" onclick={() => (open = false)}>Cancel</Button>
			<Button onclick={handleSubmit} disabled={loading}>
				{#if loading}
					<Loader2 class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				{database ? 'Update' : 'Create'}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
