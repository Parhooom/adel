<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { adminApi, type ProblemSummary } from '$lib/api.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { truncateText, getDifficultyVariant } from '$lib/utils.js';
	import { PageLayout, Table, Badge, Button, Alert } from '$lib/components/index.js';

	let problems = $state<ProblemSummary[]>([]);
	let loading = $state(true);
	let error = $state('');
	let success = $state('');
	let deletingId = $state<number | null>(null);

	$effect(() => {
		if (!authStore.isAuthenticated) {
			goto('/login');
		} else if (!authStore.isAdmin) {
			goto('/');
		}
	});

	onMount(async () => {
		await loadProblems();
	});

	async function loadProblems() {
		try {
			loading = true;
			if (!authStore.token) {
				error = 'You must be logged in to view problems';
				return;
			}
			const response = await adminApi.getAllProblemsForAdmin(authStore.token);
			problems = response.problems;
		} catch (err: any) {
			error = err.message || 'Failed to load problems';
		} finally {
			loading = false;
		}
	}

	async function deleteProblem(id: number, title: string) {
		if (!confirm(`Are you sure you want to delete "${title}"? This action cannot be undone.`)) {
			return;
		}

		if (!authStore.token) {
			error = 'You must be logged in to delete problems';
			return;
		}

		try {
			deletingId = id;
			await adminApi.deleteProblem(id, authStore.token);
			success = `Problem "${title}" deleted successfully`;

			problems = problems.filter((p) => p.id !== id);
		} catch (err: any) {
			error = err.message || 'Failed to delete problem';
		} finally {
			deletingId = null;
		}
	}

	function clearMessages() {
		error = '';
		success = '';
	}
</script>

<PageLayout
	isAdmin={true}
	currentPage="problems"
	title="Manage Problems"
	subtitle="Create, edit, and manage coding problems"
>
	{#if error}
		<div class="mb-6">
			<Alert type="error" message={error} dismissible ondismiss={clearMessages} />
		</div>
	{/if}

	{#if success}
		<div class="mb-6">
			<Alert type="success" message={success} dismissible ondismiss={clearMessages} />
		</div>
	{/if}

	<div class="mb-6 flex justify-end">
		<a href="/admin/problems/create">
			<Button variant="primary">Create New Problem</Button>
		</a>
	</div>

	<div class="rounded-lg bg-white shadow-sm">
		<Table
			headers={['ID', 'Title', 'Description', 'Difficulty', 'Status', 'Actions']}
			{loading}
			empty={problems.length === 0}
			emptyMessage="No problems created yet"
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
						<Badge variant={getDifficultyVariant(problem.difficulty)}>
							{problem.difficulty}
						</Badge>
					</td>
					<td class="whitespace-nowrap px-6 py-4">
						<Badge variant={problem.is_active ? 'success' : 'error'}>
							{problem.is_active ? 'Active' : 'Inactive'}
						</Badge>
					</td>
					<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
						<div class="flex space-x-2">
							<a href="/admin/problems/{problem.id}/edit">
								<Button variant="ghost" size="sm">Edit</Button>
							</a>
							<Button
								variant="danger"
								size="sm"
								disabled={deletingId === problem.id}
								loading={deletingId === problem.id}
								onclick={() => deleteProblem(problem.id, problem.title)}
							>
								Delete
							</Button>
						</div>
					</td>
				</tr>
			{/each}
		</Table>
	</div>
</PageLayout>
