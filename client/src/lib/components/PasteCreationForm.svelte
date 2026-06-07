<script lang="ts">
	import { enhance } from '$app/forms';
</script>

<form method="POST" action="?/create" use:enhance>
	<textarea name="content" required></textarea>
	<select name="expiration">
		<option value="30m">30 minutes</option>
		<option value="1h">1 hour</option>
		<option value="1d">1 day</option>
		<option value="7d">7 days</option>
		<option value="10d">10 days</option>
		<option value="30d">30 days</option>
	</select>
	<input name="password" type="password" autocomplete="off" />
	<input type="submit" />
</form>
