<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { authApi } from '$lib/api.js';
	import { PageLayout, Button } from '$lib/components/index.js';

	let problemsSolved = $state(0);
	let successRate = $state(0);
	let loading = $state(true);

	$effect(() => {
		if (!authStore.isAuthenticated) {
			goto('/login');
		}
	});

	async function fetchUserStats() {
		if (!authStore.token) return;

		try {
			loading = true;
			const response = await authApi.getUserStats(authStore.token);
			problemsSolved = response.stats.problems_solved;
			successRate = Math.round(response.stats.success_rate * 100) / 100;
		} catch (error) {
			console.error('Error fetching user stats:', error);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		if (authStore.isAuthenticated) {
			fetchUserStats();
		}
	});

	function handleLogout() {
		authStore.clearAuth();
		goto('/');
	}
</script>

<PageLayout
	currentPage="dashboard"
	title="Dashboard"
	subtitle="Welcome back to Adel Online Judge, {authStore.user?.username}!"
>
	<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
		<div class="overflow-hidden rounded-lg bg-white shadow">
			<div class="p-5">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<svg
							class="h-6 w-6 text-blue-600"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
							/>
						</svg>
					</div>
					<div class="ml-5 w-0 flex-1">
						<dl>
							<dt class="truncate text-sm font-medium text-gray-500">Problems Solved</dt>
							<dd class="text-lg font-medium text-gray-900">
								{loading ? '...' : problemsSolved}
							</dd>
						</dl>
					</div>
				</div>
			</div>
			<div class="bg-gray-50 px-5 py-3">
				<div class="text-sm">
					<a href="/problems" class="font-medium text-blue-600 hover:text-blue-500">View problems</a
					>
				</div>
			</div>
		</div>

		<div class="overflow-hidden rounded-lg bg-white shadow">
			<div class="p-5">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<svg
							class="h-6 w-6 text-green-600"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
							/>
						</svg>
					</div>
					<div class="ml-5 w-0 flex-1">
						<dl>
							<dt class="truncate text-sm font-medium text-gray-500">Success Rate</dt>
							<dd class="text-lg font-medium text-gray-900">
								{loading ? '...' : `${successRate}%`}
							</dd>
						</dl>
					</div>
				</div>
			</div>
			<div class="bg-gray-50 px-5 py-3">
				<div class="text-sm">
					<a href="/submissions" class="font-medium text-green-600 hover:text-green-500"
						>View submissions</a
					>
				</div>
			</div>
		</div>

		<div class="overflow-hidden rounded-lg bg-white shadow">
			<div class="p-5">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<svg
							class="h-6 w-6 text-purple-600"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M13 10V3L4 14h7v7l9-11h-7z"
							/>
						</svg>
					</div>
					<div class="ml-5 w-0 flex-1">
						<dl>
							<dt class="truncate text-sm font-medium text-gray-500">Quick Start</dt>
							<dd class="text-sm text-gray-900">Ready to practice?</dd>
						</dl>
					</div>
				</div>
			</div>
			<div class="bg-gray-50 px-5 py-3">
				<div class="text-sm">
					<a href="/problems" class="font-medium text-purple-600 hover:text-purple-500"
						>Start coding</a
					>
				</div>
			</div>
		</div>
	</div>

	<div class="mt-8 rounded-lg bg-white p-6 shadow">
		<h2 class="mb-4 text-lg font-semibold text-gray-900">Getting Started</h2>
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<div class="flex items-start space-x-3">
				<div class="flex-shrink-0">
					<div class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100">
						<span class="text-sm font-medium text-blue-800">1</span>
					</div>
				</div>
				<div>
					<h3 class="text-sm font-medium text-gray-900">Browse Problems</h3>
					<p class="text-sm text-gray-500">
						Explore our collection of algorithmic challenges suited for all skill levels.
					</p>
				</div>
			</div>
			<div class="flex items-start space-x-3">
				<div class="flex-shrink-0">
					<div class="flex h-8 w-8 items-center justify-center rounded-full bg-green-100">
						<span class="text-sm font-medium text-green-800">2</span>
					</div>
				</div>
				<div>
					<h3 class="text-sm font-medium text-gray-900">Submit Solutions</h3>
					<p class="text-sm text-gray-500">
						Write your code in C, Python, or Go and submit for automated testing.
					</p>
				</div>
			</div>
			<div class="flex items-start space-x-3">
				<div class="flex-shrink-0">
					<div class="flex h-8 w-8 items-center justify-center rounded-full bg-yellow-100">
						<span class="text-sm font-medium text-yellow-800">3</span>
					</div>
				</div>
				<div>
					<h3 class="text-sm font-medium text-gray-900">Track Progress</h3>
					<p class="text-sm text-gray-500">
						Monitor your submission history and see your success rate improve.
					</p>
				</div>
			</div>
			<div class="flex items-start space-x-3">
				<div class="flex-shrink-0">
					<div class="flex h-8 w-8 items-center justify-center rounded-full bg-purple-100">
						<span class="text-sm font-medium text-purple-800">4</span>
					</div>
				</div>
				<div>
					<h3 class="text-sm font-medium text-gray-900">Keep Learning</h3>
					<p class="text-sm text-gray-500">
						Challenge yourself with increasingly difficult problems and improve your skills.
					</p>
				</div>
			</div>
		</div>
	</div>
</PageLayout>
