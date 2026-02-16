<script>
  import { selectedFile, currentDiff } from '../lib/stores.js'

  function closeDiff() {
    selectedFile.set(null)
    currentDiff.set("")
  }

  function lineClass(line) {
    if (line.startsWith('+') && !line.startsWith('+++')) return 'line-added'
    if (line.startsWith('-') && !line.startsWith('---')) return 'line-removed'
    if (line.startsWith('@@')) return 'line-hunk'
    return 'line-context'
  }

  $: diffLines = $currentDiff ? $currentDiff.split('\n') : []
</script>

<div class="diff-viewer">
  {#if $selectedFile}
    <div class="diff-header">
      <span class="diff-filepath">{$selectedFile.Path}</span>
      <button class="close-btn" on:click={closeDiff}>&times;</button>
    </div>
    <div class="diff-content">
      {#each diffLines as line}
        <div class="diff-line {lineClass(line)}">{line}</div>
      {/each}
      {#if diffLines.length === 0}
        <div class="empty-state">No changes to display</div>
      {/if}
    </div>
  {:else}
    <div class="empty-state">Select a file to view diff</div>
  {/if}
</div>

<style>
  .diff-viewer {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: #fff;
  }

  .diff-header {
    height: 40px;
    padding: 12px 20px;
    background: #f8f8f8;
    border-bottom: 1px solid #e0e0e0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-sizing: border-box;
  }

  .diff-filepath {
    font-size: 14px;
    font-weight: 600;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    color: #333;
  }

  .close-btn {
    width: 24px;
    height: 24px;
    border: none;
    background: transparent;
    color: #666;
    font-size: 20px;
    cursor: pointer;
    border-radius: 3px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .close-btn:hover {
    background: #e0e0e0;
  }

  .diff-content {
    flex: 1;
    overflow: auto;
    padding: 20px;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.5;
    tab-size: 4;
    white-space: pre;
  }

  .diff-line {
    padding: 0 8px;
  }

  .line-added {
    background: #e6ffed;
    color: #24292e;
  }

  .line-removed {
    background: #ffeef0;
    color: #24292e;
  }

  .line-hunk {
    background: #f6f8fa;
    color: #0366d6;
    font-weight: 600;
    padding: 4px 8px;
    border-radius: 3px;
  }

  .line-context {
    color: #586069;
  }

  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #999;
    font-size: 14px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  }
</style>
