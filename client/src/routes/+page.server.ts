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

		const res = await createPaste({
			content,
			expiration,
			password
		});

		if (!res.ok) {
			return {
				ok: res.ok,
				status: res.status,
				message: res.message
			};
		}

		// if (!res.ok) {
		// 	return fail(400, {
		// 		message: res.message
		// 	});
		// }

		return {
			slug: res.slug
		};
	}
};
