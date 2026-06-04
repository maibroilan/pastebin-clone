<script lang="ts">
	async function handleSubmit(e: Event) {
		e.preventDefault();
	}
</script>

<form class="flex flex-col p-5" onsubmit={handleSubmit}>
	<textarea class="min-h-64" placeholder="Paste your text here..." rows="10" required></textarea>

	<div class="md:flex md:flex-row md:justify-between">
		<div class="md:flex md:flex-row md:gap-2">
			<div class="flex flex-col">
				<label for="expiration">🕒Expires in:</label>
				<select id="expiration">
					<option value="30m">30 minutes</option>
					<option value="1h">1 hour</option>
					<option value="1d">1 day</option>
					<option value="7d">7 days</option>
					<option value="10d">10 days</option>
					<option value="30d">30 day</option>
				</select>
			</div>
			<div class="flex flex-col">
				<label for="password">🔒Password (optional):</label>
				<input id="password" />
			</div>
		</div>
		<button class="m-2 w-full rounded-3xl bg-sky-300 p-2 md:w-auto" type="submit"
			>Create Paste</button
		>
	</div>
</form>
