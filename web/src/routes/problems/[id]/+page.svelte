<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { problemsApi, submissionsApi, type Problem, type SubmissionRequest } from '$lib/api.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import Button from '$lib/components/Button.svelte';
	import Alert from '$lib/components/Alert.svelte';
	import CodeEditor from '$lib/components/CodeEditor.svelte';
	import ResponsiveNavbar from '$lib/components/ResponsiveNavbar.svelte';

	let problem = $state<Problem | null>(null);
	let loading = $state(true);
	let error = $state('');
	let submissionSuccess = $state('');
	let submissionError = $state('');
	let selectedLanguage = $state('c');

	const problemId = $derived(parseInt($page.params.id));

	const languages = [
		{
			value: 'c',
			label: 'C'
		},
		{
			value: 'python',
			label: 'Python'
		},
		{
			value: 'go',
			label: 'Go'
		}
	];

	onMount(async () => {
		try {
			const response = await problemsApi.getProblemById(problemId);
			problem = response.problem;
		} catch (err: any) {
			error = err.message || 'Failed to load problem';
		} finally {
			loading = false;
		}
	});

	function getDifficultyVariant(difficulty: string): 'easy' | 'medium' | 'hard' {
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

	const difficultyColors: Record<string, string> = {
		easy: 'bg-green-100 text-green-800',
		medium: 'bg-yellow-100 text-yellow-800',
		hard: 'bg-red-100 text-red-800'
	};

	function goBack() {
		goto('/problems');
	}

	async function handleSubmission(submission: SubmissionRequest) {
		if (!authStore.token) {
			submissionError = 'You must be logged in to submit solutions';
			return;
		}

		try {
			submissionError = '';
			submissionSuccess = '';

			const response = await submissionsApi.submitSolution(submission, authStore.token);
			submissionSuccess = `Solution submitted successfully! Submission ID: ${response.submission.id}. Status: ${response.submission.status}`;
		} catch (err: any) {
			submissionError = err.message || 'Failed to submit solution';
		}
	}

	function clearMessages() {
		submissionSuccess = '';
		submissionError = '';
	}
</script>

<div class="min-h-screen bg-gray-50">
	<ResponsiveNavbar currentPage="problems" />

	<main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<svg class="h-8 w-8 animate-spin text-blue-600" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
					></circle>
					<path
						class="opacity-75"
						fill="currentColor"
						d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
					></path>
				</svg>
				<span class="ml-2 text-gray-500">Loading problem...</span>
			</div>
		{:else if error}
			<div class="py-12 text-center">
				<Alert type="error" message={error} />
				<div class="mt-4">
					<Button variant="secondary" onclick={goBack}>← Back to Problems</Button>
				</div>
			</div>
		{:else if problem}
			<div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
				<!-- Problem Details -->
				<div class="space-y-6">
					<div class="rounded-lg bg-white p-6 shadow">
						<div class="mb-4 flex items-center justify-between">
							<h1 class="text-2xl font-bold text-gray-900">{problem.title}</h1>
							<span
								class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium {difficultyColors[
									problem.difficulty
								]}"
							>
								{problem.difficulty}
							</span>
						</div>

						<div class="prose max-w-none">
							<div class="whitespace-pre-wrap text-gray-700">{problem.description}</div>
						</div>

						<div class="mt-6 grid grid-cols-2 gap-4 text-sm">
							<div>
								<span class="font-medium text-gray-500">Time Limit:</span>
								<span class="ml-2 text-gray-900">{problem.time_limit_ms}ms</span>
							</div>
							<div>
								<span class="font-medium text-gray-500">Memory Limit:</span>
								<span class="ml-2 text-gray-900">{problem.memory_limit_mb}MB</span>
							</div>
						</div>
					</div>

					<!-- Sample Test Cases -->
					{#if problem.testcases && problem.testcases.length > 0}
						<div class="rounded-lg bg-white p-6 shadow-sm">
							<h2 class="mb-4 text-lg font-semibold text-gray-900">Sample Test Cases</h2>
							<div class="space-y-4">
								{#each problem.testcases.slice(0, 3) as testcase, index}
									<div class="rounded-lg border border-gray-200 p-4">
										<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
											<div>
												<h4 class="mb-2 text-sm font-medium text-gray-700">Input:</h4>
												<pre
													class="overflow-x-auto rounded-md bg-gray-50 p-3 text-sm text-gray-900">{testcase.input_data}</pre>
											</div>
											<div>
												<h4 class="mb-2 text-sm font-medium text-gray-700">Output:</h4>
												<pre
													class="overflow-x-auto rounded-md bg-gray-50 p-3 text-sm text-gray-900">{testcase.output_data}</pre>
											</div>
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>

				<!-- Submission Panel -->
				<div class="space-y-6">
					{#if submissionSuccess}
						<Alert
							type="success"
							message={submissionSuccess}
							dismissible
							ondismiss={clearMessages}
						/>
					{/if}

					{#if submissionError}
						<Alert type="error" message={submissionError} dismissible ondismiss={clearMessages} />
					{/if}

					<div class="rounded-lg bg-white p-6 shadow-sm">
						<div class="mb-4 flex items-center justify-between">
							<h2 class="text-lg font-semibold text-gray-900">Submit Solution</h2>
							{#if authStore.isAuthenticated}
								<div class="flex items-center space-x-2">
									<label for="language-select" class="text-sm font-medium text-gray-700">
										Language:
									</label>
									<select
										id="language-select"
										bind:value={selectedLanguage}
										class="rounded-md border-gray-300 py-1 pl-3 pr-8 text-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
									>
										{#each languages as language}
											<option value={language.value}>{language.label}</option>
										{/each}
									</select>
								</div>
							{/if}
						</div>
						{#if authStore.isAuthenticated}
							<CodeEditor
								problemId={problem.id}
								onsubmit={handleSubmission}
								{loading}
								{selectedLanguage}
							/>
						{:else}
							<div class="py-8 text-center text-gray-500">
								<svg
									class="mx-auto h-12 w-12 text-gray-400"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
									/>
								</svg>
								<h3 class="mt-2 text-sm font-medium text-gray-900">Login Required</h3>
								<p class="mt-1 text-sm text-gray-500">
									You need to sign in to submit solutions and track your progress.
								</p>
								<div class="mt-6 space-x-3">
									<a href="/login">
										<Button variant="primary">Sign In</Button>
									</a>
									<a href="/register">
										<Button variant="secondary">Register</Button>
									</a>
								</div>
							</div>
						{/if}
					</div>

					<!-- Quick Stats -->
					<div class="rounded-lg bg-white p-6 shadow-sm">
						<h2 class="mb-4 text-lg font-semibold text-gray-900">Problem Stats</h2>
						<div class="space-y-3">
							<div class="flex justify-between">
								<span class="text-sm text-gray-600">Total Submissions</span>
								<span class="text-sm font-medium text-gray-900">0</span>
							</div>
							<div class="flex justify-between">
								<span class="text-sm text-gray-600">Accepted</span>
								<span class="text-sm font-medium text-green-600">0</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</main>
</div>
