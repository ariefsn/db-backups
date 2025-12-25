<script lang="ts">
	import { page } from '$app/stores';
	import ModeToggle from '$lib/components/ModeToggle.svelte';
	import { Button } from '$lib/components/ui/button';
	import { BarChart3, Database, LogOut, Menu, Server } from '@lucide/svelte';

	let { children } = $props();

	const sidebarItems = [
		{ href: '/dashboard', icon: Database, label: 'Backups' },
		{ href: '/dashboard/databases', icon: Server, label: 'My Databases' },
		{ href: '/dashboard/stats', icon: BarChart3, label: 'Reports' }
	];

	function isActive(href: string) {
		if (href === '/dashboard' && $page.url.pathname === '/dashboard') return true;
		if (href !== '/dashboard' && $page.url.pathname.startsWith(href)) return true;
		return false;
	}
</script>

<div class="flex h-screen bg-background text-foreground">
	<!-- Sidebar -->
	<aside class="hidden w-64 flex-col border-r md:flex">
		<div class="flex h-14 items-center border-b px-4 font-bold tracking-tight">
			<Database class="mr-2 h-5 w-5 text-primary" />
			<span>DbBackup</span>
		</div>
		<nav class="flex-1 space-y-1 p-4">
			{#each sidebarItems as item}
				<a
					href={item.href}
					class="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-muted {isActive(
						item.href
					)
						? 'bg-secondary text-secondary-foreground'
						: 'text-muted-foreground'}"
				>
					<item.icon class="h-4 w-4" />
					{item.label}
				</a>
			{/each}
		</nav>
		<div class="flex items-center gap-2 border-t p-4">
			<ModeToggle />
			<Button variant="ghost" class="flex-1 justify-start gap-2" href="/">
				<LogOut class="h-4 w-4" />
				Exit
			</Button>
		</div>
	</aside>

	<!-- Mobile Header (Visible only on small screens) -->
	<!-- TODO: Implement Sheet for mobile sidebar -->

	<main class="flex-1 overflow-y-auto">
		<header class="flex h-14 items-center justify-between border-b px-6 lg:hidden">
			<div class="flex items-center">
				<Button variant="ghost" size="icon" class="mr-2">
					<Menu class="h-5 w-5" />
				</Button>
				<span class="font-bold">DbBackup</span>
			</div>
			<ModeToggle />
		</header>
		<div class="p-6">
			{@render children()}
		</div>
	</main>
</div>
