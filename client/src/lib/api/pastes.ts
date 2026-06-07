import { api } from '$lib/api/client';
import type { apiError } from '$lib/types/error';
import type { CreatePasteRequest, CreatePasteResponse, GetPasteResponse } from '$lib/types/paste';
import { handleError } from './errors';

// Variable only set to true during development for testing error handling
const SIMULATE_FAIL = false;
const SAMPLE_ERROR: apiError = {
	ok: false,
	status: 469,
	message: 'error testing is on'
};

export async function createPaste(
	data: CreatePasteRequest
): Promise<CreatePasteResponse | apiError> {
	if (SIMULATE_FAIL) {
		console.warn('ERROR TESTING IS ON !');
		return SAMPLE_ERROR;
	}

	const password = data.password;

	const response = await api(`pastes`, {
		method: 'POST',
		body: JSON.stringify({
			content: data.content,
			expiration: data.expiration,
			...(password ? { password } : {})
		})
	});

	const error = handleError(response.status);
	if (!error.ok) {
		return error;
	}

	const respjs = await response.json();

	return {
		ok: true,
		slug: respjs.slug,
		expires_at: respjs.expires_at
	};
}

export async function getPaste(
	slug: string,
	password?: string
): Promise<GetPasteResponse | apiError> {
	if (SIMULATE_FAIL) {
		console.warn('ERROR TESTING IS ON !');
		return SAMPLE_ERROR;
	}

	const response = await api(`pastes/${slug}`, {
		method: 'GET',
		headers: password ? { 'X-Paste-Password': password } : {}
	});

	const error = handleError(response.status);
	if (!error.ok) {
		return error;
	}

	const respjs = await response.json();

	return {
		ok: true,
		status: 200,
		content: respjs.content,
		expires_at: respjs.expires_at
	};
}
