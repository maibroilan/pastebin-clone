<script lang="ts">
	import { timeRemaining } from '$lib/utils/time';
	import { onDestroy } from 'svelte';
	import CopyButton from './CopyButton.svelte';

	let { expires_at, content } = $props();

	// svelte-ignore state_referenced_locally
	let remaining = $state(timeRemaining(expires_at));

	const interval = setInterval(() => {
		remaining = timeRemaining(expires_at);
	}, 1000);

	onDestroy(() => {
		clearInterval(interval);
	});
</script>

<div class="mx-auto max-w-3xl px-4 pt-5">
	<p class="mb-3 text-sm text-zinc-400">⏰ Expires in : {remaining}</p>
	<div class="flex flex-col rounded-2xl border border-zinc-800 bg-zinc-900 px-6 py-2">
		<div class="self-end">
			<CopyButton variant="paste" text={content} title="Copy Content" />
		</div>
		<pre class="leading-relaxed wrap-break-word whitespace-pre-wrap text-zinc-100">{content}</pre>
	</div>
</div>
