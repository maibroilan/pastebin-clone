<script lang="ts">
	let { text, title, variant } = $props();

	const style = $derived(
		variant === 'success'
			? 'border-emerald-800/40 bg-emerald-900 hover:bg-emerald-800'
			: 'border-zinc-500/40 bg-zinc-800 hover:bg-zinc-700'
	);

	let copied = $state(false);

	async function copyLink() {
		await navigator.clipboard.writeText(text);

		copied = true;

		setTimeout(() => {
			copied = false;
		}, 2000);
	}
</script>

<button
	class={`${style} flex rounded-lg border p-1 text-sm font-bold hover:cursor-pointer`}
	onclick={copyLink}
	{title}
>
	{#if copied}
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="size-5 self-center"
		>
			<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
		</svg>
	{:else}
		<svg
			width="20"
			height="20"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			class="size-5 self-center"
		>
			<path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path>
			<rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect>
		</svg>
	{/if}
</button>
