<script lang="ts">
	import { goto } from '$app/navigation';
	import { authApi, type AuthRequest } from '$lib/api.js';
	import AuthCard from '$lib/components/AuthCard.svelte';
	import Input from '$lib/components/Input.svelte';
	import Button from '$lib/components/Button.svelte';
	import Alert from '$lib/components/Alert.svelte';

	let username = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let isLoading = $state(false);
	let error = $state('');
	let success = $state('');

	let usernameError = $derived(
		username && username.length < 3 ? 'Username must be at least 3 characters' : ''
	);

	let passwordError = $derived(
		password && password.length < 8 ? 'Password must be at least 8 characters' : ''
	);

	let confirmPasswordError = $derived(
		confirmPassword && password !== confirmPassword ? 'Passwords do not match' : ''
	);

	let isFormValid = $derived(
		username.length >= 3 && password.length >= 6 && password === confirmPassword && !isLoading
	);

	async function handleSubmit(e: Event) {
		e.preventDefault();

		if (!isFormValid) return;

		isLoading = true;
		error = '';
		success = '';

		try {
			const credentials: AuthRequest = { username, password };
			const response = await authApi.signup(credentials);

			success = `Welcome ${response.user.username}! Your account has been created successfully.`;

			setTimeout(() => {
				goto('/login');
			}, 2000);
		} catch (err: any) {
			error = err.message || 'Registration failed. Please try again.';
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
	<title>Register - Adel Online Judge</title>
</svelte:head>

<AuthCard title="Create your account">
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

		{#if success}
			<Alert type="success" message={success} />
		{/if}

		<Input
			label="Username"
			type="text"
			placeholder="Enter your username"
			required
			bind:value={username}
			error={usernameError}
			help="Must be at least 3 characters long"
		/>

		<Input
			label="Password"
			type="password"
			placeholder="Enter your password"
			required
			bind:value={password}
			error={passwordError}
			help="Must be at least 8 characters long"
		/>

		<Input
			label="Confirm Password"
			type="password"
			placeholder="Confirm your password"
			required
			bind:value={confirmPassword}
			error={confirmPasswordError}
		/>

		<div class="flex justify-end">
			<Button type="submit" variant="primary" size="lg" disabled={!isFormValid} loading={isLoading}>
				{isLoading ? 'Creating Account...' : 'Create Account'}
			</Button>
		</div>

		<div class="text-center">
			<p class="text-sm text-gray-600">
				Already have an account?
				<a href="/login" class="font-medium text-blue-600 transition-colors hover:text-blue-500">
					Sign in here
				</a>
			</p>
		</div>
	</form>
</AuthCard>
