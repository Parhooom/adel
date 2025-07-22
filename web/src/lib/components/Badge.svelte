<script lang="ts">
	interface Props {
		variant?: 'easy' | 'medium' | 'hard' | 'success' | 'warning' | 'error' | 'info';
		children: any;
	}

	let { variant = 'info', children }: Props = $props();

	const variantClasses = {
		easy: 'bg-green-100 text-green-800',
		medium: 'bg-yellow-100 text-yellow-800',
		hard: 'bg-red-100 text-red-800',
		success: 'bg-green-100 text-green-800',
		warning: 'bg-yellow-100 text-yellow-800',
		error: 'bg-red-100 text-red-800',
		info: 'bg-blue-100 text-blue-800'
	};

	const classes = `inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${variantClasses[variant]}`;
</script>

<span class={classes}>
	{@render children()}
</span>
