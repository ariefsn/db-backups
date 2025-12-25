<script lang="ts">
	import * as Card from '$lib/components/ui/card';

	let { 
		title, 
		description, 
		data = {}, 
		total = 0,
		getColor = (key: string) => 'bg-primary'
	}: { 
		title: string, 
		description: string, 
		data: Record<string, number>, 
		total: number,
		getColor?: (key: string) => string
	} = $props();

	function getPercentage(value: number, totalCount: number) {
		if (totalCount === 0) return 0;
		return Math.round((value / totalCount) * 100);
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{title}</Card.Title>
		<Card.Description>{description}</Card.Description>
	</Card.Header>
	<Card.Content class="space-y-4">
		{#each Object.entries(data) as [key, count]}
			<div class="space-y-1">
				<div class="flex items-center justify-between text-sm">
					<span class="font-medium capitalize">{key}</span>
					<span class="text-muted-foreground"
						>{count} ({getPercentage(count, total)}%)</span
					>
				</div>
				<div class="h-2 w-full rounded-full bg-secondary">
					<div
						class="h-2 rounded-full transition-all {getColor(key)}"
						style="width: {getPercentage(count, total)}%"
					></div>
				</div>
			</div>
		{:else}
			<p class="text-muted-foreground text-sm">No data available.</p>
		{/each}
	</Card.Content>
</Card.Root>
