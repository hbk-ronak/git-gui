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
    background: #fafafa;
    border-right: 1px solid #e0e0e0;
  }

  .file-list-header {
    height: 40px;
    padding: 12px 16px;
    background: #f0f0f0;
    border-bottom: 1px solid #e0e0e0;
    font-size: 14px;
    font-weight: 600;
    color: #333;
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
    border-bottom: 1px solid #e8e8e8;
    cursor: pointer;
    box-sizing: border-box;
  }

  .file-item:hover {
    background: #f0f0f0;
  }

  .file-item.selected {
    background: #e3f2fd;
  }

  input[type="checkbox"] {
    width: 16px;
    height: 16px;
    cursor: pointer;
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

  .status-modified { background: #fff3cd; color: #856404; }
  .status-added { background: #d4edda; color: #155724; }
  .status-deleted { background: #f8d7da; color: #721c24; }
  .status-untracked { background: #d1ecf1; color: #0c5460; }
  .status-renamed { background: #d1ecf1; color: #0c5460; }

  .filename {
    font-size: 13px;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    color: #333;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .empty-state {
    padding: 20px;
    text-align: center;
    color: #999;
    font-size: 14px;
  }
</style>
