<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
</script>

<div class="flex min-h-screen items-center justify-center bg-zinc-950 px-4 text-zinc-100">
	<div class="max-w-md text-center">
		<h1 class="text-4xl font-bold tracking-tight">
			{page.status}
		</h1>

		<p class="mt-2 text-zinc-400">
			{page.error?.message ?? 'Something went wrong'}
		</p>

		<a
			href={resolve('/')}
			class="mt-6 inline-block rounded-lg bg-purple-600 px-4 py-2 text-sm font-medium text-white transition-colors duration-150 hover:bg-purple-500"
		>
			Go home
		</a>
	</div>
</div>
