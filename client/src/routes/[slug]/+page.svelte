<script lang="ts">
	let { data } = $props();

	let paste = $derived(data.paste);
	let state = $derived(data.state);
	let error = $derived(data.message);
</script>

{#if error}
	<p class="error">{error}</p>
{/if}

{#if state === 'unlocked'}
	<h1>Paste</h1>
	<pre>{paste?.content}</pre>
{/if}

{#if state === 'needs_password'}
	<p>need password</p>
	<!-- TO BE IMPLEMENTED LATER -->
{/if}
