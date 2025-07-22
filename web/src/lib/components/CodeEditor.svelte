<script lang="ts">
	import type { SubmissionRequest } from '$lib/api.js';
	import Button from './Button.svelte';

	interface Props {
		onsubmit: (submission: SubmissionRequest) => Promise<void>;
		problemId: number;
		loading?: boolean;
		selectedLanguage?: string;
		onLanguageChange?: (language: string) => void;
	}

	let {
		onsubmit,
		problemId,
		loading = false,
		selectedLanguage = 'c',
		onLanguageChange
	}: Props = $props();
	let code = $state('');
	let submitting = $state(false);
	let hasUserModifiedCode = $state(false);

	const languages = [
		{
			value: 'c',
			label: 'C',
			template: '#include <stdio.h>\n\nint main() {\n    // Your code here\n    return 0;\n}'
		},
		{
			value: 'python',
			label: 'Python',
			template: '# Your code here\n'
		},
		{
			value: 'go',
			label: 'Go',
			template: 'package main\n\nimport "fmt"\n\nfunc main() {\n    // Your code here\n}'
		}
	];

	function getTemplate(language: string): string {
		return languages.find((lang) => lang.value === language)?.template || '';
	}

	let previousLanguage = $state(selectedLanguage);

	$effect(() => {
		if (code === '') {
			code = getTemplate(selectedLanguage);
			hasUserModifiedCode = false;
		} else if (previousLanguage !== selectedLanguage) {
			const previousTemplate = getTemplate(previousLanguage);
			if (!hasUserModifiedCode || code.trim() === previousTemplate.trim()) {
				code = getTemplate(selectedLanguage);
				hasUserModifiedCode = false;
			}
		}
		previousLanguage = selectedLanguage;
	});

	function handleCodeChange() {
		const currentTemplate = getTemplate(selectedLanguage);
		if (code.trim() !== currentTemplate.trim() && code.trim() !== '') {
			hasUserModifiedCode = true;
		} else if (code.trim() === '') {
			hasUserModifiedCode = false;
		}
	}

	function forceTemplateChange() {
		code = getTemplate(selectedLanguage);
		hasUserModifiedCode = false;
	}

	async function handleSubmit() {
		submitting = true;
		try {
			await onsubmit({
				problem_id: problemId,
				code: code,
				language: selectedLanguage
			});
		} finally {
			submitting = false;
		}
	}

	function insertTab(event: KeyboardEvent) {
		if (event.key === 'Tab') {
			event.preventDefault();
			const textarea = event.target as HTMLTextAreaElement;
			const start = textarea.selectionStart;
			const end = textarea.selectionEnd;
			const spaces = '    ';

			textarea.value = textarea.value.substring(0, start) + spaces + textarea.value.substring(end);
			textarea.selectionStart = textarea.selectionEnd = start + spaces.length;
			code = textarea.value;
			handleCodeChange();
		}
	}

	let lineCount = $derived(Math.max(20, code.split('\n').length + 5));
</script>

<div class="space-y-4">
	<!-- Code Editor Header -->
	<div class="flex items-center justify-between">
		<label for="code-editor" class="text-sm font-medium text-gray-700"> Solution Code </label>
		{#if hasUserModifiedCode}
			<button
				type="button"
				onclick={forceTemplateChange}
				class="text-xs text-blue-600 underline hover:text-blue-500"
			>
				Reset to {languages.find((l) => l.value === selectedLanguage)?.label} template
			</button>
		{/if}
	</div>

	<!-- Code Editor -->
	<div class="relative">
		<div class="relative overflow-hidden rounded-lg border border-gray-300 shadow-sm">
			<textarea
				id="code-editor"
				bind:value={code}
				oninput={handleCodeChange}
				onkeydown={insertTab}
				placeholder="Write your solution here..."
				class="block w-full resize-none rounded-lg border-0 bg-gray-50 px-4 py-4 font-mono text-sm focus:border-transparent focus:ring-2 focus:ring-blue-500"
				rows={Math.min(lineCount, 30)}
				disabled={submitting || loading}
			></textarea>

			<!-- Line numbers -->
			<div
				class="pointer-events-none absolute left-2 top-4 select-none font-mono text-xs text-gray-400"
			>
				{#each Array(lineCount) as _, i}
					<div class="h-5 pr-2 text-right leading-5" style="min-width: 2rem;">{i + 1}</div>
				{/each}
			</div>
		</div>
	</div>

	<!-- Editor Actions -->
	<div class="flex items-center justify-between">
		<div class="flex items-center space-x-4 text-xs text-gray-500">
			<span>Lines: {code.split('\n').length}</span>
			<span>Characters: {code.length}</span>
			{#if hasUserModifiedCode}
				<span class="text-orange-600">Modified</span>
			{/if}
		</div>

		<Button
			variant="primary"
			onclick={handleSubmit}
			disabled={submitting || loading}
			loading={submitting}
		>
			{submitting ? 'Submitting...' : 'Submit Solution'}
		</Button>
	</div>
</div>

<style>
	#code-editor {
		padding-left: 3.5rem;
		line-height: 1.25rem;
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
	}
</style>
