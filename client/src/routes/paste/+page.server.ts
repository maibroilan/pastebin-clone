import { api } from '$lib/api.js';

export const actions = {
	default: async ({ request }) => {
		const formData = await request.formData();

		const content = formData.get('content');
		const expiration = formData.get('expiration');
		const password = formData.get('password')?.toString() ?? '';

		const payload: Record<string, unknown> = {
			content,
			expiration
		};

		if (password.trim() !== '') {
			payload.password = password;
		}

		const response = await api('/pastes', {
			method: 'POST',
			body: JSON.stringify(payload)
		});

		const result = await response;

		return result;
	}
};
