<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { ArrowRight, BarChart3, Clock, Database, Shield } from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
</script>

<svelte:head>
	<title>DB Backup</title>
</svelte:head>

<div
	class="flex min-h-screen flex-col bg-background selection:bg-primary selection:text-primary-foreground"
>
	<!-- Header -->
	<header class="container mx-auto flex items-center justify-between py-6">
		<div class="flex items-center gap-2 text-xl font-bold tracking-tighter">
			<Database class="h-6 w-6 text-primary" />
			<span>DbBackup</span>
		</div>
		<nav>
			<Button variant="ghost" href="https://github.com/ariefsn/db-backups" target="_blank"
				>GitHub</Button
			>
			<Button href="/dashboard">Dashboard</Button>
		</nav>
	</header>

	<!-- Hero -->
	<main class="flex-1">
		<section
			class="container mx-auto grid min-h-[80vh] grid-cols-1 items-center gap-12 lg:grid-cols-2"
		>
			<div class="space-y-6" in:fly={{ y: 20, duration: 800, delay: 200 }}>
				<div
					class="inline-block rounded-full bg-secondary px-3 py-1 text-sm font-medium text-secondary-foreground"
				>
					v1.0.0 Now Available
				</div>
				<h1 class="text-5xl font-extrabold tracking-tight lg:text-7xl">
					Secure Your Data <br />
					<span class="bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent"
						>With Confidence</span
					>
				</h1>
				<p class="max-w-[600px] text-lg text-muted-foreground md:text-xl">
					The ultimate tool for managing database backups locally and in the cloud. Support for
					PostgreSQL, MySQL, MongoDB, and Redis.
				</p>
				<div class="flex gap-4">
					<Button size="lg" href="/dashboard" class="group">
						Get Started
						<ArrowRight class="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
					</Button>
					<Button size="lg" variant="outline" href="https://github.com/ariefsn/db-backups"
						>Documentation</Button
					>
				</div>
			</div>

			<div class="relative hidden lg:block" in:fade={{ duration: 1000, delay: 400 }}>
				<div
					class="absolute -inset-1 rounded-2xl bg-gradient-to-r from-blue-600 to-indigo-600 opacity-20 blur-3xl"
				></div>
				<Card.Root class="relative overflow-hidden border-border/50 bg-card/50 backdrop-blur-sm">
					<Card.Header>
						<Card.Title>Backup Status</Card.Title>
						<Card.Description>Real-time overview of your database backups.</Card.Description>
					</Card.Header>
					<Card.Content class="space-y-4">
						<div class="flex items-center justify-between rounded-lg border p-4">
							<div class="flex items-center gap-3">
								<div class="rounded-full bg-green-500/10 p-2 text-green-500">
									<Database class="h-4 w-4" />
								</div>
								<div>
									<p class="font-medium">PostgreSQL Production</p>
									<p class="text-xs text-muted-foreground">Last backup: 10 mins ago</p>
								</div>
							</div>
							<div
								class="rounded-full bg-green-500/10 px-2.5 py-0.5 text-xs font-medium text-green-500"
							>
								Completed
							</div>
						</div>

						<div class="flex items-center justify-between rounded-lg border p-4">
							<div class="flex items-center gap-3">
								<div class="rounded-full bg-blue-500/10 p-2 text-blue-500">
									<Database class="h-4 w-4" />
								</div>
								<div>
									<p class="font-medium">MongoDB Analytics</p>
									<p class="text-xs text-muted-foreground">Last backup: 1 hour ago</p>
								</div>
							</div>
							<div
								class="rounded-full bg-blue-500/10 px-2.5 py-0.5 text-xs font-medium text-blue-500"
							>
								Syncing
							</div>
						</div>

						<div class="flex items-center justify-between rounded-lg border p-4 opacity-50">
							<div class="flex items-center gap-3">
								<div class="rounded-full bg-orange-500/10 p-2 text-orange-500">
									<Database class="h-4 w-4" />
								</div>
								<div>
									<p class="font-medium">Redis Cache</p>
									<p class="text-xs text-muted-foreground">Scheduled: in 2 hours</p>
								</div>
							</div>
							<div
								class="rounded-full bg-secondary px-2.5 py-0.5 text-xs font-medium text-secondary-foreground"
							>
								Pending
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</section>

		<!-- Features -->
		<section class="container mx-auto py-24">
			<div class="grid gap-8 md:grid-cols-3">
				<Card.Root class="bg-card/50 backdrop-blur-sm transition-colors hover:bg-card">
					<Card.Header>
						<Shield class="mb-4 h-12 w-12 text-primary" />
						<Card.Title>Secure Storage</Card.Title>
						<Card.Description>
							AES-256 encryption support with seamless integration for Cloudflare R2 and AWS S3
							storage.
						</Card.Description>
					</Card.Header>
				</Card.Root>

				<Card.Root class="bg-card/50 backdrop-blur-sm transition-colors hover:bg-card">
					<Card.Header>
						<Clock class="mb-4 h-12 w-12 text-primary" />
						<Card.Title>Automated Scheduling</Card.Title>
						<Card.Description>
							Set it and forget it. Flexible cron-style scheduling for all your database instances.
						</Card.Description>
					</Card.Header>
				</Card.Root>

				<Card.Root class="bg-card/50 backdrop-blur-sm transition-colors hover:bg-card">
					<Card.Header>
						<BarChart3 class="mb-4 h-12 w-12 text-primary" />
						<Card.Title>Insightful Reports</Card.Title>
						<Card.Description>
							Comprehensive analytics and status reports to keep you informed about your data
							health.
						</Card.Description>
					</Card.Header>
				</Card.Root>
			</div>
		</section>
	</main>

	<!-- Footer -->
	<footer class="border-t py-6 md:py-8">
		<div class="container mx-auto flex flex-col items-center justify-between gap-4 md:flex-row">
			<p class="text-sm text-muted-foreground">© 2025 DbBackup Manager. All rights reserved.</p>
			<!-- <div class="flex gap-4">
				<a href="/" class="text-sm font-medium hover:underline">Privacy</a>
				<a href="/" class="text-sm font-medium hover:underline">Terms</a>
			</div> -->
		</div>
	</footer>
</div>
