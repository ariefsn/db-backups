<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Search } from '@lucide/svelte';

	let { 
		value = $bindable(''), 
		placeholder = 'Search...',
		onkeydown
	}: { 
		value: string, 
		placeholder?: string,
		onkeydown?: (e: KeyboardEvent) => void
	} = $props();
</script>

<div class="relative flex-1">
	<Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
	<Input
		type="search"
		placeholder={placeholder}
		class="pl-8"
		bind:value={value}
		onkeydown={onkeydown}
	/>
</div>
