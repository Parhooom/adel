<script lang="ts">
	import { goto } from '$app/navigation';
	import { authApi, type AuthRequest } from '$lib/api.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import AuthCard from '$lib/components/AuthCard.svelte';
	import Input from '$lib/components/Input.svelte';
	import Button from '$lib/components/Button.svelte';
	import Alert from '$lib/components/Alert.svelte';

	let username = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let error = $state('');

	let usernameError = $derived(
		username && username.length < 3 ? 'Username must be at least 3 characters' : ''
	);

	let passwordError = $derived(
		password && password.length < 6 ? 'Password must be at least 6 characters' : ''
	);

	let isFormValid = $derived(username.length >= 3 && password.length >= 6 && !isLoading);

	async function handleSubmit(e: Event) {
		e.preventDefault();

		if (!isFormValid) return;

		isLoading = true;
		error = '';

		try {
			const credentials: AuthRequest = { username, password };
			const loginResponse = await authApi.login(credentials);
			const token = loginResponse.token.token;

			const userResponse = await authApi.getCurrentUser(token);

			authStore.setAuth(userResponse.user, token);

			goto('/');
		} catch (err: any) {
			error = err.message || 'Login failed. Please check your credentials.';
		} finally {
			isLoading = false;
		}
	}

	function clearMessages() {
		error = '';
	}
</script>

<svelte:head>
	<title>Sign in - Adel Online Judge</title>
</svelte:head>

<AuthCard title="Sign in to your account">
	<div class="mb-4">
		<a
			href="/"
			class="inline-flex items-center text-sm text-gray-600 transition-colors hover:text-gray-900"
		>
			<svg class="mr-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M10 19l-7-7m0 0l7-7m-7 7h18"
				/>
			</svg>
			Back to Home
		</a>
	</div>
	<form onsubmit={handleSubmit} class="space-y-6">
		{#if error}
			<Alert type="error" message={error} dismissible ondismiss={clearMessages} />
		{/if}

		<Input
			label="Username"
			type="text"
			placeholder="Enter your username"
			required
			bind:value={username}
			error={usernameError}
		/>

		<Input
			label="Password"
			type="password"
			placeholder="Enter your password"
			required
			bind:value={password}
			error={passwordError}
		/>

		<div class="flex justify-end">
			<Button type="submit" variant="primary" size="lg" disabled={!isFormValid} loading={isLoading}>
				{isLoading ? 'Signing in...' : 'Sign In'}
			</Button>
		</div>

		<div class="text-center">
			<p class="text-sm text-gray-600">
				Don't have an account?
				<a href="/register" class="font-medium text-blue-600 transition-colors hover:text-blue-500">
					Create one here
				</a>
			</p>
		</div>
	</form>
</AuthCard>
