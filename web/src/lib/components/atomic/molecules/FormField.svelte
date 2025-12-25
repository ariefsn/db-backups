<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import type { Snippet } from 'svelte';

	let { 
		label, 
		error, 
		children,
		description
	}: { 
		label?: string, 
		error?: string, 
		description?: string,
		children: Snippet 
	} = $props();
</script>

<div class="space-y-2">
	{#if label}
		<Label>{label}</Label>
	{/if}
	
	{@render children()}

	{#if description}
		<p class="text-xs text-muted-foreground">{description}</p>
	{/if}

	{#if error}
		<p class="text-xs text-destructive">{error}</p>
	{/if}
</div>
