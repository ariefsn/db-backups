<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { CheckCircle2, Clock, RefreshCw, ShieldAlert } from '@lucide/svelte';

	let { status = '' }: { status?: string } = $props();
</script>

{#if status === 'completed'}
	<Badge
		class="bg-green-500/15 text-green-700 hover:bg-green-500/25 dark:bg-green-500/10 dark:text-green-400"
	>
		<CheckCircle2 class="mr-1 h-3 w-3" /> Completed
	</Badge>
{:else if status === 'failed'}
	<Badge
		variant="destructive"
		class="bg-red-500/15 text-red-700 hover:bg-red-500/25 dark:bg-red-500/10 dark:text-red-400"
	>
		<ShieldAlert class="mr-1 h-3 w-3" /> Failed
	</Badge>
{:else if status === 'generating'}
	<Badge variant="secondary" class="animate-pulse text-blue-500">
		<RefreshCw class="mr-1 h-3 w-3 animate-spin" /> Generating
	</Badge>
{:else}
	<Badge variant="outline">
		<Clock class="mr-1 h-3 w-3" /> Pending
	</Badge>
{/if}
