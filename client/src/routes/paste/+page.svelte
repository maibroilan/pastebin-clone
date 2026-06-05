<script lang="ts">
	import { resolve } from '$app/paths';
	import { formatTimeRemaining } from '$lib/utils/date_helpers';

	let props = $props();
	let paste_data = $derived(props.form);
</script>

<article class="flex flex-col p-3 text-center">
	<h3 class="font-bold md:text-3xl">Your Paste has been successfully created !</h3>
	<p>
		access it <a
			class="font-bold text-sky-500 hover:text-sky-800"
			href={resolve(`/p/${paste_data.slug}`)}>here</a
		>
	</p>
	<p>it will be expired in {formatTimeRemaining(paste_data.expires_at)}</p>
</article>
