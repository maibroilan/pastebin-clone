export interface CreatePasteRequest {
	content: string;
	expiration: string;
	password?: string;
}

export interface CreatePasteResponse {
	ok: true;
	status: 200;
	slug: string;
}

export interface GetPasteResponse {
	ok: true;
	status: 200;
	content: string;
	expires_at: string;
}

export interface Paste {
	slug: string;
	content: string;
	expires_at: string;
}
