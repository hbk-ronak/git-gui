<script>
  import { branches, currentBranch, files, selectedFile, currentDiff, errorMessage, successMessage, isLoading } from '../lib/stores.js'
  import { SwitchBranch, CreateBranch, GetBranches, GetGitStatus } from '../../wailsjs/go/main/App.js'

  let dropdownOpen = false
  let showNewBranchModal = false
  let newBranchName = ""

  function toggleDropdown() {
    dropdownOpen = !dropdownOpen
  }

  async function refreshAfterBranchChange() {
    const [branchResult, statusResult] = await Promise.all([
      GetBranches(),
      GetGitStatus()
    ])
    branches.set(branchResult || [])
    files.set(statusResult || [])
    selectedFile.set(null)
    currentDiff.set("")
  }

  async function selectBranch(branchName) {
    dropdownOpen = false
    if (branchName === $currentBranch) return
    try {
      isLoading.set(true)
      await SwitchBranch(branchName)
      successMessage.set(`Switched to branch ${branchName}`)
      await refreshAfterBranchChange()
    } catch (err) {
      errorMessage.set(`Failed to switch branch: ${err.message || err}`)
    } finally {
      isLoading.set(false)
    }
  }

  function openNewBranchModal() {
    dropdownOpen = false
    showNewBranchModal = true
    newBranchName = ""
  }

  function closeModal() {
    showNewBranchModal = false
    newBranchName = ""
  }

  async function createBranch() {
    if (!newBranchName.trim()) return
    const name = newBranchName.trim()
    closeModal()
    try {
      isLoading.set(true)
      await CreateBranch(name)
      successMessage.set(`Created and switched to branch ${name}`)
      await refreshAfterBranchChange()
    } catch (err) {
      errorMessage.set(`Failed to create branch: ${err.message || err}`)
    } finally {
      isLoading.set(false)
    }
  }
</script>

<div class="branch-selector">
  <div class="branch-dropdown" role="button" tabindex="0" on:click={toggleDropdown} on:keydown={(e) => e.key === 'Enter' && toggleDropdown()}>
    <span class="branch-name">{$currentBranch}</span>
    <span class="chevron">&#9662;</span>
  </div>

  {#if dropdownOpen}
    <div class="dropdown-menu">
      {#each $branches as branch (branch.Name)}
        <div
          class="dropdown-item"
          role="option"
          aria-selected={branch.IsCurrent}
          tabindex="0"
          on:click={() => selectBranch(branch.Name)}
          on:keydown={(e) => e.key === 'Enter' && selectBranch(branch.Name)}
        >
          {#if branch.IsCurrent}
            <span class="check">&#10003;</span>
          {:else}
            <span class="check"></span>
          {/if}
          <span class:current={branch.IsCurrent}>{branch.Name}</span>
        </div>
      {/each}
    </div>
  {/if}

  <button class="new-branch-btn" on:click={openNewBranchModal}>New Branch</button>
</div>

{#if showNewBranchModal}
  <div class="modal-overlay" role="dialog" on:click={closeModal} on:keydown={(e) => e.key === 'Escape' && closeModal()}>
    <div class="modal-box" on:click|stopPropagation on:keydown|stopPropagation>
      <h3 class="modal-header">Create New Branch</h3>
      <input
        class="modal-input"
        type="text"
        placeholder="feature/my-branch"
        bind:value={newBranchName}
        on:keydown={(e) => e.key === 'Enter' && createBranch()}
      />
      <div class="modal-buttons">
        <button class="cancel-btn" on:click={closeModal}>Cancel</button>
        <button class="create-btn" on:click={createBranch}>Create</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .branch-selector {
    display: flex;
    align-items: center;
    gap: 12px;
    position: relative;
  }

  .branch-dropdown {
    width: 200px;
    height: 36px;
    border: 1px solid #d0d0d0;
    border-radius: 4px;
    padding: 8px 12px;
    font-size: 14px;
    background: #f5f5f5;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    box-sizing: border-box;
  }

  .branch-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: #333;
  }

  .chevron {
    font-size: 8px;
    color: #666;
  }

  .dropdown-menu {
    position: absolute;
    top: 40px;
    left: 0;
    width: 200px;
    max-height: 300px;
    background: #fff;
    border: 1px solid #d0d0d0;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.15);
    overflow-y: auto;
    z-index: 100;
  }

  .dropdown-item {
    height: 36px;
    padding: 8px 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    font-size: 14px;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    color: #333;
    box-sizing: border-box;
  }

  .dropdown-item:hover {
    background: #f0f0f0;
  }

  .check {
    width: 12px;
    color: #0066cc;
    font-size: 12px;
  }

  .current {
    font-weight: 600;
  }

  .new-branch-btn {
    height: 36px;
    padding: 8px 16px;
    border: 1px solid #0066cc;
    border-radius: 4px;
    background: #0066cc;
    color: #fff;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
  }

  .new-branch-btn:hover {
    background: #0052a3;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .modal-box {
    width: 400px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 4px 16px rgba(0,0,0,0.2);
    padding: 24px;
  }

  .modal-header {
    font-size: 18px;
    font-weight: 600;
    color: #333;
    margin: 0 0 20px 0;
  }

  .modal-input {
    width: 100%;
    height: 40px;
    padding: 10px 12px;
    border: 1px solid #d0d0d0;
    border-radius: 4px;
    font-size: 14px;
    background: #fff;
    color: #333;
    margin-bottom: 20px;
    box-sizing: border-box;
  }

  .modal-buttons {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
  }

  .cancel-btn {
    height: 40px;
    padding: 10px 20px;
    border: 1px solid #d0d0d0;
    border-radius: 4px;
    background: transparent;
    color: #666;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
  }

  .create-btn {
    height: 40px;
    padding: 10px 20px;
    border: none;
    border-radius: 4px;
    background: #0066cc;
    color: #fff;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
  }

  .create-btn:hover {
    background: #0052a3;
  }
</style>
