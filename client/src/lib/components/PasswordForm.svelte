<script>
	import { enhance } from '$app/forms';
</script>

<div class="mx-auto max-w-md">
	<div class="rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
		<h1 class="text-xl font-semibold text-zinc-100">🔒 Password required</h1>
		<p class="mt-1 text-sm text-zinc-400">Enter the password to view this paste</p>
		<form method="POST" action="?/unlock" use:enhance>
			<input
				class="mt-4 w-full rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2 text-zinc-100 placeholder:text-zinc-600 focus:ring-2 focus:ring-purple-500 focus:outline-none"
				name="password"
				type="password"
			/>
			<button
				class="mt-4 w-full rounded-lg bg-purple-600 px-3 py-2 text-sm font-medium text-white transition-colors duration-150 hover:bg-purple-500 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-zinc-950 focus:outline-none"
			>
				Unlock
			</button>
		</form>
	</div>
</div>
