<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<header class="border border-black py-10 text-center">
	<h1 class="text-7xl font-bold">Pastebin</h1>
	<p>Create Pastes</p>
</header>

<main class="border border-black">
	{@render children()}
</main>

<footer class="border border-black py-5 text-center">
	<p>Developed by <a href="https://github.com/maibroilan">maibroilan</a></p>
	<p>
		Check out this project at <a href="https://github.com/maibroilan/pastebin-clone"
			>https://github.com/maibroilan/pastebin-clone</a
		>
	</p>
</footer>
