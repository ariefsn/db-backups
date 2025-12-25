<script lang="ts">
	import { Label } from '$lib/components/ui/label';

	let { 
		value = $bindable(), 
		options = [],
		label = ''
	}: { 
		value: any, 
		options: { value: any, label: string }[],
		label?: string
	} = $props();
</script>

<div class="flex items-center space-x-2">
	{#if label}
		<Label>{label}</Label>
	{/if}
	<div class="flex items-center space-x-2 rounded-md border p-1 bg-muted/20">
		{#each options as option}
			<Label
				class="cursor-pointer rounded-sm px-2 py-1 text-sm {value === option.value
					? 'bg-primary text-primary-foreground shadow-sm'
					: 'text-muted-foreground hover:bg-muted/50'}"
			>
				<input
					type="radio"
					class="hidden"
					bind:group={value}
					value={option.value}
				/>
				{option.label}
			</Label>
		{/each}
	</div>
</div>
