import { browser } from '$app/environment';
import type { User } from '../api.js';

interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}

function createAuthStore() {
	const initialToken = browser ? localStorage.getItem('auth_token') : null;
	const initialUser = browser ? JSON.parse(localStorage.getItem('auth_user') || 'null') : null;

	let state = $state<AuthState>({
		user: initialUser,
		token: initialToken,
		isAuthenticated: !!initialToken,
		isLoading: false
	});

	return {
		get user() {
			return state.user;
		},
		get token() {
			return state.token;
		},
		get isAuthenticated() {
			return state.isAuthenticated;
		},
		get isAdmin() {
			return state.user?.is_admin || false;
		},
		get isLoading() {
			return state.isLoading;
		},

		setLoading(loading: boolean) {
			state.isLoading = loading;
		},

		setAuth(user: User, token: string) {
			state.user = user;
			state.token = token;
			state.isAuthenticated = true;

			if (browser) {
				localStorage.setItem('auth_token', token);
				localStorage.setItem('auth_user', JSON.stringify(user));
			}
		},

		clearAuth() {
			state.user = null;
			state.token = null;
			state.isAuthenticated = false;

			if (browser) {
				localStorage.removeItem('auth_token');
				localStorage.removeItem('auth_user');
			}
		}
	};
}

export const authStore = createAuthStore();
