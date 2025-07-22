<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte.js';
	import Button from './Button.svelte';

	interface Props {
		currentPage?: string;
	}

	let { currentPage = '' }: Props = $props();
	let mobileMenuOpen = $state(false);

	function getLinkClass(page: string): string {
		const baseClass = 'rounded-md px-3 py-2 text-sm font-medium transition-colors';
		if (currentPage === page) {
			return `${baseClass} bg-red-100 text-red-800`;
		}
		return `${baseClass} text-gray-600 hover:text-gray-900 hover:bg-gray-100`;
	}

	function getMobileLinkClass(page: string): string {
		const baseClass = 'block rounded-md px-3 py-2 text-base font-medium transition-colors';
		if (currentPage === page) {
			return `${baseClass} bg-red-100 text-red-800`;
		}
		return `${baseClass} text-gray-600 hover:text-gray-900 hover:bg-gray-100`;
	}

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}
</script>

<nav class="border-b border-gray-200 bg-white shadow-sm">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 justify-between">
			<div class="flex items-center">
				<div class="flex flex-shrink-0 items-center">
					<a href="/" class="flex items-center">
						<span class="text-xl font-bold text-gray-900">Adel OJ</span>
						<span class="ml-2 rounded-full bg-red-100 px-2 py-1 text-xs font-medium text-red-800"
							>Admin</span
						>
					</a>
				</div>

				{#if authStore.isAuthenticated}
					<!-- Desktop navigation -->
					<div class="ml-10 hidden md:block">
						<div class="flex items-baseline space-x-4">
							<a href="/admin" class={getLinkClass('dashboard')}>Dashboard</a>
							<a href="/admin/problems" class={getLinkClass('problems')}>Manage Problems</a>
							<a href="/admin/users" class={getLinkClass('users')}>Manage Users</a>
							<a href="/submissions" class={getLinkClass('submissions')}>Submissions</a>
						</div>
					</div>
				{/if}
			</div>

			<div class="flex items-center space-x-4">
				{#if authStore.isAuthenticated}
					<!-- Desktop user info -->
					<div class="hidden md:flex md:items-center md:space-x-4">
						<span class="text-sm text-gray-700">
							Admin: <span class="font-semibold">{authStore.user?.username}</span>
						</span>
						<Button variant="ghost" size="sm" onclick={() => authStore.clearAuth()}>Logout</Button>
					</div>

					<!-- Mobile menu button -->
					<div class="md:hidden">
						<button
							type="button"
							onclick={toggleMobileMenu}
							class="inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
							aria-expanded="false"
						>
							<span class="sr-only">Open main menu</span>
							{#if mobileMenuOpen}
								<!-- X icon -->
								<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M6 18L18 6M6 6l12 12"
									/>
								</svg>
							{:else}
								<!-- Hamburger icon -->
								<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M4 6h16M4 12h16M4 18h16"
									/>
								</svg>
							{/if}
						</button>
					</div>
				{/if}
			</div>
		</div>
	</div>

	<!-- Mobile menu -->
	{#if mobileMenuOpen && authStore.isAuthenticated}
		<div class="md:hidden">
			<div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
				<a href="/admin" class={getMobileLinkClass('dashboard')} onclick={closeMobileMenu}
					>Dashboard</a
				>
				<a href="/admin/problems" class={getMobileLinkClass('problems')} onclick={closeMobileMenu}
					>Manage Problems</a
				>
				<a href="/admin/users" class={getMobileLinkClass('users')} onclick={closeMobileMenu}
					>Manage Users</a
				>
				<a href="/submissions" class={getMobileLinkClass('submissions')} onclick={closeMobileMenu}
					>Submissions</a
				>
			</div>
			<div class="border-t border-gray-200 pb-3 pt-4">
				<div class="flex items-center px-5">
					<div class="flex-shrink-0">
						<div class="flex h-8 w-8 items-center justify-center rounded-full bg-red-100">
							<span class="text-sm font-medium text-red-800">
								{authStore.user?.username?.charAt(0).toUpperCase()}
							</span>
						</div>
					</div>
					<div class="ml-3">
						<div class="text-base font-medium text-gray-800">{authStore.user?.username}</div>
						<div class="text-sm text-gray-500">Admin</div>
					</div>
				</div>
				<div class="mt-3 space-y-1 px-2">
					<button
						type="button"
						onclick={() => {
							authStore.clearAuth();
							closeMobileMenu();
						}}
						class="block w-full rounded-md px-3 py-2 text-left text-base font-medium text-gray-600 hover:bg-gray-100 hover:text-gray-900"
					>
						Logout
					</button>
				</div>
			</div>
		</div>
	{/if}
</nav>
