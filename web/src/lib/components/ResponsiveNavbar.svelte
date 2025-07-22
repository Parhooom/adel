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
			return `${baseClass} text-blue-600 hover:text-blue-800`;
		}
		return `${baseClass} text-gray-600 hover:text-gray-900`;
	}

	function getMobileLinkClass(page: string): string {
		const baseClass = 'block rounded-md px-3 py-2 text-base font-medium transition-colors';
		if (currentPage === page) {
			return `${baseClass} text-blue-600 hover:text-blue-800`;
		}
		return `${baseClass} text-gray-600 hover:text-gray-900 hover:bg-gray-100`;
	}

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}

	function handleLogout() {
		authStore.clearAuth();
		closeMobileMenu();
	}
</script>

<nav class="border-b border-gray-200 bg-white shadow-sm">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 justify-between">
			<div class="flex items-center">
				<div class="flex flex-shrink-0 items-center">
					<a href="/" class="flex items-center">
						<span class="ml-2 text-xl font-bold text-gray-900">Adel OJ</span>
					</a>
				</div>

				{#if authStore.isAuthenticated}
					<!-- Desktop navigation -->
					<div class="ml-10 hidden md:block">
						<div class="flex items-baseline space-x-4">
							<a href="/dashboard" class={getLinkClass('dashboard')}>Dashboard</a>
							{#if authStore.isAdmin}
								<a href="/admin" class={getLinkClass('admin')}>Admin Dashboard</a>
							{/if}
							<a href="/problems" class={getLinkClass('problems')}>Problems</a>
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
							Welcome, <span class="font-semibold">{authStore.user?.username}</span>!
						</span>
						<Button variant="ghost" size="sm" onclick={handleLogout}>Logout</Button>
					</div>
				{:else}
					<!-- Non-authenticated desktop actions -->
					<div class="hidden md:flex md:items-center md:space-x-4">
						<a href="/login">
							<Button variant="ghost" size="sm">Login</Button>
						</a>
						<a href="/register">
							<Button variant="primary" size="sm">Register</Button>
						</a>
					</div>
				{/if}

				<!-- Mobile menu button -->
				<div class="md:hidden">
					<button
						type="button"
						onclick={toggleMobileMenu}
						class="inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
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
			</div>
		</div>
	</div>

	<!-- Mobile menu -->
	{#if mobileMenuOpen}
		<div class="md:hidden">
			{#if authStore.isAuthenticated}
				<!-- Authenticated mobile menu -->
				<div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
					<a href="/dashboard" class={getMobileLinkClass('dashboard')} onclick={closeMobileMenu}
						>Dashboard</a
					>
					{#if authStore.isAdmin}
						<a href="/admin" class={getMobileLinkClass('admin')} onclick={closeMobileMenu}
							>Admin Dashboard</a
						>
					{/if}
					<a href="/problems" class={getMobileLinkClass('problems')} onclick={closeMobileMenu}
						>Problems</a
					>
					<a href="/submissions" class={getMobileLinkClass('submissions')} onclick={closeMobileMenu}
						>Submissions</a
					>
				</div>
				<div class="border-t border-gray-200 pb-3 pt-4">
					<div class="flex items-center px-5">
						<div class="flex-shrink-0">
							<div class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100">
								<span class="text-sm font-medium text-blue-800">
									{authStore.user?.username?.charAt(0).toUpperCase()}
								</span>
							</div>
						</div>
						<div class="ml-3">
							<div class="text-base font-medium text-gray-800">{authStore.user?.username}</div>
							<div class="text-sm text-gray-500">User</div>
						</div>
					</div>
					<div class="mt-3 space-y-1 px-2">
						<button
							type="button"
							onclick={handleLogout}
							class="block w-full rounded-md px-3 py-2 text-left text-base font-medium text-gray-600 hover:bg-gray-100 hover:text-gray-900"
						>
							Logout
						</button>
					</div>
				</div>
			{:else}
				<!-- Non-authenticated mobile menu -->
				<div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
					<a href="/problems" class={getMobileLinkClass('problems')} onclick={closeMobileMenu}
						>Problems</a
					>
				</div>
				<div class="border-t border-gray-200 pb-3 pt-4">
					<div class="space-y-1 px-2">
						<a
							href="/login"
							onclick={closeMobileMenu}
							class="block rounded-md px-3 py-2 text-base font-medium text-gray-600 hover:bg-gray-100 hover:text-gray-900"
						>
							Login
						</a>
						<a
							href="/register"
							onclick={closeMobileMenu}
							class="block rounded-md px-3 py-2 text-base font-medium text-blue-600 hover:bg-blue-100 hover:text-blue-800"
						>
							Register
						</a>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</nav>
