const API_BASE_URL = 'http://localhost:8080';

export interface User {
	id: number;
	username: string;
	is_admin: boolean;
	created_at: string;
	updated_at: string;
}

export interface AuthRequest {
	username: string;
	password: string;
}

export interface SignupResponse {
	user: User;
}

export interface TokenData {
	token: string;
	expiry: string;
}

export interface LoginResponse {
	token: TokenData;
}

export interface TestCase {
	id: number;
	user_id: number;
	problem_id: number;
	is_active: boolean;
	input_data: string;
	output_data: string;
}

export interface ProblemSummary {
	id: number;
	title: string;
	description: string;
	difficulty: string;
	is_active: boolean;
}

export interface Problem {
	id: number;
	user_id: number;
	title: string;
	description: string;
	difficulty: string;
	time_limit_ms: number;
	memory_limit_mb: number;
	is_active: boolean;
	testcases: TestCase[];
}

export interface ProblemsResponse {
	problems: ProblemSummary[];
}

export interface ProblemResponse {
	problem: Problem;
}

export interface SubmissionRequest {
	problem_id: number;
	code: string;
	language: string;
}

export interface Submission {
	id: number;
	user_id: number;
	problem_id: number;
	code: string;
	language: string;
	status: string;
	execution_time_ms: number;
	memory_usage_mb: number;
	error_message: string;
	created_at: string;
	updated_at: string;
}

export interface SubmissionResponse {
	submission: Submission;
}

export interface SubmissionsResponse {
	submissions: Submission[];
}

export interface UsersResponse {
	users: User[];
}

export interface AdminStats {
	total_problems: number;
	registered_users: number;
	total_submissions: number;
	active_problems: number;
}

export interface AdminStatsResponse {
	stats: AdminStats;
}

export interface TestCaseRequest {
	input_data: string;
	output_data: string;
	is_active: boolean;
}

export interface ProblemCreateRequest {
	title: string;
	description: string;
	difficulty: string;
	time_limit_ms: number;
	memory_limit_mb: number;
	is_active: boolean;
	testcases: TestCaseRequest[];
}

export interface ProblemUpdateRequest extends ProblemCreateRequest {}

export interface ProblemCreateResponse {
	problem: Problem;
}

class ApiError extends Error {
	constructor(
		message: string,
		public status: number
	) {
		super(message);
		this.name = 'ApiError';
	}
}

async function apiRequest<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
	const url = `${API_BASE_URL}${endpoint}`;

	const response = await fetch(url, {
		headers: {
			'Content-Type': 'application/json',
			...options.headers
		},
		...options
	});

	if (!response.ok) {
		const errorText = await response.text();
		throw new ApiError(errorText || 'An error occurred', response.status);
	}

	if (response.status === 204) {
		return undefined as T;
	}

	const contentType = response.headers.get('content-type');
	if (contentType && contentType.includes('application/json')) {
		return response.json();
	}

	return undefined as T;
}

function getAuthHeaders(token?: string): Record<string, string> {
	if (!token) return {};
	return {
		Authorization: `Bearer ${token}`
	};
}

export const authApi = {
	async signup(credentials: AuthRequest): Promise<SignupResponse> {
		return apiRequest<SignupResponse>('/users/register', {
			method: 'POST',
			body: JSON.stringify(credentials)
		});
	},

	async adminSignup(credentials: AuthRequest): Promise<SignupResponse> {
		return apiRequest<SignupResponse>('/users/register-admin', {
			method: 'POST',
			body: JSON.stringify(credentials)
		});
	},

	async login(credentials: AuthRequest): Promise<LoginResponse> {
		return apiRequest<LoginResponse>('/users/login', {
			method: 'POST',
			body: JSON.stringify(credentials)
		});
	},

	async getCurrentUser(token: string): Promise<{ user: User }> {
		return apiRequest<{ user: User }>('/users/me', {
			headers: getAuthHeaders(token)
		});
	},

	async getUserStats(
		token: string
	): Promise<{ stats: { problems_solved: number; success_rate: number } }> {
		return apiRequest<{ stats: { problems_solved: number; success_rate: number } }>(
			'/users/stats',
			{
				headers: getAuthHeaders(token)
			}
		);
	}
};

export const problemsApi = {
	async getAllProblems(): Promise<ProblemsResponse> {
		return apiRequest<ProblemsResponse>('/problems');
	},

	async getProblemById(id: number): Promise<ProblemResponse> {
		return apiRequest<ProblemResponse>(`/problems/${id}`);
	}
};

export const submissionsApi = {
	async submitSolution(request: SubmissionRequest, token: string): Promise<SubmissionResponse> {
		return apiRequest<SubmissionResponse>('/submissions', {
			method: 'POST',
			headers: getAuthHeaders(token),
			body: JSON.stringify(request)
		});
	},

	async getUserSubmissions(token: string): Promise<SubmissionsResponse> {
		return apiRequest<SubmissionsResponse>('/submissions/user', {
			headers: getAuthHeaders(token)
		});
	},

	async getSubmissionById(id: number, token: string): Promise<SubmissionResponse> {
		return apiRequest<SubmissionResponse>(`/submissions/${id}`, {
			headers: getAuthHeaders(token)
		});
	}
};

export const adminApi = {
	async createProblem(
		request: ProblemCreateRequest,
		token: string
	): Promise<ProblemCreateResponse> {
		return apiRequest<ProblemCreateResponse>('/problems', {
			method: 'POST',
			headers: getAuthHeaders(token),
			body: JSON.stringify(request)
		});
	},

	async updateProblem(
		id: number,
		request: ProblemUpdateRequest,
		token: string
	): Promise<ProblemCreateResponse> {
		return apiRequest<ProblemCreateResponse>(`/problems/${id}`, {
			method: 'PUT',
			headers: getAuthHeaders(token),
			body: JSON.stringify(request)
		});
	},

	async deleteProblem(id: number, token: string): Promise<void> {
		return apiRequest<void>(`/problems/${id}`, {
			method: 'DELETE',
			headers: getAuthHeaders(token)
		});
	},

	async getAllUsers(token: string): Promise<UsersResponse> {
		return apiRequest<UsersResponse>('/users', {
			headers: getAuthHeaders(token)
		});
	},

	async deleteUser(id: number, token: string): Promise<void> {
		return apiRequest<void>(`/users/${id}`, {
			method: 'DELETE',
			headers: getAuthHeaders(token)
		});
	},

	async getAdminStats(token: string): Promise<AdminStatsResponse> {
		return apiRequest<AdminStatsResponse>('/admin/stats', {
			headers: getAuthHeaders(token)
		});
	},

	async getAllProblemsForAdmin(token: string): Promise<ProblemsResponse> {
		return apiRequest<ProblemsResponse>('/admin/problems', {
			headers: getAuthHeaders(token)
		});
	}
};
