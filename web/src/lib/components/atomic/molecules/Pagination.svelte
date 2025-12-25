<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { ChevronLeft, ChevronRight, MoreHorizontal } from '@lucide/svelte';

	interface Props {
		page: number;
		limit: number;
		total: number;
		onPageChange: (page: number) => void;
		class?: string;
	}

	let { page, limit, total, onPageChange, class: className }: Props = $props();

	const totalPages = $derived(Math.ceil(total / limit));
	const hasPrevious = $derived(page > 1);
	const hasNext = $derived(page < totalPages);

	function getPages() {
		const pages: (number | 'ellipsis')[] = [];
		const maxVisible = 5;

		if (totalPages <= maxVisible) {
			for (let i = 1; i <= totalPages; i++) pages.push(i);
		} else {
			pages.push(1);
			if (page > 3) pages.push('ellipsis');
			
			const start = Math.max(2, page - 1);
			const end = Math.min(totalPages - 1, page + 1);
			
			for (let i = start; i <= end; i++) {
				if (!pages.includes(i)) pages.push(i);
			}
			
			if (page < totalPages - 2) pages.push('ellipsis');
			if (!pages.includes(totalPages)) pages.push(totalPages);
		}
		return pages;
	}

	const pages = $derived(getPages());
</script>

{#if totalPages > 1}
	<nav class={cn('flex items-center justify-center space-x-2 py-4', className)} aria-label="Pagination">
		<Button
			variant="outline"
			size="icon"
			onclick={() => onPageChange(page - 1)}
			disabled={!hasPrevious}
			aria-label="Previous page"
		>
			<ChevronLeft class="h-4 w-4" />
		</Button>

		<div class="flex items-center space-x-1">
			{#each pages as p}
				{#if p === 'ellipsis'}
					<div class="flex h-9 w-9 items-center justify-center">
						<MoreHorizontal class="h-4 w-4 text-muted-foreground" />
					</div>
				{:else}
					<Button
						variant={page === p ? 'default' : 'outline'}
						size="icon"
						onclick={() => onPageChange(p as number)}
						aria-label="Page {p}"
						aria-current={page === p ? 'page' : undefined}
					>
						{p}
					</Button>
				{/if}
			{/each}
		</div>

		<Button
			variant="outline"
			size="icon"
			onclick={() => onPageChange(page + 1)}
			disabled={!hasNext}
			aria-label="Next page"
		>
			<ChevronRight class="h-4 w-4" />
		</Button>
	</nav>
{/if}
