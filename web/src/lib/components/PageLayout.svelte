<script lang="ts">
	import ResponsiveNavbar from './ResponsiveNavbar.svelte';
	import AdminNavbar from './AdminNavbar.svelte';

	interface Props {
		currentPage?: string;
		title?: string;
		subtitle?: string;
		isAdmin?: boolean;
		children: any;
	}

	let { currentPage = '', title, subtitle, isAdmin = false, children }: Props = $props();
</script>

<svelte:head>
	{#if title}
		<title>{title} - Adel Online Judge</title>
	{/if}
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<!-- Navigation Bar -->
	{#if isAdmin}
		<AdminNavbar {currentPage} />
	{:else}
		<ResponsiveNavbar {currentPage} />
	{/if}

	<!-- Main Content -->
	<main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
		{#if title}
			<div class="mb-8">
				<h1 class="text-3xl font-bold text-gray-900">{title}</h1>
				{#if subtitle}
					<p class="mt-2 text-gray-600">{subtitle}</p>
				{/if}
			</div>
		{/if}

		{@render children()}
	</main>
</div>
