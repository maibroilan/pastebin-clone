import { getPaste } from '$lib/api/pastes.js';

async function resolvePaste(
	slug: string,
	password?: string,
	wrongPasswordState = 'needs_password'
) {
	const res = await getPaste(slug, password);

	if (res.status === 401) {
		return {
			state: wrongPasswordState,
			message: wrongPasswordState === 'wrong_password' ? 'Wrong Password !' : undefined
		};
	}

	if (res.status === 404) {
		return {
			state: 'error',
			message: 'Paste Not Found !'
		};
	}

	if (res.status === 410) {
		return {
			state: 'error',
			message: 'Paste Expired !'
		};
	}

	if (!res.ok) {
		return {
			state: 'error',
			message: res.message
		};
	}

	return {
		state: 'unlocked',
		paste: {
			content: res.content,
			expires_at: res.expires_at
		}
	};
}

export const load = async ({ params }) => {
	return resolvePaste(params.slug);
};

export const actions = {
	unlock: async ({ request, params }) => {
		const data = await request.formData();

		return resolvePaste(params.slug, data.get('password')?.toString(), 'wrong_password');
	}
};
