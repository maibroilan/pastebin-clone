import type { apiError, apiOk } from '$lib/types/error';

export function handleError(status: number): apiError | apiOk {
	if (status === 400) {
		return {
			ok: false,
			status: status,
			message: 'Invalid Request'
		};
	}

	if (status === 401) {
		return {
			ok: false,
			status: status,

			message: 'Unauthorized'
		};
	}

	if (status === 410) {
		return {
			ok: false,
			status: status,

			message: 'Gone'
		};
	}

	if (status === 404) {
		return {
			ok: false,
			status: status,

			message: 'Not Found'
		};
	}

	if (status === 500) {
		return {
			ok: false,
			status: status,
			message: 'Server Error'
		};
	}

	return {
		ok: true,
		status: status,
		message: 'success'
	};
}
