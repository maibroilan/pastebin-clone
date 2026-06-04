<!-- src/lib/components/CopyButton.svelte -->
<script lang="ts">
	let { text, classes }: { text: string; classes: string } = $props();

	let status: 'idle' | 'success' = $state('idle');
	let timeoutId: ReturnType<typeof setTimeout> | null = $state(null);

	async function handleCopy() {
		if (!text) return;

		try {
			await navigator.clipboard.writeText(text);
			status = 'success';

			if (timeoutId) clearTimeout(timeoutId);
			timeoutId = setTimeout(() => {
				status = 'idle';
				timeoutId = null;
			}, 1500);
		} catch (err) {
			console.error('Copy failed', err);
		}
	}

	$effect(() => {
		return () => {
			if (timeoutId) clearTimeout(timeoutId);
		};
	});
</script>

<button
	onclick={handleCopy}
	type="button"
	class="{classes} m-2 rounded-md border border-green-600 bg-green-100 p-1 px-2 transition-colors hover:bg-green-300"
>
	{#if status === 'success'}
		<span class="text-lg text-green-600">✅</span>
	{:else}
		<span class="text-lg">📋</span>
	{/if}
</button>
