<script lang="ts">
	import ErrorMessage from '$lib/components/ErrorMessage.svelte';
	import PasteCreationForm from '$lib/components/PasteCreationForm.svelte';
	import SuccessMessage from '$lib/components/SuccessMessage.svelte';

	let { form } = $props();
</script>

{#if form?.message}
	<ErrorMessage error={form.message} />
{/if}

{#if form?.slug}
	<SuccessMessage slug={form.slug} />
{/if}

<PasteCreationForm />
