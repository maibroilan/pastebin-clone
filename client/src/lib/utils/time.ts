export function timeRemaining(targetDateString: string) {
	const now = new Date();
	const target = new Date(targetDateString);

	const diff = target.getTime() - now.getTime();

	// If already passed
	if (diff <= 0) {
		return 'Expired';
	}

	const seconds = Math.floor(diff / 1000) % 60;
	const minutes = Math.floor(diff / (1000 * 60)) % 60;
	const hours = Math.floor(diff / (1000 * 60 * 60)) % 24;
	const days = Math.floor(diff / (1000 * 60 * 60 * 24));

	if (days > 0) {
		return `${days}d ${hours}h ${minutes}m`;
	}

	if (hours > 0) {
		return `${hours}h ${minutes}m ${seconds}s`;
	}

	if (minutes > 0) {
		return `${minutes}m ${seconds}s`;
	}

	return `${seconds}s`;
}
