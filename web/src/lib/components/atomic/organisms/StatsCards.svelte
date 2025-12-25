<script lang="ts">
	import type { database_BackupStats } from '$lib/api';
	import * as Card from '$lib/components/ui/card';
	import { Activity, BarChart3, Database } from '@lucide/svelte';

	let { stats }: { stats: database_BackupStats } = $props();

	function getPercentage(value: number, total: number | undefined) {
		if (!total || total === 0) return 0;
		return Math.round((value / total) * 100);
	}
</script>

<div class="grid gap-4 md:grid-cols-3">
	<Card.Root>
		<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
			<Card.Title class="text-sm font-medium">Total Backups</Card.Title>
			<Database class="h-4 w-4 text-muted-foreground" />
		</Card.Header>
		<Card.Content>
			<div class="text-2xl font-bold">{stats.total || 0}</div>
		</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
			<Card.Title class="text-sm font-medium">Success Rate</Card.Title>
			<Activity class="h-4 w-4 text-muted-foreground" />
		</Card.Header>
		<Card.Content>
			<div class="text-2xl font-bold">
				{getPercentage(stats.byStatus?.completed || 0, stats.total)}%
			</div>
			<p class="text-xs text-muted-foreground">
				{stats.byStatus?.completed || 0} completed
			</p>
		</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
			<Card.Title class="text-sm font-medium">Failed</Card.Title>
			<BarChart3 class="h-4 w-4 text-muted-foreground" />
		</Card.Header>
		<Card.Content>
			<div class="text-2xl font-bold">{stats.byStatus?.failed || 0}</div>
			<p class="text-xs text-muted-foreground">Needs attention</p>
		</Card.Content>
	</Card.Root>
</div>
