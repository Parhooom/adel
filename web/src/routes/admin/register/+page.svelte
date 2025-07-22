<script lang="ts">
	import { goto } from '$app/navigation';
	import { authApi, type AuthRequest } from '$lib/api.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import AdminNavbar from '$lib/components/AdminNavbar.svelte';
	import Input from '$lib/components/Input.svelte';
	import Button from '$lib/components/Button.svelte';
	import Alert from '$lib/components/Alert.svelte';

	$effect(() => {
		if (authStore.isAuthenticated) {
			goto('/admin');
		}
	});

	let username = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let isLoading = $state(false);
	let error = $state('');

	let usernameError = $derived(
		username && username.length < 3 ? 'Username must be at least 3 characters' : ''
	);
	let passwordError = $derived(
		password && password.length < 6 ? 'Password must be at least 6 characters' : ''
	);
	let confirmPasswordError = $derived(
		confirmPassword && confirmPassword !== password ? 'Passwords do not match' : ''
	);

	let isFormValid = $derived(
		username.length >= 3 && password.length >= 6 && confirmPassword === password && !isLoading
	);

	async function handleSubmit(e: Event) {
		e.preventDefault();

		if (!isFormValid) return;

		isLoading = true;
		error = '';

		try {
			const request: AuthRequest = {
				username,
				password
			};

			const response = await authApi.adminSignup(request);

			const loginResponse = await authApi.login({ username, password });
			authStore.setAuth(response.user, loginResponse.token.token);
			goto('/admin');
		} catch (err: any) {
			error = err.message || 'Registration failed';
		} finally {
			isLoading = false;
		}
	}

	function clearError() {
		error = '';
	}
</script>

<svelte:head>
	<title>Admin Register - Adel Online Judge</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<AdminNavbar />

	<main class="flex min-h-screen items-center justify-center px-4 py-12 sm:px-6 lg:px-8">
		<div class="w-full max-w-md space-y-8">
			<div class="text-center">
				<h2 class="mt-6 text-3xl font-bold tracking-tight text-gray-900">Create Admin Account</h2>
				<p class="mt-2 text-sm text-gray-600">
					Already have an account?
					<a href="/login" class="font-medium text-blue-600 hover:text-blue-500"> Sign in here </a>
				</p>
			</div>

			<div class="rounded-lg border border-gray-100 bg-white px-6 py-8 shadow-lg">
				{#if error}
					<div class="mb-6">
						<Alert type="error" message={error} dismissible ondismiss={clearError} />
					</div>
				{/if}

				<form class="space-y-6" onsubmit={handleSubmit}>
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

					<Input
						label="Confirm Password"
						type="password"
						placeholder="Confirm your password"
						required
						bind:value={confirmPassword}
						error={confirmPasswordError}
					/>

					<div class="w-full">
						<Button
							type="submit"
							variant="primary"
							size="lg"
							disabled={!isFormValid}
							loading={isLoading}
						>
							{isLoading ? 'Creating Account...' : 'Create Account'}
						</Button>
					</div>
				</form>
			</div>
		</div>
	</main>
</div>
