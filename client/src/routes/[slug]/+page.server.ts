import { getPaste } from '$lib/api/pastes.js';

export const load = async ({ params }) => {
	const res = await getPaste(params.slug);

	if (res.status === 401) {
		return {
			state: 'needs_password',
			slug: params.slug
		};
	}

	if (!res.ok) {
		console.log('yup not found');
		return {
			state: 'error',
			message: 'Paste not found'
		};
	}

	return {
		state: 'unlocked',
		paste: {
			content: res.content,
			expires_at: res.expires_at
		}
	};
};
