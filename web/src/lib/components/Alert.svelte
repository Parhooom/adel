<script lang="ts">
	interface Props {
		type: 'success' | 'error' | 'warning' | 'info';
		message: string;
		dismissible?: boolean;
		ondismiss?: () => void;
	}

	let { type, message, dismissible = false, ondismiss }: Props = $props();

	const typeStyles = {
		success: 'bg-green-50 border-green-200 text-green-800',
		error: 'bg-red-50 border-red-200 text-red-800',
		warning: 'bg-yellow-50 border-yellow-200 text-yellow-800',
		info: 'bg-blue-50 border-blue-200 text-blue-800'
	};

	const iconPaths = {
		success: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
		error: 'M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
		warning: 'M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
		info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
	};

	const iconColors = {
		success: 'text-green-400',
		error: 'text-red-400',
		warning: 'text-yellow-400',
		info: 'text-blue-400'
	};
</script>

<div class="rounded-md border p-4 {typeStyles[type]}">
	<div class="flex">
		<div class="flex-shrink-0">
			<svg class="h-5 w-5 {iconColors[type]}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={iconPaths[type]} />
			</svg>
		</div>
		<div class="ml-3 flex-1">
			<p class="text-sm font-medium">{message}</p>
		</div>
		{#if dismissible && ondismiss}
			<div class="ml-auto pl-3">
				<div class="-mx-1.5 -my-1.5">
					<button
						type="button"
						class="inline-flex rounded-md p-1.5 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-600 focus:ring-offset-2"
						onclick={ondismiss}
					>
						<span class="sr-only">Dismiss</span>
						<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							/>
						</svg>
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
