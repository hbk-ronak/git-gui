<script>
  import BranchSelector from './components/BranchSelector.svelte'
  import FileList from './components/FileList.svelte'
  import DiffViewer from './components/DiffViewer.svelte'
  import CommitPanel from './components/CommitPanel.svelte'
  import ErrorNotification from './components/ErrorNotification.svelte'
  import { files, branches, selectedFile, currentDiff, isLoading } from './lib/stores.js'
  import { GetGitStatus, GetGitDiff, GetBranches } from '../wailsjs/go/main/App.js'

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

  // Load status on mount
  loadStatus()

  // Load diff when selected file changes
  $: loadDiff($selectedFile)
</script>

<ErrorNotification />

<div class="app-container">
  <header class="header">
    <span class="app-title">Git Commit Tool</span>
    <BranchSelector />
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

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: hidden;
  }

  .header {
    height: 60px;
    background: #fff;
    border-bottom: 1px solid #e0e0e0;
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
    color: #333;
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
</style>
