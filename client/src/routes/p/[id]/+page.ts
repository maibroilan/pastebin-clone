import { api } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const { id } = params;

	const paste = await api(`/pastes/${id}`);

	return { paste };
};
