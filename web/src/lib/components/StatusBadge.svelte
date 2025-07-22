<script lang="ts">
	import Badge from './Badge.svelte';

	interface Props {
		status: string;
		type?: 'submission' | 'problem' | 'user';
	}

	let { status, type = 'submission' }: Props = $props();

	function getVariant(
		status: string,
		type: string
	): 'success' | 'error' | 'warning' | 'info' | 'easy' | 'medium' | 'hard' {
		if (type === 'submission') {
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
		} else if (type === 'problem') {
			switch (status.toLowerCase()) {
				case 'easy':
					return 'easy';
				case 'medium':
					return 'medium';
				case 'hard':
					return 'hard';
				default:
					return 'medium';
			}
		} else if (type === 'user') {
			return status === 'admin' ? 'error' : 'info';
		}
		return 'info';
	}

	function formatStatus(status: string): string {
		return status.replace('_', ' ').toUpperCase();
	}
</script>

<Badge variant={getVariant(status, type)}>
	{formatStatus(status)}
</Badge>
