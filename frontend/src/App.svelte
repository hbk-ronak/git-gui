<script>
  import { onMount, onDestroy } from 'svelte'
  import BranchSelector from './components/BranchSelector.svelte'
  import FileList from './components/FileList.svelte'
  import DiffViewer from './components/DiffViewer.svelte'
  import CommitPanel from './components/CommitPanel.svelte'
  import ErrorNotification from './components/ErrorNotification.svelte'
  import { files, branches, selectedFile, currentDiff, isLoading, commitMessage, checkedFiles } from './lib/stores.js'
  import { GetGitStatus, GetGitDiff, GetBranches, GetCurrentRepo, InitRepo } from '../wailsjs/go/backend/App.js'

  let hasRepo = false
  let projectPath = ""
  let pathError = ""

  async function checkRepo() {
    try {
      await GetCurrentRepo()
      hasRepo = true
      loadStatus()
    } catch {
      hasRepo = false
    }
  }

  async function openProject() {
    pathError = ""
    if (!projectPath.trim()) {
      pathError = "Please enter a project path"
      return
    }
    try {
      await InitRepo(projectPath.trim())
      hasRepo = true
      loadStatus()
    } catch (err) {
      pathError = err.message || "Not a valid git repository"
    }
  }

  function changeProject() {
    hasRepo = false
    projectPath = ""
    pathError = ""
    files.set([])
    branches.set([])
    selectedFile.set(null)
    currentDiff.set("")
    commitMessage.set("")
    checkedFiles.set(new Set())
  }

  async function loadStatus() {
    try {
      isLoading.set(true)
      const [statusResult, branchResult] = await Promise.all([
        GetGitStatus(),
        GetBranches()
      ])
      files.set(statusResult || [])
      branches.set(branchResult || [])
    } catch (err) {
      console.error("Failed to load status:", err)
    } finally {
      isLoading.set(false)
    }
  }

  async function loadDiff(file) {
    if (!file) {
      currentDiff.set("")
      return
    }
    try {
      const result = await GetGitDiff(file.Path)
      currentDiff.set(result?.Diff || "")
    } catch (err) {
      console.error("Failed to load diff:", err)
      currentDiff.set("")
    }
  }

  function onWindowFocus() {
    if (hasRepo) loadStatus()
  }

  function handleKeydown(e) {
    if (e.key === 'Enter') openProject()
  }

  onMount(() => {
    checkRepo()
    window.addEventListener('focus', onWindowFocus)
  })

  onDestroy(() => {
    window.removeEventListener('focus', onWindowFocus)
  })

  // Load diff when selected file changes
  $: loadDiff($selectedFile)
</script>

<ErrorNotification />

{#if hasRepo}
  <div class="app-container">
    <header class="header">
      <span class="app-title">Git Commit Tool</span>
      <div class="header-right">
        <button class="change-project-btn" on:click={changeProject}>Change Project</button>
        <BranchSelector />
      </div>
    </header>

    <div class="content">
      <div class="file-panel">
        <FileList />
      </div>
      <div class="diff-panel">
        <DiffViewer />
      </div>
    </div>

    <CommitPanel />
  </div>
{:else}
  <div class="landing-container">
    <div class="landing-card">
      <h1 class="landing-title">Git Commit Tool</h1>
      <p class="landing-subtitle">Enter the path to a git repository to get started.</p>
      <div class="landing-input-row">
        <input
          type="text"
          class="landing-input"
          placeholder="/path/to/your/project"
          bind:value={projectPath}
          on:keydown={handleKeydown}
        />
        <button class="landing-btn" on:click={openProject}>Open</button>
      </div>
      {#if pathError}
        <p class="landing-error">{pathError}</p>
      {/if}
    </div>
  </div>
{/if}

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: hidden;
  }

  .header {
    height: 60px;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-color);
    padding: 0 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-shrink: 0;
    box-sizing: border-box;
  }

  .app-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .change-project-btn {
    padding: 6px 12px;
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
  }

  .change-project-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .content {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .file-panel {
    width: 30%;
    min-width: 250px;
    max-width: 400px;
  }

  .diff-panel {
    flex: 1;
  }

  /* Landing page */
  .landing-container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100vh;
    background: var(--bg-primary);
  }

  .landing-card {
    text-align: center;
    max-width: 500px;
    width: 100%;
    padding: 40px;
  }

  .landing-title {
    font-size: 28px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0 0 8px 0;
  }

  .landing-subtitle {
    color: var(--text-secondary);
    margin: 0 0 24px 0;
    font-size: 14px;
  }

  .landing-input-row {
    display: flex;
    gap: 8px;
  }

  .landing-input {
    flex: 1;
    padding: 10px 14px;
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    font-size: 14px;
    outline: none;
  }

  .landing-input:focus {
    border-color: var(--accent);
  }

  .landing-btn {
    padding: 10px 20px;
    background: var(--accent);
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 600;
  }

  .landing-btn:hover {
    opacity: 0.9;
  }

  .landing-error {
    color: var(--status-deleted-text);
    font-size: 13px;
    margin: 12px 0 0 0;
    text-align: left;
  }
</style>
