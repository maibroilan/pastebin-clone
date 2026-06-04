<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<header class="border border-black text-center">
	<h1 class="text-4xl font-bold">Pastebin</h1>
	<p>create pastes</p>
</header>

<main class="border border-black">
	{@render children()}
</main>

<footer class="border border-black text-center">
	<p>developed by <a href="https://github.com/maibroilan">maibroilan</a></p>
</footer>
