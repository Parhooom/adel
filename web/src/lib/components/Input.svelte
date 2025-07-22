<script lang="ts">
	interface Props {
		label?: string;
		type?: string;
		placeholder?: string;
		value?: string;
		required?: boolean;
		disabled?: boolean;
		error?: string;
		help?: string;
		id?: string;
	}

	let {
		label,
		type = 'text',
		placeholder,
		value = $bindable(''),
		required = false,
		disabled = false,
		error,
		help,
		id
	}: Props = $props();

	const inputId = id || crypto.randomUUID();

	const baseClasses =
		'block w-full rounded-md border-0 px-3 py-2 text-gray-900 shadow-sm ring-1 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-inset sm:text-sm sm:leading-6 transition-colors';

	const normalClasses = 'ring-gray-300 focus:ring-blue-600';
	const errorClasses = 'ring-red-300 focus:ring-red-600';

	const inputClasses = error ? `${baseClasses} ${errorClasses}` : `${baseClasses} ${normalClasses}`;
</script>

<div class="mb-4">
	{#if label}
		<label for={inputId} class="mb-1 block text-sm font-medium leading-6 text-gray-900">
			{label}
			{#if required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}

	<input id={inputId} {type} {placeholder} {required} {disabled} bind:value class={inputClasses} />

	{#if error}
		<p class="mt-1 text-sm text-red-600">{error}</p>
	{:else if help}
		<p class="mt-1 text-sm text-gray-500">{help}</p>
	{/if}
</div>
