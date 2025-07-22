<script lang="ts">
	import { onMount } from 'svelte';
	import { problemsApi, type ProblemSummary } from '$lib/api.js';
	import { truncateText } from '$lib/utils.js';
	import { PageLayout, Table, StatusBadge } from '$lib/components/index.js';

	let problems = $state<ProblemSummary[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const response = await problemsApi.getAllProblems();
			problems = response.problems;
		} catch (err: any) {
			error = err.message || 'Failed to load problems';
		} finally {
			loading = false;
		}
	});
</script>

<PageLayout
	currentPage="problems"
	title="Problems"
	subtitle="Challenge yourself with our collection of algorithmic problems"
>
	{#if error}
		<div class="mb-6 rounded-md bg-red-50 p-4">
			<div class="flex">
				<svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
				<div class="ml-3">
					<p class="text-sm font-medium text-red-800">{error}</p>
				</div>
			</div>
		</div>
	{/if}

	<div class="rounded-lg bg-white shadow-sm">
		<Table
			headers={['ID', 'Title', 'Description', 'Difficulty']}
			{loading}
			empty={problems.length === 0}
			emptyMessage="No problems available yet"
		>
			{#each problems as problem}
				<tr class="hover:bg-gray-50">
					<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
						#{problem.id}
					</td>
					<td class="whitespace-nowrap px-6 py-4">
						<a
							href="/problems/{problem.id}"
							class="text-sm font-medium text-blue-600 hover:text-blue-500 hover:underline"
						>
							{problem.title}
						</a>
					</td>
					<td class="max-w-md px-6 py-4 text-sm text-gray-500">
						<p class="line-clamp-2">
							{truncateText(problem.description.replace(/\n/g, ' '))}
						</p>
					</td>
					<td class="whitespace-nowrap px-6 py-4">
						<StatusBadge status={problem.difficulty} type="problem" />
					</td>
				</tr>
			{/each}
		</Table>
	</div>
</PageLayout>
