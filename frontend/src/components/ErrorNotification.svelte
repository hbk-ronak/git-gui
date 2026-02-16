<script>
  import { errorMessage, successMessage } from '../lib/stores.js'

  function dismissError() {
    errorMessage.set(null)
  }

  function dismissSuccess() {
    successMessage.set(null)
  }

  // Auto-dismiss after 5 seconds
  $: if ($errorMessage) {
    setTimeout(() => errorMessage.set(null), 5000)
  }
  $: if ($successMessage) {
    setTimeout(() => successMessage.set(null), 5000)
  }
</script>

{#if $errorMessage}
  <div class="notification error">
    <span>{$errorMessage}</span>
    <button class="dismiss" on:click={dismissError}>&times;</button>
  </div>
{/if}

{#if $successMessage}
  <div class="notification success">
    <span>{$successMessage}</span>
    <button class="dismiss" on:click={dismissSuccess}>&times;</button>
  </div>
{/if}

<style>
  .notification {
    position: fixed;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    max-width: 600px;
    padding: 12px 20px;
    border-radius: 4px;
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    z-index: 2000;
  }

  .error {
    background: var(--notify-error-bg);
    border: 1px solid var(--notify-error-border);
    color: var(--notify-error-text);
  }

  .success {
    background: var(--notify-success-bg);
    border: 1px solid var(--notify-success-border);
    color: var(--notify-success-text);
  }

  .dismiss {
    border: none;
    background: transparent;
    font-size: 20px;
    cursor: pointer;
    color: inherit;
    padding: 0;
    line-height: 1;
  }
</style>
