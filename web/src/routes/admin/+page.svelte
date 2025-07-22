<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { adminApi, type AdminStats } from '$lib/api.js';
	import AdminNavbar from '$lib/components/AdminNavbar.svelte';
	import Button from '$lib/components/Button.svelte';

	$effect(() => {
		if (!authStore.isAuthenticated) {
			goto('/login');
		} else if (!authStore.isAdmin) {
			goto('/');
		}
	});

	let stats = $state<AdminStats>({
		total_problems: 0,
		registered_users: 0,
		total_submissions: 0,
		active_problems: 0
	});
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		await loadStats();
	});

	async function loadStats() {
		try {
			loading = true;
			if (!authStore.token) {
				error = 'You must be logged in to view admin stats';
				return;
			}
			const response = await adminApi.getAdminStats(authStore.token);
			stats = response.stats;
		} catch (err: any) {
			error = err.message || 'Failed to load admin statistics';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Admin Dashboard - Adel Online Judge</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<AdminNavbar currentPage="dashboard" />

	<main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
			<p class="mt-2 text-gray-600">Manage your online judge platform</p>
		</div>

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

		<!-- Stats Grid -->
		<div class="mb-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
			<div class="overflow-hidden rounded-lg bg-white shadow">
				<div class="p-6">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<svg
								class="h-8 w-8 text-blue-600"
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
								<dt class="truncate text-sm font-medium text-gray-500">Total Problems</dt>
								<dd class="text-3xl font-semibold text-gray-900">{stats.total_problems}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>

			<div class="overflow-hidden rounded-lg bg-white shadow">
				<div class="p-6">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<svg
								class="h-8 w-8 text-green-600"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
								/>
							</svg>
						</div>
						<div class="ml-5 w-0 flex-1">
							<dl>
								<dt class="truncate text-sm font-medium text-gray-500">Registered Users</dt>
								<dd class="text-3xl font-semibold text-gray-900">{stats.registered_users}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>

			<div class="overflow-hidden rounded-lg bg-white shadow">
				<div class="p-6">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<svg
								class="h-8 w-8 text-purple-600"
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
								<dt class="truncate text-sm font-medium text-gray-500">Total Submissions</dt>
								<dd class="text-3xl font-semibold text-gray-900">{stats.total_submissions}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>

			<div class="overflow-hidden rounded-lg bg-white shadow">
				<div class="p-6">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<svg
								class="h-8 w-8 text-yellow-600"
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
								<dt class="truncate text-sm font-medium text-gray-500">Active Problems</dt>
								<dd class="text-3xl font-semibold text-gray-900">{stats.active_problems}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Quick Actions -->
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<div class="rounded-lg bg-white p-6 shadow">
				<h2 class="mb-4 text-lg font-semibold text-gray-900">Quick Actions</h2>
				<div class="space-y-3">
					<a href="/admin/problems/create" class="block">
						<Button variant="primary" size="sm">Create New Problem</Button>
					</a>
					<a href="/admin/problems" class="block">
						<Button variant="secondary" size="sm">Manage Problems</Button>
					</a>
					<a href="/admin/users" class="block">
						<Button variant="secondary" size="sm">Manage Users</Button>
					</a>
				</div>
			</div>
		</div>
	</main>
</div>
