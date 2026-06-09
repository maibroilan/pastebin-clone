<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="min-h-screen bg-zinc-950 p-6 text-zinc-100">
	<header class="mb-6 text-center">
		<h1 class="text-3xl font-bold tracking-tight text-zinc-100">Pastebin</h1>
		<p class="mt-2 text-zinc-400">Share text snippets with optional password protection</p>
	</header>

	<main class="">
		{@render children()}
	</main>

	<footer class="mt-6 border-t border-zinc-800 py-6 text-center text-sm text-zinc-500">
		<p>
			Developed by <a
				class="transition-colors duration-150 hover:text-purple-500"
				href="https://github.com/maibroilan">Maibroilan</a
			>
		</p>
	</footer>
</div>
