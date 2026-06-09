<script lang="ts">
	import { enhance } from '$app/forms';

	const scrollToTop = () => {
		window.scrollTo({ top: 0, behavior: 'smooth' });
	};
</script>

<form
	class="space-y-3 rounded-2xl border border-zinc-800 p-6"
	method="POST"
	action="?/create"
	use:enhance={() => {
		return async ({ update }) => {
			await update();
			scrollToTop();
		};
	}}
>
	<textarea
		class="min-h-75 w-full rounded-xl border border-zinc-800 bg-zinc-900 p-4"
		name="content"
		required
	></textarea>

	<label class="text-sm text-zinc-400" for="expiration">⏰ Expiration</label>
	<select
		class="w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2"
		id="expiration"
		name="expiration"
	>
		<option value="30m">30 minutes</option>
		<option value="1h">1 hour</option>
		<option value="1d">1 day</option>
		<option value="7d">7 days</option>
		<option value="10d">10 days</option>
		<option value="30d">30 days</option>
	</select>

	<section>
		<label class="text-sm text-zinc-400" for="password">🔒 Password (optional)</label>
		<input
			class="w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2"
			name="password"
			id="password"
			type="password"
			autocomplete="off"
		/>
	</section>

	<button
		class="mt-2 w-full rounded-lg bg-purple-600 px-4 py-2 font-medium transition-colors duration-150 hover:cursor-pointer hover:bg-purple-500"
		type="submit"
	>
		Create
	</button>
</form>
