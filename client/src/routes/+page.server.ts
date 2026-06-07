import { createPaste } from '$lib/api/pastes.js';
import { fail } from '@sveltejs/kit';

export const actions = {
	create: async ({ request }) => {
		const data = await request.formData();

		const content = data.get('content')?.toString();
		const expiration = data.get('expiration')?.toString();
		const password = data.get('password')?.toString();

		if (!content || !expiration) {
			return fail(400);
		}

		const response = await createPaste({
			content,
			expiration,
			password
		});

		if (!response.ok) {
			return fail(400, {
				message: response.message
			});
		}

		return {
			slug: response.slug,
			expires_at: response.expires_at
		};
	}
};
