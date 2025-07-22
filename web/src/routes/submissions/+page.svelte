<script lang="ts">
	import { onMount } from 'svelte';
	import { submissionsApi, type Submission } from '$lib/api.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { goto } from '$app/navigation';
	import { formatDate, formatExecutionTime, formatMemoryUsage, truncateText } from '$lib/utils.js';
	import { PageLayout, Table, Alert, StatusBadge } from '$lib/components/index.js';

	let submissions = $state<Submission[]>([]);
	let loading = $state(true);
	let error = $state('');

	$effect(() => {
		if (!authStore.isAuthenticated) {
			goto('/login');
		}
	});

	onMount(async () => {
		if (!authStore.token) {
			error = 'You must be logged in to view submissions';
			loading = false;
			return;
		}

		try {
			const response = await submissionsApi.getUserSubmissions(authStore.token);
			submissions = response.submissions || [];
		} catch (err: any) {
			error = err.message || 'Failed to load submissions';
		} finally {
			loading = false;
		}
	});
</script>

<PageLayout
	currentPage="submissions"
	title="My Submissions"
	subtitle="Track your solution submissions and their status"
>
	{#if error}
		<div class="mb-6">
			<Alert type="error" message={error} />
		</div>
	{/if}

	<div class="rounded-lg bg-white shadow-sm">
		<Table
			headers={['ID', 'Problem', 'Language', 'Code', 'Status', 'Time', 'Memory', 'Submitted']}
			{loading}
			empty={submissions.length === 0}
			emptyMessage="No submissions yet. Start solving problems to see your submissions here!"
		>
			{#each submissions as submission}
				<tr class="hover:bg-gray-50">
					<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
						#{submission.id}
					</td>
					<td class="whitespace-nowrap px-6 py-4">
						<a
							href="/problems/{submission.problem_id}"
							class="text-sm font-medium text-blue-600 hover:text-blue-500 hover:underline"
						>
							Problem #{submission.problem_id}
						</a>
					</td>
					<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
						<span
							class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-medium text-gray-800"
						>
							{submission.language.toUpperCase()}
						</span>
					</td>
					<td class="max-w-md px-6 py-4 text-sm text-gray-500">
						<div
							class="overflow-hidden whitespace-pre-wrap break-all rounded bg-gray-50 px-2 py-1 font-mono text-xs"
							title={submission.code}
						>
							{truncateText(submission.code, 80)}
						</div>
					</td>
					<td class="whitespace-nowrap px-6 py-4">
						<StatusBadge status={submission.status} type="submission" />
					</td>
					<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
						{formatExecutionTime(submission.execution_time_ms)}
					</td>
					<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
						{formatMemoryUsage(submission.memory_usage_mb)}
					</td>
					<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
						{formatDate(submission.created_at)}
					</td>
				</tr>
			{/each}
		</Table>
	</div>
</PageLayout>
