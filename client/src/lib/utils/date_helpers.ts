export function getTimeRemaining(targetDate: Date | string) {
	const target = new Date(targetDate);
	const now = new Date();
	const total = target.getTime() - now.getTime();

	if (total <= 0) {
		return { total: 0, days: 0, hours: 0, minutes: 0, seconds: 0 };
	}

	const seconds = Math.floor((total / 1000) % 60);
	const minutes = Math.floor((total / (1000 * 60)) % 60);
	const hours = Math.floor((total / (1000 * 60 * 60)) % 24);
	const days = Math.floor(total / (1000 * 60 * 60 * 24));

	return { total, days, hours, minutes, seconds };
}

export function formatTimeRemaining(targetDate: Date | string): string {
	const { days, hours, minutes, seconds } = getTimeRemaining(targetDate);

	const units = [];
	if (days > 0) units.push({ value: days, unit: 'day' });
	if (hours > 0) units.push({ value: hours, unit: 'hour' });
	if (minutes > 0) units.push({ value: minutes, unit: 'minute' });
	if (seconds > 0 && units.length === 0) units.push({ value: seconds, unit: 'second' });

	if (units.length === 0) return '0 seconds';

	const topUnits = units.slice(0, 2);
	const formatted = topUnits.map((u) => `${u.value} ${u.unit}${u.value !== 1 ? 's' : ''}`);

	return formatted.join(' and ');
}
