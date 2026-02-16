<script>
  import { commitMessage, checkedFiles, hasCheckedFiles, files, branches, selectedFile, currentDiff, errorMessage, successMessage, isLoading } from '../lib/stores.js'
  import { CommitFiles, CommitAndPush, GetGitStatus, GetBranches } from '../../wailsjs/go/backend/App.js'

  $: canCommit = $hasCheckedFiles && $commitMessage.trim().length > 0

  async function refreshStatus() {
    const [statusResult, branchResult] = await Promise.all([
      GetGitStatus(),
      GetBranches()
    ])
    files.set(statusResult || [])
    branches.set(branchResult || [])
    selectedFile.set(null)
    currentDiff.set("")
    commitMessage.set("")
    checkedFiles.set(new Set())
  }

  async function commit() {
    if (!canCommit) return
    try {
      isLoading.set(true)
      const fileList = Array.from($checkedFiles)
      const result = await CommitFiles(fileList, $commitMessage.trim())
      successMessage.set(`Committed ${result.CommitSHA}: ${result.Message}`)
      await refreshStatus()
    } catch (err) {
      errorMessage.set(`Commit failed: ${err.message || err}`)
    } finally {
      isLoading.set(false)
    }
  }

  async function commitAndPush() {
    if (!canCommit) return
    try {
      isLoading.set(true)
      const fileList = Array.from($checkedFiles)
      const result = await CommitAndPush(fileList, $commitMessage.trim())
      successMessage.set(`Committed and pushed ${result.CommitSHA}: ${result.Message}`)
      await refreshStatus()
    } catch (err) {
      errorMessage.set(`Commit & push failed: ${err.message || err}`)
    } finally {
      isLoading.set(false)
    }
  }
</script>

<div class="commit-panel">
  <label class="commit-label" for="commit-message">Commit Message</label>
  <textarea
    id="commit-message"
    class="commit-textarea"
    placeholder="Enter commit message..."
    bind:value={$commitMessage}
  ></textarea>
  <div class="button-container">
    <button class="commit-btn" disabled={!canCommit} on:click={commit}>
      Commit
    </button>
    <button class="commit-push-btn" disabled={!canCommit} on:click={commitAndPush}>
      Commit & Push
    </button>
  </div>
</div>

<style>
  .commit-panel {
    height: 180px;
    background: var(--bg-secondary);
    border-top: 1px solid var(--border-color);
    padding: 20px;
    box-sizing: border-box;
  }

  .commit-label {
    display: block;
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
  }

  .commit-textarea {
    width: 100%;
    height: 80px;
    padding: 12px;
    border: 1px solid var(--border-light);
    border-radius: 4px;
    font-size: 14px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    resize: vertical;
    background: var(--bg-input);
    color: var(--text-primary);
    box-sizing: border-box;
  }

  .commit-textarea::placeholder {
    color: var(--text-muted);
  }

  .button-container {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    margin-top: 12px;
  }

  .commit-btn {
    height: 40px;
    padding: 10px 24px;
    border: 1px solid var(--accent);
    border-radius: 4px;
    background: transparent;
    color: var(--accent);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
  }

  .commit-btn:hover:not(:disabled) {
    background: var(--accent-subtle);
  }

  .commit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .commit-push-btn {
    height: 40px;
    padding: 10px 24px;
    border: 1px solid var(--accent);
    border-radius: 4px;
    background: var(--accent);
    color: #1e1e1e;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
  }

  .commit-push-btn:hover:not(:disabled) {
    background: var(--accent-hover);
  }

  .commit-push-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
