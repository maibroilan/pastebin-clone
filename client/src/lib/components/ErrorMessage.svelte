<script lang="ts">
	let { error } = $props();
</script>

<section
	class="mb-3 flex flex-col justify-between rounded-lg border border-red-800/40 bg-red-950/30 p-4 md:flex-row"
>
	<p>{error}</p>
</section>
