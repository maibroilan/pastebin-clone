import { getPaste } from '$lib/api/pastes.js';

export const load = async ({ params }) => {
	const res = await getPaste(params.slug);

	if (res.status === 401) {
		return {
			state: 'needs_password'
			// message: 'password needed'
		};
	}

	if (res.status === 404) {
		return {
			state: 'error',
			message: 'not found'
		};
	}

	if (res.status === 410) {
		return {
			state: 'error',
			message: 'paste expired'
		};
	}

	if (!res.ok) {
		return {
			state: 'error',
			message: 'Status: ' + res.status + ' Unknown Error : ' + res.message
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

export const actions = {
	unlock: async ({ request, params }) => {
		const data = await request.formData();
		const password = data.get('password')?.toString();
		const paste = await getPaste(params.slug, password);

		if (paste.status === 401) {
			return {
				state: 'wrong_password',
				message: 'wrong password'
			};
		}

		if (paste.status === 404) {
			return {
				state: 'error',
				message: 'not found'
			};
		}

		if (paste.status === 410) {
			return {
				state: 'error',
				message: 'paste expired'
			};
		}

		if (!paste.ok) {
			return {
				state: 'error',
				message: 'Status: ' + paste.status + ' Unknown Error : ' + paste.message
			};
		}

		return {
			state: 'unlocked',
			paste
		};
	}
};
