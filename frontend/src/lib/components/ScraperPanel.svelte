<script>
  import { onDestroy } from 'svelte';
  import { Events } from '@wailsio/runtime';
  import { JobService } from '../../../bindings/jobdash';

  export let onScrape;

  let scraping = false;
  let result = null;
  let customQuery = '';
  let logs = [];
  let consoleOpen = false;
  let consoleRef;

  let unsub = Events.On('scrape_log', (msg) => {
    logs = [...logs, { time: new Date().toLocaleTimeString(), text: msg.data }];
    // auto-scroll
    if (consoleRef) {
      setTimeout(() => { if (consoleRef) consoleRef.scrollTop = consoleRef.scrollHeight; }, 50);
    }
  });

  onDestroy(() => { if (unsub) unsub(); });

  function runScrape() {
    scraping = true;
    result = null;
    logs = [];
    consoleOpen = true;
    JobService.ScrapeJobs(customQuery)
      .then(r => { result = r; scraping = false; onScrape(); })
      .catch(e => { result = { jobs_found: 0, errors: [String(e)] }; scraping = false; });
  }

  function clearLogs() { logs = []; }
</script>

<div class="scraper-panel">
  <div class="scraper-row">
    <input type="text" placeholder="Custom query (or leave empty to use target positions)..." bind:value={customQuery} class="query-input" disabled={scraping} />
    <button class="btn btn-scrape" on:click={runScrape} disabled={scraping}>
      {#if scraping}&#9881; Scraping...{:else}&#128269; Scrape Jobs{/if}
    </button>
  </div>
  {#if result}
    <div class="scrape-result">
      {#if result.jobs_found > 0}<span class="result-ok">+{result.jobs_found} new jobs</span>{/if}
      {#each result.errors || [] as err}<span class="result-err">{err}</span>{/each}
    </div>
  {/if}
  {#if logs.length > 0}
    <div class="console-toggle" on:click={() => consoleOpen = !consoleOpen}>
      &#9660; Console ({logs.length} lines)
    </div>
    {#if consoleOpen}
      <div class="console" bind:this={consoleRef}>
        {#each logs as log}
          <div class="log-line" class:log-error={log.text.includes('ERROR')} class:log-warn={log.text.includes('WARN')}>
            <span class="log-time">{log.time}</span>
            <span class="log-msg">{log.text}</span>
          </div>
        {/each}
      </div>
      <button class="btn-clear" on:click={clearLogs}>Clear</button>
    {/if}
  {/if}
</div>

<style>
  .scraper-panel { background: #0d0f1a; border: 1px solid #1a1d30; border-radius: 8px; padding: 16px; margin-bottom: 16px; }
  .scraper-row { display: flex; gap: 10px; }
  .query-input { flex: 1; background: #151828; border: 1px solid #1e2240; color: #ccc; padding: 9px 14px; border-radius: 6px; font-size: 13px; outline: none; }
  .query-input:focus { border-color: #3a4480; }
  .btn-scrape { background: #2a3a8a; color: #ccd5ff; padding: 9px 20px; border: none; border-radius: 6px; font-size: 13px; cursor: pointer; white-space: nowrap; }
  .btn-scrape:hover:not(:disabled) { background: #3450b0; }
  .btn-scrape:disabled { opacity: 0.5; cursor: not-allowed; }
  .scrape-result { margin-top: 10px; display: flex; flex-wrap: wrap; gap: 8px; font-size: 12px; }
  .result-ok { color: #66bb6a; }
  .result-err { color: #ef5350; }
  .console-toggle { margin-top: 10px; font-size: 11px; color: #666; cursor: pointer; user-select: none; padding: 4px 0; }
  .console-toggle:hover { color: #999; }
  .console { background: #05070d; border: 1px solid #1a1d30; border-radius: 4px; margin-top: 8px; padding: 10px; max-height: 200px; overflow-y: auto; font-family: 'SF Mono', 'Menlo', 'Monaco', monospace; font-size: 11px; line-height: 1.6; }
  .log-line { display: flex; gap: 8px; }
  .log-time { color: #444; flex-shrink: 0; }
  .log-msg { color: #888; word-break: break-all; }
  .log-error .log-msg { color: #ef5350; }
  .log-warn .log-msg { color: #ffb74d; }
  .btn-clear { background: none; border: none; color: #555; font-size: 10px; cursor: pointer; margin-top: 4px; padding: 2px 0; }
  .btn-clear:hover { color: #999; }
</style>
