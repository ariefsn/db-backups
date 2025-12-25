<script lang="ts">
	import { OpenAPI } from '$lib/api';
	import { Toaster } from "$lib/components/ui/sonner";
	import '../app.css';

	OpenAPI.BASE = '/api';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href="/favicon.png" />
</svelte:head>
<Toaster />
{@render children()}
