<script lang="ts">
	import ErrorMessage from '$lib/components/ErrorMessage.svelte';
	import PasswordForm from '$lib/components/PasswordForm.svelte';
	import PasteCard from '$lib/components/PasteCard.svelte';

	let { data, form } = $props();

	let state = $derived(form?.state ?? data.state);
	let paste = $derived(form?.paste ?? data.paste);
	let error = $derived(form?.message ?? data.message);
</script>

{#if error || state === 'wrong_password'}
	<ErrorMessage {error} />
{/if}

{#if state === 'unlocked'}
	<PasteCard content={paste?.content} expires_at={paste?.expires_at} />
{/if}

{#if state === 'needs_password' || state === 'wrong_password'}
	<PasswordForm />
{/if}
