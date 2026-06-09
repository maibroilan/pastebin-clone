const API_URL = import.meta.env.VITE_API_URL;

export async function api(path: string, options?: RequestInit) {
	try {
		const response = await fetch(`${API_URL}${path}`, {
			headers: {
				'Content-Type': 'application/json',
				...options?.headers
			},
			...options
		});

		return response;
	} catch (err) {
		return new Response(
			JSON.stringify({
				ok: false,
				status: 500,
				message: err
			}),
			{
				status: 500,
				headers: {
					'Content-Type': 'application/json'
				}
			}
		);
	}
}
