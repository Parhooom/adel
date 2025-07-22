<script lang="ts">
	interface Props {
		title: string;
		subtitle?: string;
		children: any;
	}

	let { title, subtitle, children }: Props = $props();
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md space-y-8">
		<div>
			<!-- <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-blue-100"> -->
			<!-- Adel OJ Logo/Icon -->
			<!-- <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
					/>
				</svg>
			</div> -->
			<h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">
				{title}
			</h2>
			{#if subtitle}
				<p class="mt-2 text-center text-sm text-gray-600">
					{subtitle}
				</p>
			{/if}
		</div>

		<div class="rounded-lg border border-gray-100 bg-white px-6 py-8 shadow-lg">
			{@render children()}
		</div>

		<!-- <div class="text-center">
			<p class="text-xs text-gray-500">
				Welcome to <span class="font-semibold text-blue-600">Adel Online Judge</span>
			</p>
		</div> -->
	</div>
</div>
