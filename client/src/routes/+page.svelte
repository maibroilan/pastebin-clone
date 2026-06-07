<script lang="ts">
	import PasteCreationForm from '$lib/components/PasteCreationForm.svelte';

	let { form } = $props();

	let dismissed = $state(false);
</script>

{#if form?.message && !dismissed}
	<p class="error">{form.message}</p>
	<button type="button" onclick={() => (dismissed = true)} aria-label="Dismiss error"> × </button>
{/if}

{#if form?.slug}
	<p>Paste created!</p>
	<p>Slug: {form.slug}</p>
	<button type="button" onclick={() => (dismissed = true)} aria-label="Dismiss error"> × </button>
{/if}

<PasteCreationForm />
