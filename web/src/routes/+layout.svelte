<script lang="ts">
	import { OpenAPI } from '$lib/api';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher } from 'mode-watcher';
	import '../app.css';

	OpenAPI.BASE = '/api';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href="/favicon.png" />
</svelte:head>
<ModeWatcher />
<Toaster />
{@render children()}
