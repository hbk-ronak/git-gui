<script>
  import { commitMessage, checkedFiles, hasCheckedFiles, files, branches, selectedFile, currentDiff, errorMessage, successMessage, isLoading } from '../lib/stores.js'
  import { CommitFiles, CommitAndPush, GetGitStatus, GetBranches } from '../../wailsjs/go/main/App.js'

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
    background: #f8f8f8;
    border-top: 1px solid #e0e0e0;
    padding: 20px;
    box-sizing: border-box;
  }

  .commit-label {
    display: block;
    font-size: 14px;
    font-weight: 600;
    color: #333;
    margin-bottom: 8px;
  }

  .commit-textarea {
    width: 100%;
    height: 80px;
    padding: 12px;
    border: 1px solid #d0d0d0;
    border-radius: 4px;
    font-size: 14px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    resize: vertical;
    background: #fff;
    color: #333;
    box-sizing: border-box;
  }

  .commit-textarea::placeholder {
    color: #999;
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
    border: 1px solid #0066cc;
    border-radius: 4px;
    background: #fff;
    color: #0066cc;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
  }

  .commit-btn:hover:not(:disabled) {
    background: #f0f7ff;
  }

  .commit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .commit-push-btn {
    height: 40px;
    padding: 10px 24px;
    border: 1px solid #0066cc;
    border-radius: 4px;
    background: #0066cc;
    color: #fff;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
  }

  .commit-push-btn:hover:not(:disabled) {
    background: #0052a3;
  }

  .commit-push-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
