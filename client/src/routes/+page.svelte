<form
	class="flex w-full flex-col justify-self-center p-5 md:max-w-1/2"
	action="/paste"
	method="POST"
>
	<div class="flex w-full justify-center-safe">
		<textarea
			class="min-h-64 w-full rounded-xl"
			placeholder="Paste your text here..."
			rows="10"
			name="content"
			required
		></textarea>
	</div>

	<div class="md:flex md:flex-row md:justify-between">
		<div class="md:flex md:flex-row md:gap-2">
			<div class="my-2 flex flex-col">
				<label for="expiration">🕒Expires in:</label>
				<select class="rounded-xl" id="expiration" name="expiration">
					<option value="30m">30 minutes</option>
					<option value="1h">1 hour</option>
					<option value="1d">1 day</option>
					<option value="7d">7 days</option>
					<option value="10d">10 days</option>
					<option value="30d">30 day</option>
				</select>
			</div>
			<div class="my-2 flex flex-col">
				<label for="password">🔒Password (optional):</label>
				<input
					class="rounded-xl"
					name="password"
					type="password"
					placeholder="••••••••"
					id="password"
				/>
			</div>
		</div>
		<button class="m-2 my-4 w-full rounded-3xl bg-sky-300 p-2 md:w-auto" type="submit"
			>Create Paste</button
		>
	</div>
</form>
