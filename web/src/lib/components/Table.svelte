<script lang="ts">
	interface Props {
		headers: string[];
		children: any;
		loading?: boolean;
		empty?: boolean;
		emptyMessage?: string;
	}

	let {
		headers,
		children,
		loading = false,
		empty = false,
		emptyMessage = 'No data available'
	}: Props = $props();
</script>

<div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
	<table class="min-w-full divide-y divide-gray-300">
		<thead class="bg-gray-50">
			<tr>
				{#each headers as header}
					<th
						class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
					>
						{header}
					</th>
				{/each}
			</tr>
		</thead>
		<tbody class="divide-y divide-gray-200 bg-white">
			{#if loading}
				<tr>
					<td colspan={headers.length} class="px-6 py-12 text-center">
						<div class="flex items-center justify-center">
							<svg class="h-8 w-8 animate-spin text-blue-600" fill="none" viewBox="0 0 24 24">
								<circle
									class="opacity-25"
									cx="12"
									cy="12"
									r="10"
									stroke="currentColor"
									stroke-width="4"
								></circle>
								<path
									class="opacity-75"
									fill="currentColor"
									d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
								></path>
							</svg>
							<span class="ml-2 text-gray-500">Loading...</span>
						</div>
					</td>
				</tr>
			{:else if empty}
				<tr>
					<td colspan={headers.length} class="px-6 py-12 text-center text-gray-500">
						{emptyMessage}
					</td>
				</tr>
			{:else}
				{@render children()}
			{/if}
		</tbody>
	</table>
</div>
