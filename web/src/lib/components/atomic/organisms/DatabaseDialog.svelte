<script lang="ts">
	import type { model_Database } from '$lib/api';
	import { DatabaseService, model_BackupType } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { parseConnectionString } from '$lib/db-utils';
	import { Loader2 } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import FormField from '../molecules/FormField.svelte';
	import SegmentedToggle from '../molecules/SegmentedToggle.svelte';

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
	let type = $state<model_BackupType | string>('');
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

	// Handle connection URI parsing
	$effect(() => {
		if (useConnectionString && connectionUri) {
			const parsed = parseConnectionString(connectionUri);
			if (parsed) {
        // parsed.type is not used
				// type = parsed.type;
				host = parsed.host;
				port = parsed.port;
				username = parsed.username;
				password = parsed.password;
				dbName = parsed.database;
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
				host,
				port,
				username,
				password,
				database: dbName,
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
				<FormField label="Name">
					<Input bind:value={name} placeholder="Production DB" />
				</FormField>
				<FormField label="Type">
					<select 
						bind:value={type} 
						class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
					>
						<option value="" disabled>Select Type</option>
						{#each dbTypes as t}
							<option value={t.value}>{t.label}</option>
						{/each}
					</select>
				</FormField>
			</div>

			<SegmentedToggle 
				label="Connection Method"
				bind:value={useConnectionString}
				options={[
					{ value: true, label: 'URI' },
					{ value: false, label: 'Manual' }
				]}
			/>

			{#if useConnectionString}
				<FormField label="Connection URI" description="Enter the full connection string.">
					<Input bind:value={connectionUri} placeholder="postgresql://user:pass@host:port/db" />
				</FormField>
			{:else}
				<div class="grid grid-cols-2 gap-4">
					<FormField label="Host">
						<Input bind:value={host} placeholder="localhost" />
					</FormField>
					<FormField label="Port">
						<Input bind:value={port} placeholder="5432" />
					</FormField>
					<FormField label="Username">
						<Input bind:value={username} placeholder="postgres" />
					</FormField>
					<FormField label="Password">
						<Input type="password" bind:value={password} placeholder="••••••••" />
					</FormField>
					<FormField label="Database Name">
						<Input bind:value={dbName} placeholder="postgres" />
					</FormField>
				</div>
			{/if}

			<div class="grid grid-cols-2 gap-4">
				<FormField label="Cron Schedule" description="Standard cron expression (e.g., 0 0 * * *)">
					<Input bind:value={cronExpression} placeholder="0 0 * * *" />
				</FormField>
				<FormField label="Webhook URL">
					<Input bind:value={webhookUrl} placeholder="https://api.example.com/webhook" />
				</FormField>
			</div>

			<div class="flex items-center space-x-2">
				<input type="checkbox" id="isActive" bind:checked={isActive} class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary" />
				<label for="isActive" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
					Active (Enable automated backups)
				</label>
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
