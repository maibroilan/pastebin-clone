<script lang="ts">
	import { enhance } from '$app/forms';

	let { data, form } = $props();

	let state = $derived(form?.state ?? data.state);
	let paste = $derived(form?.paste ?? data.paste);
	let error = $derived(form?.message ?? data.message);
</script>

{#if error || state === 'wrong_password'}
	<p class="error">{error}</p>
{/if}

{#if state === 'unlocked'}
	<h1>Paste</h1>
	<pre>{paste?.content}</pre>
{/if}

{#if state === 'needs_password' || state === 'wrong_password'}
	<form method="POST" action="?/unlock" use:enhance>
		<input name="password" type="password" />
		<button>Unlock</button>
	</form>
{/if}
