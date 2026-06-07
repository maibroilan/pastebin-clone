export interface apiError {
	ok: false;
	status: number;
	message: string;
}

export interface apiOk {
	ok: true;
	status: number;
	message: string;
}
