export function truncateText(text: string, maxLength = 100): string {
	if (text.length <= maxLength) return text;
	return text.substring(0, maxLength) + '...';
}

export function formatDate(dateString: string): string {
	try {
		if (dateString === '0001-01-01T00:00:00Z') return 'N/A';
		return new Date(dateString).toLocaleString();
	} catch {
		return 'N/A';
	}
}

export function formatExecutionTime(timeMs: number): string {
	const defaultTimeMs = 100;
	if (timeMs === 0) return `${defaultTimeMs}ms`;
	return `${timeMs}ms`;
}

export function formatMemoryUsage(memoryMb: number): string {
	const defaultMemoryMb = 256;
	if (memoryMb === 0) return `${defaultMemoryMb}MB`;
	return `${memoryMb}MB`;
}

export function getDifficultyVariant(difficulty: string): 'easy' | 'medium' | 'hard' {
	switch (difficulty.toLowerCase()) {
		case 'easy':
			return 'easy';
		case 'medium':
			return 'medium';
		case 'hard':
			return 'hard';
		default:
			return 'medium';
	}
}

export function getStatusVariant(status: string): 'success' | 'error' | 'warning' | 'info' {
	switch (status.toLowerCase()) {
		case 'accepted':
		case 'correct':
			return 'success';
		case 'wrong_answer':
		case 'runtime_error':
		case 'compilation_error':
			return 'error';
		case 'time_limit_exceeded':
		case 'memory_limit_exceeded':
			return 'warning';
		case 'pending':
		case 'running':
		default:
			return 'info';
	}
}

export function validateForm(
	fields: { [key: string]: any },
	rules: { [key: string]: (value: any) => string }
): { [key: string]: string } {
	const errors: { [key: string]: string } = {};

	for (const [field, value] of Object.entries(fields)) {
		if (rules[field]) {
			const error = rules[field](value);
			if (error) {
				errors[field] = error;
			}
		}
	}

	return errors;
}

export function isFormValid(errors: { [key: string]: string }): boolean {
	return Object.values(errors).every((error) => error === '');
}
