<script lang="ts">
	import ModeToggle from '$lib/components/ModeToggle.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { ArrowRight, BarChart3, Clock, Database, Shield } from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
</script>

<svelte:head>
	<title>DB Backup Manager</title>
</svelte:head>

<div
	class="flex min-h-screen flex-col bg-background text-foreground selection:bg-primary selection:text-primary-foreground transition-colors duration-300"
>
	<!-- Header -->
	<header class="container mx-auto flex items-center justify-between py-6">
		<div class="flex items-center gap-2 text-xl font-bold tracking-tighter">
			<Database class="h-6 w-6 text-primary" />
			<span>DbBackup</span>
		</div>
		<nav class="flex items-center gap-2">
			<Button variant="ghost" href="https://github.com/ariefsn/db-backups" target="_blank"
				>GitHub</Button
			>
			<ModeToggle />
			<Button href="/dashboard">Dashboard</Button>
		</nav>
	</header>

	<!-- Hero -->
	<main class="flex-1">
		<section class="relative overflow-hidden pt-16 pb-24 md:pt-24 md:pb-32">
			<div class="container relative z-10 mx-auto grid grid-cols-1 items-center gap-12 lg:grid-cols-2">
				<div class="space-y-8" in:fly={{ y: 20, duration: 800, delay: 200 }}>
					<div class="inline-flex items-center gap-2 rounded-full bg-primary/10 px-4 py-1.5 text-sm font-semibold text-primary">
						<span class="relative flex h-2 w-2">
							<span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
							<span class="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
						</span>
						New: Real-time status tracking
					</div>
					<h1 class="text-6xl font-black leading-[1.1] tracking-tight md:text-8xl">
						Enterprise <br />
						<span class="bg-linear-to-br from-blue-500 via-indigo-600 to-violet-700 bg-clip-text text-transparent">
							Backup Engine
						</span>
					</h1>
					<p class="max-w-[600px] text-xl leading-relaxed text-muted-foreground/80">
						Automate, monitor, and scale your database backups with ease. A professional-grade toolkit for PostgreSQL, MySQL, MongoDB, and Redis.
					</p>
					<div class="flex flex-wrap gap-4 pt-4">
						<Button size="lg" href="/dashboard" class="h-14 px-8 text-base font-bold shadow-xl shadow-primary/20 group">
							Launch Dashboard
							<ArrowRight class="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1" />
						</Button>
						<Button size="lg" variant="outline" href="https://github.com/ariefsn/db-backups" class="h-14 px-8 text-base font-bold">
							View Source
						</Button>
					</div>
				</div>

				<div class="relative" in:fade={{ duration: 1000, delay: 400 }}>
					<div class="absolute -inset-4 rounded-4xl bg-linear-to-tr from-blue-600 to-indigo-600 opacity-20 blur-3xl"></div>
					<div class="relative rounded-2xl border border-border bg-card shadow-2xl backdrop-blur-xl">
						<div class="flex items-center gap-2 border-b border-border bg-muted/30 px-4 py-3">
							<div class="flex gap-1.5">
								<div class="h-3 w-3 rounded-full bg-destructive/50"></div>
								<div class="h-3 w-3 rounded-full bg-orange-500/50"></div>
								<div class="h-3 w-3 rounded-full bg-green-500/50"></div>
							</div>
							<div class="mx-auto text-[10px] font-medium tracking-widest text-muted-foreground uppercase text-center">LIVE MONITOR</div>
						</div>
						<div class="p-6 space-y-4">
							{#each [
								{ name: 'Production DB', type: 'PostgreSQL', status: 'Completed', time: '2m ago', color: 'green' },
								{ name: 'Analytics Cluster', type: 'MongoDB', status: 'In Progress', time: 'Active', color: 'blue' },
								{ name: 'Session Cache', type: 'Redis', status: 'Pending', time: 'Scheduled', color: 'slate' }
							] as item}
								<div class="flex items-center justify-between rounded-xl border border-border bg-muted/10 p-4 transition-colors hover:bg-muted/20">
									<div class="flex items-center gap-4">
										<div class="rounded-lg bg-{item.color}-500/10 p-2.5 text-{item.color}-500">
											<Database class="h-5 w-5" />
										</div>
										<div>
											<p class="font-bold text-sm">{item.name}</p>
											<p class="text-xs text-muted-foreground">{item.type} • {item.time}</p>
										</div>
									</div>
									<div class="rounded-full bg-{item.color}-500/10 px-3 py-1 text-[10px] font-bold uppercase tracking-wider text-{item.color}-500 border border-{item.color}-500/20">
										{item.status}
									</div>
								</div>
							{/each}
						</div>
					</div>
				</div>
			</div>

			<!-- Background Blurs -->
			<div class="pointer-events-none absolute left-1/2 top-0 -z-10 h-[1000px] w-[1000px] -translate-x-1/2 bg-linear-to-b from-blue-500/10 to-transparent opacity-50 blur-[120px]"></div>
		</section>

		<!-- Features Grid -->
		<section class="container mx-auto py-24 md:py-32">
			<div class="mb-16 text-center space-y-4">
				<h2 class="text-3xl font-black md:text-5xl">Built for Reliability</h2>
				<p class="mx-auto max-w-2xl text-lg text-muted-foreground">Every component is engineered to ensure your data is safe, accessible, and easily managed.</p>
			</div>
			
			<div class="grid gap-6 md:grid-cols-4 md:grid-rows-2 h-full">
				<!-- Big Card -->
				<Card.Root class="md:col-span-2 md:row-span-2 bg-card border-border overflow-hidden group shadow-sm">
					<div class="pointer-events-none absolute inset-0 bg-linear-to-br from-blue-600/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
					<Card.Header class="p-8">
						<div class="mb-6 rounded-2xl bg-blue-500/10 p-4 text-blue-500 w-fit">
							<Shield class="h-10 w-10" />
						</div>
						<Card.Title class="text-3xl font-black">Multi-Layer Security</Card.Title>
						<Card.Description class="text-lg leading-relaxed pt-2">
							Your backups are encrypted at rest and in transit. Seamlessly store your data in Cloudflare R2, AWS S3, or local storage with full control over your keys.
						</Card.Description>
					</Card.Header>
				</Card.Root>

				<!-- Small Cards -->
				<Card.Root class="md:col-span-2 bg-card border-border group shadow-sm">
					<Card.Header class="p-8">
						<Clock class="mb-4 h-8 w-8 text-indigo-500" />
						<Card.Title class="text-xl font-bold">Smart Scheduling</Card.Title>
						<Card.Description>
							Granular cron expressions for every database. Automate your backup windows without touching a server.
						</Card.Description>
					</Card.Header>
				</Card.Root>

				<Card.Root class="md:col-span-1 bg-card border-border group shadow-sm">
					<Card.Header class="p-8">
						<BarChart3 class="mb-4 h-8 w-8 text-violet-500" />
						<Card.Title class="text-xl font-bold">Live Logs</Card.Title>
						<Card.Description>
							Real-time status tracking for every backup job.
						</Card.Description>
					</Card.Header>
				</Card.Root>

				<Card.Root class="md:col-span-1 bg-card border-border group shadow-sm">
					<Card.Header class="p-8">
						<Database class="mb-4 h-8 w-8 text-emerald-500" />
						<Card.Title class="text-xl font-bold">Multi-DB</Card.Title>
						<Card.Description>
							Postgre, MySQL, MongoDB, and Redis.
						</Card.Description>
					</Card.Header>
				</Card.Root>
			</div>
		</section>

		<!-- Technical Breakdown -->
		<section class="border-y border-border bg-muted/5 py-24">
			<div class="container mx-auto">
				<div class="grid gap-12 lg:grid-cols-2 lg:items-center">
					<div class="space-y-8">
						<h2 class="text-4xl font-black md:text-6xl">Production Ready Infrastructure</h2>
						<ul class="space-y-6">
							{#each [
								{ title: 'Atomic Architecture', desc: 'Built with a modular frontend and a high-performance Go backend.' },
								{ title: 'Stateless Design', desc: 'Easily deployable via Docker for consistent results across environments.' },
								{ title: 'Automated Cleanup', desc: 'Managed lifecycle policies to prevent storage bloat automatically.' }
							] as feature}
								<li class="flex gap-4">
									<div class="mt-1 flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
										<ArrowRight class="h-4 w-4" />
									</div>
									<div class="space-y-1">
										<p class="font-bold text-lg">{feature.title}</p>
										<p class="text-muted-foreground">{feature.desc}</p>
									</div>
								</li>
							{/each}
						</ul>
					</div>
					<div class="relative aspect-square lg:aspect-video rounded-3xl border border-border bg-card p-8 shadow-2xl">
						<div class="font-mono text-xs text-primary/60">
							<pre>
{`{
  "status": "success",
  "engine": "docker-optimized",
  "databases": [
    "postgresql",
    "mysql",
    "mongodb",
    "redis"
  ],
  "storage": "r2-compatible",
  "monitoring": "real-time"
}`}
							</pre>
						</div>
						<div class="absolute bottom-12 right-12 h-24 w-24 rounded-full bg-blue-600 blur-[60px] opacity-10"></div>
					</div>
				</div>
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
