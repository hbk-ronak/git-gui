<script>
  import { files, selectedFile, checkedFiles } from '../lib/stores.js'

  function toggleFile(file) {
    checkedFiles.update(set => {
      const next = new Set(set)
      if (next.has(file.Path)) {
        next.delete(file.Path)
      } else {
        next.add(file.Path)
      }
      return next
    })
  }

  function selectFile(file) {
    selectedFile.set(file)
  }

  function statusLabel(status) {
    return status.charAt(0).toUpperCase() + status.slice(1)
  }
</script>

<div class="file-list">
  <div class="file-list-header">
    Changed Files ({$files.length})
  </div>
  <div class="file-list-items">
    {#each $files as file (file.Path)}
      <div
        class="file-item"
        class:selected={$selectedFile?.Path === file.Path}
        role="option"
        aria-selected={$selectedFile?.Path === file.Path}
        tabindex="0"
        on:click={() => selectFile(file)}
        on:keydown={(e) => e.key === 'Enter' && selectFile(file)}
      >
        <input
          type="checkbox"
          checked={$checkedFiles.has(file.Path)}
          on:click|stopPropagation={() => toggleFile(file)}
        />
        <span class="status-badge status-{file.Status}">{statusLabel(file.Status)}</span>
        <span class="filename">{file.Path}</span>
      </div>
    {/each}
    {#if $files.length === 0}
      <div class="empty-state">No changed files</div>
    {/if}
  </div>
</div>

<style>
  .file-list {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-secondary);
    border-right: 1px solid var(--border-color);
  }

  .file-list-header {
    height: 40px;
    padding: 12px 16px;
    background: var(--bg-tertiary);
    border-bottom: 1px solid var(--border-color);
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
    box-sizing: border-box;
    display: flex;
    align-items: center;
  }

  .file-list-items {
    flex: 1;
    overflow-y: auto;
  }

  .file-item {
    height: 40px;
    padding: 8px 16px;
    display: flex;
    align-items: center;
    gap: 12px;
    border-bottom: 1px solid var(--border-color);
    cursor: pointer;
    box-sizing: border-box;
  }

  .file-item:hover {
    background: var(--bg-tertiary);
  }

  .file-item.selected {
    background: var(--accent-subtle);
  }

  input[type="checkbox"] {
    width: 16px;
    height: 16px;
    cursor: pointer;
    accent-color: var(--accent);
  }

  .status-badge {
    width: 60px;
    height: 20px;
    border-radius: 3px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    text-align: center;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .status-modified { background: var(--status-modified-bg); color: var(--status-modified-text); }
  .status-added { background: var(--status-added-bg); color: var(--status-added-text); }
  .status-deleted { background: var(--status-deleted-bg); color: var(--status-deleted-text); }
  .status-untracked { background: var(--status-untracked-bg); color: var(--status-untracked-text); }
  .status-renamed { background: var(--status-untracked-bg); color: var(--status-untracked-text); }

  .filename {
    font-size: 13px;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    color: var(--text-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .empty-state {
    padding: 20px;
    text-align: center;
    color: var(--text-muted);
    font-size: 14px;
  }
</style>
