<script lang="ts">
	import { goto } from '$app/navigation';
	import { adminApi, type ProblemCreateRequest, type TestCaseRequest } from '$lib/api.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import AdminNavbar from '$lib/components/AdminNavbar.svelte';
	import Input from '$lib/components/Input.svelte';
	import Button from '$lib/components/Button.svelte';
	import Alert from '$lib/components/Alert.svelte';

	$effect(() => {
		if (!authStore.isAuthenticated) {
			goto('/login');
		} else if (!authStore.isAdmin) {
			goto('/');
		}
	});

	let title = $state('');
	let description = $state('');
	let difficulty = $state('easy');
	let timeLimitMs = $state(3000);
	let memoryLimitMb = $state(256);
	let isActive = $state(true);

	let testCases = $state<TestCaseRequest[]>([{ input_data: '', output_data: '', is_active: true }]);

	let isLoading = $state(false);
	let error = $state('');
	let success = $state('');

	
	let titleError = $derived(title && title.length < 3 ? 'Title must be at least 3 characters' : '');
	let descriptionError = $derived(
		description && description.length < 10 ? 'Description must be at least 10 characters' : ''
	);
	let timeLimitError = $derived(timeLimitMs < 100 ? 'Time limit must be at least 100ms' : '');
	let memoryLimitError = $derived(memoryLimitMb < 1 ? 'Memory limit must be at least 1MB' : '');

	let isFormValid = $derived(
		title.length >= 3 &&
			description.length >= 10 &&
			timeLimitMs >= 100 &&
			memoryLimitMb >= 1 &&
			testCases.length > 0 &&
			testCases.every((tc) => tc.input_data.trim() && tc.output_data.trim()) &&
			!isLoading
	);

	function addTestCase() {
		testCases = [...testCases, { input_data: '', output_data: '', is_active: true }];
	}

	function removeTestCase(index: number) {
		if (testCases.length > 1) {
			testCases = testCases.filter((_, i) => i !== index);
		}
	}

	function updateTestCase(index: number, field: 'input_data' | 'output_data', value: string) {
		testCases = testCases.map((tc, i) => (i === index ? { ...tc, [field]: value } : tc));
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();

		if (!isFormValid) return;

		if (!authStore.token) {
			error = 'You must be logged in to create problems';
			return;
		}

		isLoading = true;
		error = '';
		success = '';

		try {
			const request: ProblemCreateRequest = {
				title,
				description,
				difficulty,
				time_limit_ms: timeLimitMs,
				memory_limit_mb: memoryLimitMb,
				is_active: isActive,
				testcases: testCases.filter((tc) => tc.input_data.trim() && tc.output_data.trim())
			};

			const response = await adminApi.createProblem(request, authStore.token);
			success = `Problem "${title}" created successfully! ID: ${response.problem.id}`;

			title = '';
			description = '';
			difficulty = 'easy';
			timeLimitMs = 3000;
			memoryLimitMb = 256;
			isActive = true;
			testCases = [{ input_data: '', output_data: '', is_active: true }];

			setTimeout(() => {
				goto('/admin/problems');
			}, 2000);
		} catch (err: any) {
			error = err.message || 'Failed to create problem';
		} finally {
			isLoading = false;
		}
	}

	function clearMessages() {
		error = '';
		success = '';
	}
</script>

<svelte:head>
	<title>Create Problem - Admin - Adel Online Judge</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<AdminNavbar currentPage="problems" />

	<main class="mx-auto max-w-4xl px-4 py-8 sm:px-6 lg:px-8">
		<div class="mb-8">
			<div class="flex items-center space-x-4">
				<a
					href="/admin/problems"
					class="text-gray-600 hover:text-gray-900"
					aria-label="Back to Problems"
				>
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M10 19l-7-7m0 0l7-7m-7 7h18"
						/>
					</svg>
				</a>
				<div>
					<h1 class="text-3xl font-bold text-gray-900">Create New Problem</h1>
					<p class="mt-2 text-gray-600">Add a new problem to the platform</p>
				</div>
			</div>
		</div>

		{#if success}
			<div class="mb-6">
				<Alert type="success" message={success} dismissible ondismiss={clearMessages} />
			</div>
		{/if}

		{#if error}
			<div class="mb-6">
				<Alert type="error" message={error} dismissible ondismiss={clearMessages} />
			</div>
		{/if}

		<div class="rounded-lg bg-white shadow-sm">
			<form onsubmit={handleSubmit} class="space-y-6 p-6">
				<!-- Basic Information -->
				<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
					<div class="lg:col-span-2">
						<Input
							label="Problem Title"
							type="text"
							placeholder="e.g., Find Maximum Element"
							required
							bind:value={title}
							error={titleError}
						/>
					</div>

					<div>
						<label for="difficulty" class="mb-1 block text-sm font-medium text-gray-700">
							Difficulty
						</label>
						<select
							id="difficulty"
							bind:value={difficulty}
							class="block w-full rounded-md border-gray-300 py-2 pl-3 pr-10 text-base focus:border-blue-500 focus:outline-none focus:ring-blue-500 sm:text-sm"
						>
							<option value="easy">Easy</option>
							<option value="medium">Medium</option>
							<option value="hard">Hard</option>
						</select>
					</div>

					<div class="flex items-center">
						<input
							id="is-active"
							type="checkbox"
							bind:checked={isActive}
							class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
						/>
						<label for="is-active" class="ml-2 block text-sm text-gray-900">
							Active (visible to users)
						</label>
					</div>
				</div>

				<!-- Description -->
				<div>
					<label for="description" class="mb-1 block text-sm font-medium text-gray-700">
						Problem Description
					</label>
					<textarea
						id="description"
						bind:value={description}
						placeholder="Describe the problem, input format, output format, and constraints..."
						rows="10"
						class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
					></textarea>
					{#if descriptionError}
						<p class="mt-1 text-sm text-red-600">{descriptionError}</p>
					{/if}
				</div>

				<!-- Limits -->
				<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
					<div>
						<label class="mb-1 block text-sm font-medium text-gray-700" for="time-limit">
							Time Limit (ms)
						</label>
						<input
							type="number"
							min="100"
							max="10000"
							bind:value={timeLimitMs}
							class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
							required
						/>
						{#if timeLimitError}
							<p class="mt-1 text-sm text-red-600">{timeLimitError}</p>
						{/if}
					</div>

					<div>
						<label class="mb-1 block text-sm font-medium text-gray-700" for="memory-limit">
							Memory Limit (MB)
						</label>
						<input
							type="number"
							min="1"
							max="1024"
							bind:value={memoryLimitMb}
							class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
							required
						/>
						{#if memoryLimitError}
							<p class="mt-1 text-sm text-red-600">{memoryLimitError}</p>
						{/if}
					</div>
				</div>

				<!-- Test Cases -->
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<h3 class="text-lg font-medium text-gray-900">Test Cases</h3>
						<Button type="button" variant="secondary" size="sm" onclick={addTestCase}>
							Add Test Case
						</Button>
					</div>

					{#each testCases as testCase, index}
						<div class="rounded-md border border-gray-200 p-4">
							<div class="mb-3 flex items-center justify-between">
								<h4 class="text-sm font-medium text-gray-700">Test Case {index + 1}</h4>
								{#if testCases.length > 1}
									<Button
										type="button"
										variant="ghost"
										size="sm"
										onclick={() => removeTestCase(index)}
									>
										Remove
									</Button>
								{/if}
							</div>

							<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
								<div>
									<label class="mb-1 block text-sm font-medium text-gray-700" for="input-data">
										Input
									</label>
									<textarea
										bind:value={testCase.input_data}
										oninput={(e) => {
											const target = e.target as HTMLTextAreaElement;
											updateTestCase(index, 'input_data', target.value);
										}}
										placeholder="Input data for this test case..."
										rows="4"
										class="block w-full rounded-md border-gray-300 font-mono shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
									></textarea>
								</div>

								<div>
									<label class="mb-1 block text-sm font-medium text-gray-700" for="output-data">
										Expected Output
									</label>
									<textarea
										bind:value={testCase.output_data}
										oninput={(e) => {
											const target = e.target as HTMLTextAreaElement;
											updateTestCase(index, 'output_data', target.value);
										}}
										placeholder="Expected output for this test case..."
										rows="4"
										class="block w-full rounded-md border-gray-300 font-mono shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
									></textarea>
								</div>
							</div>
						</div>
					{/each}
				</div>

				<!-- Submit Button -->
				<div class="flex justify-end space-x-3">
					<a href="/admin/problems">
						<Button type="button" variant="ghost">Cancel</Button>
					</a>
					<Button type="submit" variant="primary" disabled={!isFormValid} loading={isLoading}>
						{isLoading ? 'Creating...' : 'Create Problem'}
					</Button>
				</div>
			</form>
		</div>
	</main>
</div>
