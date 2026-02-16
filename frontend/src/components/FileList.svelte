<script>
  import { files, selectedFile, checkedFiles } from '../lib/stores.js'

  let collapsedFolders = {}

  function buildTree(fileList) {
    const root = { name: '', path: '', children: {}, files: [] }
    for (const file of fileList) {
      const parts = file.Path.split('/')
      let current = root
      for (let i = 0; i < parts.length - 1; i++) {
        const folderName = parts[i]
        if (!current.children[folderName]) {
          current.children[folderName] = {
            name: folderName,
            path: parts.slice(0, i + 1).join('/'),
            children: {},
            files: []
          }
        }
        current = current.children[folderName]
      }
      current.files.push(file)
    }
    return root
  }

  function getAllFilesInFolder(node) {
    let result = [...node.files]
    for (const child of Object.values(node.children)) {
      result = result.concat(getAllFilesInFolder(child))
    }
    return result
  }

  function flattenTree(node, depth) {
    const items = []
    const folders = Object.values(node.children).sort((a, b) => a.name.localeCompare(b.name))
    const sortedFiles = [...node.files].sort((a, b) => {
      const aName = a.Path.split('/').pop()
      const bName = b.Path.split('/').pop()
      return aName.localeCompare(bName)
    })
    for (const folder of folders) {
      items.push({ type: 'folder', node: folder, depth })
      if (!collapsedFolders[folder.path]) {
        items.push(...flattenTree(folder, depth + 1))
      }
    }
    for (const file of sortedFiles) {
      items.push({ type: 'file', file, depth })
    }
    return items
  }

  $: tree = buildTree($files)
  $: flatItems = flattenTree(tree, 0)
  $: allChecked = $files.length > 0 && $files.every(f => $checkedFiles.has(f.Path))
  $: someChecked = !allChecked && $files.some(f => $checkedFiles.has(f.Path))

  function toggleAll() {
    if (allChecked) {
      checkedFiles.set(new Set())
    } else {
      checkedFiles.set(new Set($files.map(f => f.Path)))
    }
  }

  function folderCheckState(node, checked) {
    const folderFiles = getAllFilesInFolder(node)
    const checkedCount = folderFiles.filter(f => checked.has(f.Path)).length
    if (checkedCount === 0) return 'none'
    if (checkedCount === folderFiles.length) return 'all'
    return 'some'
  }

  function toggleFolder(node) {
    const folderFiles = getAllFilesInFolder(node)
    checkedFiles.update(set => {
      const next = new Set(set)
      const allIn = folderFiles.every(f => next.has(f.Path))
      if (allIn) {
        folderFiles.forEach(f => next.delete(f.Path))
      } else {
        folderFiles.forEach(f => next.add(f.Path))
      }
      return next
    })
  }

  function toggleCollapse(path) {
    collapsedFolders[path] = !collapsedFolders[path]
    collapsedFolders = collapsedFolders
  }

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

  function fileName(path) {
    return path.split('/').pop()
  }
</script>

<div class="file-list">
  <div class="file-list-header">
    <input
      type="checkbox"
      checked={allChecked}
      indeterminate={someChecked}
      on:click|stopPropagation={toggleAll}
    />
    <span>Changed Files ({$files.length})</span>
  </div>
  <div class="file-list-items">
    {#each flatItems as item (item.type === 'folder' ? 'f:' + item.node.path : item.file.Path)}
      {#if item.type === 'folder'}
        {@const state = folderCheckState(item.node, $checkedFiles)}
        <div
          class="folder-item"
          style="padding-left: {12 + item.depth * 16}px"
        >
          <span
            class="folder-chevron"
            role="button"
            tabindex="0"
            on:click={() => toggleCollapse(item.node.path)}
            on:keydown={(e) => e.key === 'Enter' && toggleCollapse(item.node.path)}
          >{collapsedFolders[item.node.path] ? '\u25B8' : '\u25BE'}</span>
          <input
            type="checkbox"
            checked={state === 'all'}
            indeterminate={state === 'some'}
            on:click|stopPropagation={() => toggleFolder(item.node)}
          />
          <span class="folder-name">{item.node.name}/</span>
        </div>
      {:else}
        <div
          class="file-item"
          class:selected={$selectedFile?.Path === item.file.Path}
          style="padding-left: {12 + item.depth * 16}px"
          role="option"
          aria-selected={$selectedFile?.Path === item.file.Path}
          tabindex="0"
          on:click={() => selectFile(item.file)}
          on:keydown={(e) => e.key === 'Enter' && selectFile(item.file)}
        >
          <input
            type="checkbox"
            checked={$checkedFiles.has(item.file.Path)}
            on:click|stopPropagation={() => toggleFile(item.file)}
          />
          <span class="status-badge status-{item.file.Status}">{statusLabel(item.file.Status)}</span>
          <span class="filename">{fileName(item.file.Path)}</span>
        </div>
      {/if}
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
    gap: 10px;
  }

  .file-list-items {
    flex: 1;
    overflow-y: auto;
  }

  .folder-item {
    height: 32px;
    padding-right: 16px;
    display: flex;
    align-items: center;
    gap: 6px;
    box-sizing: border-box;
  }

  .folder-chevron {
    width: 16px;
    font-size: 10px;
    color: var(--text-secondary);
    cursor: pointer;
    text-align: center;
    flex-shrink: 0;
    user-select: none;
  }

  .folder-name {
    font-size: 13px;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    color: var(--text-secondary);
    font-weight: 600;
  }

  .file-item {
    height: 36px;
    padding-right: 16px;
    display: flex;
    align-items: center;
    gap: 8px;
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
    flex-shrink: 0;
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
