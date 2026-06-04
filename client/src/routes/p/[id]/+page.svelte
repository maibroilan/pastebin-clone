<script lang="ts">
	import CopyButton from '$lib/components/CopyButton.svelte';
	import PasteCard from '$lib/components/PasteCard.svelte';
	import { formatTimeRemaining } from '$lib/utils/date_helpers';
	import { onDestroy } from 'svelte';

	let props = $props();
	let paste = $derived(props.data.paste);

	// countdown that updates every second
	let formatted = $derived(formatTimeRemaining(paste.expires_at));
	let interval;
	function updateCountdown() {
		formatted = formatTimeRemaining(paste.expires_at);
	}
	interval = setInterval(updateCountdown, 1000);
	updateCountdown();
	onDestroy(() => {
		if (interval) clearInterval(interval);
	});
</script>

<div class="flex flex-col">
	<CopyButton classes="self-end" text={paste.content} />
	<PasteCard content={paste.content} />
	<p class="m-1 p-1">🕒 {formatted} until expiration</p>
</div>
