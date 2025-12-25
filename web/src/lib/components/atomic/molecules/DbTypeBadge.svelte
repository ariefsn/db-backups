<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Database } from '@lucide/svelte';

	let { type = '' }: { type?: string } = $props();
</script>

<div class="flex items-center gap-2 font-medium">
	<Database class="h-4 w-4 text-muted-foreground" />
	<Badge variant="outline" class="capitalize">{type || '-'}</Badge>
</div>
