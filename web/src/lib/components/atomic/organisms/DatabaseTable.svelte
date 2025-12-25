<script lang="ts">
	import type { model_Database } from '$lib/api';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Table from '$lib/components/ui/table';
	import { Loader2, MoreVertical, Pencil, Play, Trash2 } from '@lucide/svelte';

	let { 
		databases = [], 
		loading = false,
		ontrigger,
		onedit,
		ondelete
	}: { 
		databases: model_Database[], 
		loading?: boolean,
		ontrigger: (db: model_Database) => void,
		onedit: (db: model_Database) => void,
		ondelete: (db: model_Database) => void
	} = $props();
</script>

<div class="rounded-md border bg-card">
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
			{:else if databases.length === 0}
				<Table.Row>
					<Table.Cell colspan={6} class="text-center py-8 text-muted-foreground">
						No databases found. Add one to get started.
					</Table.Cell>
				</Table.Row>
			{:else}
				{#each databases as db}
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
									<DropdownMenu.Item onclick={() => ontrigger(db)}>
										<Play class="mr-2 h-4 w-4" />
										Run Backup
									</DropdownMenu.Item>
									<DropdownMenu.Item onclick={() => onedit(db)}>
										<Pencil class="mr-2 h-4 w-4" />
										Edit
									</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item class="text-destructive" onclick={() => ondelete(db)}>
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
