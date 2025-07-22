<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { adminApi, type User } from '$lib/api.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import AdminNavbar from '$lib/components/AdminNavbar.svelte';
	import Table from '$lib/components/Table.svelte';
	import Badge from '$lib/components/Badge.svelte';
	import Button from '$lib/components/Button.svelte';
	import Alert from '$lib/components/Alert.svelte';

	let users = $state<User[]>([]);
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
		await loadUsers();
	});

	async function loadUsers() {
		try {
			loading = true;
			if (!authStore.token) {
				error = 'You must be logged in to view users';
				return;
			}
			const response = await adminApi.getAllUsers(authStore.token);
			users = response.users;
		} catch (err: any) {
			error = err.message || 'Failed to load users';
		} finally {
			loading = false;
		}
	}

	async function deleteUser(id: number, username: string) {
		if (
			!confirm(`Are you sure you want to delete user "${username}"? This action cannot be undone.`)
		) {
			return;
		}

		if (!authStore.token) {
			error = 'You must be logged in to delete users';
			return;
		}

		try {
			deletingId = id;
			await adminApi.deleteUser(id, authStore.token);
			success = `User "${username}" deleted successfully`;
			
			users = users.filter((u) => u.id !== id);
		} catch (err: any) {
			error = err.message || 'Failed to delete user';
		} finally {
			deletingId = null;
		}
	}

	function clearMessages() {
		error = '';
		success = '';
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getRoleVariant(isAdmin: boolean): 'error' | 'info' {
		return isAdmin ? 'error' : 'info';
	}
</script>

<svelte:head>
	<title>Manage Users - Admin - Adel Online Judge</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<AdminNavbar currentPage="users" />

	<main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
		<div class="mb-8 flex items-center justify-between">
			<div>
				<h1 class="text-3xl font-bold text-gray-900">Manage Users</h1>
				<p class="mt-2 text-gray-600">View and manage registered users</p>
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
			<Table
				headers={['ID', 'Username', 'Role', 'Created At', 'Actions']}
				{loading}
				empty={users.length === 0}
				emptyMessage="No users found."
			>
				{#each users as user}
					<tr class="hover:bg-gray-50">
						<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
							#{user.id}
						</td>
						<td class="whitespace-nowrap px-6 py-4">
							<span class="text-sm font-medium text-gray-900">
								{user.username}
							</span>
						</td>
						<td class="whitespace-nowrap px-6 py-4">
							<Badge variant={getRoleVariant(user.is_admin)}>
								{user.is_admin ? 'Admin' : 'User'}
							</Badge>
						</td>
						<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
							{formatDate(user.created_at)}
						</td>
						<td class="whitespace-nowrap px-6 py-4 text-sm font-medium">
							<div class="flex space-x-2">
								{#if user.id !== authStore.user?.id}
									<Button
										variant="danger"
										size="sm"
										onclick={() => deleteUser(user.id, user.username)}
										disabled={deletingId === user.id}
										loading={deletingId === user.id}
									>
										{deletingId === user.id ? 'Deleting...' : 'Delete'}
									</Button>
								{:else}
									<span class="text-xs text-gray-400">Current User</span>
								{/if}
							</div>
						</td>
					</tr>
				{/each}
			</Table>
		</div>

		{#if users.length > 0}
			<div class="mt-6 text-center">
				<p class="text-sm text-gray-500">
					Showing {users.length} user{users.length === 1 ? '' : 's'}
				</p>
			</div>
		{/if}
	</main>
</div>
