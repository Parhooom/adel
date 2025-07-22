<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte.js';
	import Button from '$lib/components/Button.svelte';
	import ResponsiveNavbar from '$lib/components/ResponsiveNavbar.svelte';

	function handleLogout() {
		authStore.clearAuth();
	}
</script>

<div class="min-h-screen bg-gray-50">
	<!-- Navigation Bar -->
	<ResponsiveNavbar />

	<!-- Main Content -->
	<main class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
		<!-- Landing page for all users -->
		<div class="text-center">
			<div class="mx-auto max-w-2xl">
				<h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
					Welcome to <span class="text-blue-600">Adel</span> Online Judge
				</h1>
				<p class="mt-6 text-lg leading-8 text-gray-600">
					Challenge yourself with programming problems, improve your coding skills, and compete with
					developers worldwide. Start your journey today!
				</p>
				<div class="mt-10 flex items-center justify-center gap-x-6">
					<a href="/problems">
						<Button variant="primary" size="lg">Explore Problems</Button>
					</a>
					{#if !authStore.isAuthenticated}
						<a href="/register">
							<Button variant="secondary" size="lg">Get Started</Button>
						</a>
					{:else}
						<a href="/dashboard">
							<Button variant="secondary" size="lg">Go to Dashboard</Button>
						</a>
					{/if}
				</div>
			</div>
		</div>

		<!-- Features Section -->
		<div class="mx-auto mt-32 max-w-7xl px-6 sm:mt-40 lg:px-8">
			<div class="mx-auto max-w-2xl lg:text-center">
				<p class="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
					Everything you need to excel in programming
				</p>
				<p class="mt-6 text-lg leading-8 text-gray-600">
					From beginner-friendly problems to advanced algorithmic challenges, our platform provides
					a comprehensive learning experience for programmers of all levels.
				</p>
			</div>
			<div class="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-4xl">
				<dl
					class="grid max-w-xl grid-cols-1 gap-x-8 gap-y-10 lg:max-w-none lg:grid-cols-2 lg:gap-y-16"
				>
					<div class="relative pl-16">
						<dt class="text-base font-semibold leading-7 text-gray-900">
							<div
								class="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-blue-600"
							>
								<svg
									class="h-6 w-6 text-white"
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
							Diverse Problem Set
						</dt>
						<dd class="mt-2 text-base leading-7 text-gray-600">
							Hundreds of carefully curated problems covering algorithms, data structures, and
							competitive programming concepts.
						</dd>
					</div>
					<div class="relative pl-16">
						<dt class="text-base font-semibold leading-7 text-gray-900">
							<div
								class="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-blue-600"
							>
								<svg
									class="h-6 w-6 text-white"
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
							Real-time Judging
						</dt>
						<dd class="mt-2 text-base leading-7 text-gray-600">
							Get instant feedback on your submissions with detailed execution results and
							performance metrics.
						</dd>
					</div>
					<div class="relative pl-16">
						<dt class="text-base font-semibold leading-7 text-gray-900">
							<div
								class="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-blue-600"
							>
								<svg
									class="h-6 w-6 text-white"
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
							Progress Tracking
						</dt>
						<dd class="mt-2 text-base leading-7 text-gray-600">
							Monitor your improvement with comprehensive statistics and achievement tracking.
						</dd>
					</div>
					<div class="relative pl-16">
						<dt class="text-base font-semibold leading-7 text-gray-900">
							<div
								class="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-blue-600"
							>
								<svg
									class="h-6 w-6 text-white"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
									/>
								</svg>
							</div>
							Global Community
						</dt>
						<dd class="mt-2 text-base leading-7 text-gray-600">
							Join a vibrant community of programmers, participate in contests, and climb the
							leaderboards.
						</dd>
					</div>
				</dl>
			</div>
		</div>
	</main>
</div>
