<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import CopyButton from './CopyButton.svelte';

	let props = $props();
	let slug = $derived(props.slug);
</script>

<section
	class="mb-3 flex flex-col justify-between rounded-lg border border-emerald-800/40 bg-emerald-950/30 p-4 md:flex-row"
>
	<p>Paste created succesfully !</p>

	<div class="flex space-x-2 self-end">
		<a
			class="flex w-20 text-sm font-bold transition-all duration-150 hover:font-extrabold hover:text-purple-500"
			href={resolve('/[slug]', { slug })}
		>
			<p class="w-full self-center text-center">View Paste</p>
		</a>
		<CopyButton variant="success" title="Copy to Clipboard" text={`${page.url.origin}/${slug}`} />
	</div>
</section>
