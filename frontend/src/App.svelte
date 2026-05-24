<script>
  import { onMount } from 'svelte';
  import Dashboard from './lib/components/Dashboard.svelte';
  import Settings from './lib/components/Settings.svelte';
  import Insights from './lib/components/Insights.svelte';
  import { JobService } from '../bindings/jobdash';

  let view = 'dashboard';
  let counts = {};

  function loadCounts() {
    JobService.GetCounts().then(c => counts = c).catch(e => console.error(e));
  }

  onMount(loadCounts);
  $: view, loadCounts();
</script>

<div class="app">
  <nav class="sidebar">
    <div class="logo">
      <span class="logo-icon">&#9889;</span>
      <span class="logo-text">JobDash</span>
    </div>
    <button class="nav-btn" class:active={view === 'dashboard'} on:click={() => view = 'dashboard'}>
      <span class="nav-icon">&#9632;</span> Jobs
      {#if counts.new}<span class="badge">{counts.new}</span>{/if}
    </button>
    <button class="nav-btn" class:active={view === 'insights'} on:click={() => view = 'insights'}>
      <span class="nav-icon">&#128200;</span> Insights
    </button>
    <div class="spacer"></div>
    <button class="nav-btn" class:active={view === 'settings'} on:click={() => view = 'settings'}>
      <span class="nav-icon">&#9881;</span> Settings
    </button>
  </nav>

  <main class="content">
    {#if view === 'settings'}
      <Settings />
    {:else if view === 'insights'}
      <Insights />
    {:else}
      <Dashboard />
    {/if}
  </main>
</div>

<style>
  :global(*) { margin: 0; padding: 0; box-sizing: border-box; }
  :global(body) { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0a0c14; color: #e0e0e0; overflow: hidden; }
  .app { display: flex; height: 100vh; }
  .sidebar { width: 180px; background: #0d0f1a; border-right: 1px solid #1a1d2e; display: flex; flex-direction: column; padding: 16px 0; flex-shrink: 0; }
  .logo { padding: 8px 16px 24px; font-size: 18px; font-weight: 700; display: flex; align-items: center; gap: 8px; color: #7c8aff; }
  .logo-icon { font-size: 22px; }
  .nav-btn { display: flex; align-items: center; gap: 10px; padding: 10px 16px; margin: 2px 8px; background: none; border: none; color: #8888aa; font-size: 13px; border-radius: 6px; cursor: pointer; transition: all 0.15s; text-align: left; width: calc(100% - 16px); }
  .nav-btn:hover { background: #151828; color: #c0c0dd; }
  .nav-btn.active { background: #1a1f35; color: #8b9dff; }
  .nav-icon { font-size: 15px; width: 20px; text-align: center; }
  .badge { margin-left: auto; background: #2a2f50; color: #8b9dff; font-size: 11px; padding: 2px 8px; border-radius: 10px; }
  .spacer { flex: 1; }
  .content { flex: 1; overflow-y: auto; padding: 24px; }
</style>
