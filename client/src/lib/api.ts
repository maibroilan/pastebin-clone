const API = 'http://localhost:8080';

export async function api<T>(path: string, options?: RequestInit): Promise<T> {
	const res = await fetch(`${API}${path}`, {
		headers: {
			'Content-Type': 'application/json'
		},
		...options
	});

	if (!res.ok) {
		console.log(res.status);
		throw new Error('Request failed');
	}

	return res.json();
}
