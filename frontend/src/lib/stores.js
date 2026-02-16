import { writable, derived } from 'svelte/store'

/** @type {import('svelte/store').Writable<import('../../wailsjs/go/main/App').GitRepo | null>} */
export const repo = writable(null)

/** @type {import('svelte/store').Writable<import('../../wailsjs/go/main/App').FileStatus[]>} */
export const files = writable([])

/** @type {import('svelte/store').Writable<import('../../wailsjs/go/main/App').Branch[]>} */
export const branches = writable([])

/** @type {import('svelte/store').Writable<import('../../wailsjs/go/main/App').FileStatus | null>} */
export const selectedFile = writable(null)

/** @type {import('svelte/store').Writable<string>} */
export const currentDiff = writable("")

/** @type {import('svelte/store').Writable<string>} */
export const commitMessage = writable("")

/** @type {import('svelte/store').Writable<boolean>} */
export const isLoading = writable(false)

/** @type {import('svelte/store').Writable<string | null>} */
export const errorMessage = writable(null)

/** @type {import('svelte/store').Writable<string | null>} */
export const successMessage = writable(null)

export const currentBranch = derived(
    branches,
    $branches => $branches.find(b => b.IsCurrent)?.Name || "unknown"
)

export const checkedFiles = writable(new Set())

export const hasCheckedFiles = derived(
    checkedFiles,
    $checkedFiles => $checkedFiles.size > 0
)
