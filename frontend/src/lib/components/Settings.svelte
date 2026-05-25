<script>
  import { onMount } from 'svelte';
  import { Dialogs } from '@wailsio/runtime';
  import { JobService } from '../../../bindings/jobdash';

  let config = { serpapi_key: '', location: 'Remote', resume_path: '', keywords: [], analytics_enabled: false };
  let saved = false;
  let extracting = false;
  let extractResult = '';
  let extractLogs = [];
  let analyticsEnabled = false;

  onMount(() => {
    JobService.GetConfig()
      .then(c => { if (c) { config = c; analyticsEnabled = c.analytics_enabled || false; } })
      .catch(e => console.error(e));
  });

  function toggleAnalytics() {
    analyticsEnabled = !analyticsEnabled;
    JobService.SetAnalyticsEnabled(analyticsEnabled)
      .then(() => JobService.GetConfig())
      .then(c => { if (c) config = c; })
      .catch(e => { console.error(e); analyticsEnabled = !analyticsEnabled; });
  }

  function saveConfig() {
    JobService.SetConfig('location', config.location)
      .then(() => JobService.SetConfig('resume_path', config.resume_path))
      .then(() => { saved = true; setTimeout(() => saved = false, 2000); })
      .catch(e => console.error(e));
  }

  function extractKeywords() {
    extracting = true;
    extractResult = '';
    extractLogs = [];
    JobService.ExtractResume(config.resume_path)
      .then(result => {
        extractLogs = result?.errors || [];
        extractResult = result?.jobs_found ? `${result.jobs_found} keywords extracted via Hermes AI` : 'Done';
        return JobService.GetConfig();
      })
      .then(c => { if (c) config = c; extracting = false; })
      .catch(e => { extractResult = `Error: ${e}`; extracting = false; });
  }

  function browseResume() {
    Dialogs.OpenFile({
      title: 'Select Resume PDF',
      filters: [{ name: 'PDF Files', extensions: ['pdf'] }]
    }).then(path => { if (path) config.resume_path = path; })
      .catch(e => console.error(e));
  }
</script>

<div class="settings">
  <h2>Settings</h2>

  <div class="setting-group env-status">
    <label>SerpAPI Key</label>
    <div class="api-key-row">
      <input type="password" bind:value={config.serpapi_key} placeholder="Paste your SerpAPI key here" class="setting-input api-key-input" />
      <button class="btn btn-save-key" on:click={() => JobService.SetConfig('serpapi_key', config.serpapi_key).then(() => saved = true).catch(e => console.error(e))}>
        Save Key
      </button>
    </div>
    <span class="env-msg" class:ok={config.serpapi_key} class:err={!config.serpapi_key}>
      {config.serpapi_key ? '✓ Key saved' : '✗ Paste your SerpAPI key above'}
    </span>
  </div>

  <div class="setting-group env-status">
    <label>Usage Analytics</label>
    <p class="hint">Help improve JobDash with anonymous feature usage stats. Never includes your resume, job listings, or personal data.</p>
    <div class="toggle-row">
      <span class="toggle-label">{analyticsEnabled ? 'Enabled' : 'Disabled'}</span>
      <button
        class="toggle-switch"
        class:active={analyticsEnabled}
        on:click={toggleAnalytics}
        role="switch"
        aria-checked={analyticsEnabled}
      >
        <span class="toggle-knob"></span>
      </button>
    </div>
  </div>

  <div class="setting-group">
    <label for="location">Default Location</label>
    <input id="location" type="text" bind:value={config.location} placeholder="e.g. Remote, San Francisco" class="setting-input" />
  </div>

  <div class="setting-group">
    <label for="resume">Resume PDF Path</label>
    <p class="hint">Hermes AI reads your resume and extracts keywords for job matching</p>
    <div class="resume-row">
      <input id="resume" type="text" bind:value={config.resume_path}
        placeholder="/path/to/your/resume.pdf" class="setting-input" />
      <button class="btn btn-browse" on:click={browseResume}>Browse</button>
      <button class="btn btn-extract" on:click={extractKeywords} disabled={extracting}>
        {extracting ? 'Extracting...' : 'Extract Keywords'}
      </button>
    </div>
    {#if extractResult}<p class="extract-msg">{extractResult}</p>{/if}
    {#if extractLogs.length > 0}
      <div class="console-log">
        <div class="console-title">Hermes Console</div>
        {#each extractLogs as line}
          <div class="console-line">{line}</div>
        {/each}
      </div>
    {/if}
    {#if config.keywords?.length}
      <div class="keyword-cloud">
        {#each config.keywords.slice(0, 40) as kw}
          <span class="kw-tag">{kw.trim()}</span>
        {/each}
        {#if config.keywords.length > 40}<span class="kw-more">+{config.keywords.length - 40} more</span>{/if}
      </div>
    {/if}
  </div>

  <button class="btn btn-save" on:click={saveConfig}>
    {saved ? '&#10003; Saved' : 'Save Settings'}
  </button>

  <div class="info-box">
    <h3>How it works</h3>
    <ul>
      <li>API key in <code>jobdash/.env</code> — never touches the UI</li>
      <li>Hermes AI reads your resume PDF and extracts keywords (saved to DB)</li>
      <li>Set target job titles on Dashboard — or let Hermes recommend them</li>
      <li>Scraping uses your positions as search queries across Google Jobs</li>
      <li>Track each job through your pipeline: New → Applied → Interviewing → Offer</li>
      <li>Click &#8599; to open the posting and apply manually</li>
    </ul>
  </div>
</div>

<style>
  .settings { max-width: 600px; }
  h2 { font-size: 22px; font-weight: 600; margin-bottom: 24px; color: #e8e8f0; }
  .setting-group { margin-bottom: 20px; }
  label { display: block; font-size: 13px; color: #999; margin-bottom: 6px; font-weight: 600; }
  .hint { font-size: 11px; color: #555; margin-bottom: 8px; }
  .env-status { padding: 12px 16px; background: #0d0f1a; border-radius: 6px; border: 1px solid #1a1d30; }
  .env-ok { font-size: 12px; color: #66bb6a; }
  .env-err { font-size: 12px; color: #ef5350; }
  .env-err code { background: #1a1d30; padding: 1px 6px; border-radius: 3px; font-size: 11px; }
  .env-msg { font-size: 12px; }
  .env-msg.ok { color: #66bb6a; }
  .env-msg.err { color: #ef5350; }
  .toggle-row { display: flex; align-items: center; justify-content: space-between; margin-top: 8px; }
  .toggle-label { font-size: 12px; color: #888; }
  .toggle-switch { position: relative; width: 44px; height: 24px; border-radius: 12px; border: none; background: #1e2240; cursor: pointer; transition: background 0.2s; }
  .toggle-switch.active { background: #2a3a8a; }
  .toggle-knob { position: absolute; top: 2px; left: 2px; width: 20px; height: 20px; border-radius: 50%; background: #ccc; transition: transform 0.2s; }
  .toggle-switch.active .toggle-knob { transform: translateX(20px); background: #ccd5ff; }
  .api-key-row { display: flex; gap: 8px; }
  .api-key-input { flex: 1; }
  .btn-save-key { background: #2a3a8a; color: #ccd5ff; padding: 8px 14px; }
  .btn-save-key:hover { background: #3450b0; }
  .setting-input { width: 100%; background: #151828; border: 1px solid #1e2240; color: #ccc; padding: 10px 14px; border-radius: 6px; font-size: 13px; outline: none; font-family: monospace; }
  .setting-input:focus { border-color: #3a4480; }
  .resume-row { display: flex; gap: 8px; }
  .resume-row .setting-input { flex: 1; }
  .btn { padding: 8px 16px; border: none; border-radius: 6px; font-size: 13px; cursor: pointer; white-space: nowrap; }
  .btn-save { background: #2a3a8a; color: #ccd5ff; padding: 10px 24px; margin-top: 8px; }
  .btn-save:hover { background: #3450b0; }
  .btn-extract { background: #1a3a2a; color: #66bb6a; }
  .btn-extract:hover:not(:disabled) { background: #244a34; }
  .btn-extract:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn-browse { background: #1e2240; color: #8b9dff; }
  .btn-browse:hover { background: #2a3058; }
  .extract-msg { font-size: 12px; color: #66bb6a; margin-top: 6px; }
  .keyword-cloud { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 10px; }
  .kw-tag { background: #1a1f35; color: #8b9dff; font-size: 11px; padding: 3px 10px; border-radius: 12px; }
  .kw-more { font-size: 11px; color: #555; padding: 3px 0; }
  .info-box { margin-top: 32px; padding: 20px; background: #0d0f1a; border: 1px solid #1a1d30; border-radius: 8px; }
  .info-box h3 { font-size: 14px; color: #999; margin-bottom: 12px; }
  .info-box ul { padding-left: 18px; font-size: 12px; color: #777; line-height: 1.8; }
  .info-box code { background: #1a1d30; padding: 1px 6px; border-radius: 3px; font-size: 11px; }
  .console-log { margin-top: 12px; background: #0a0c14; border: 1px solid #1a1d30; border-radius: 6px; padding: 12px; max-height: 200px; overflow-y: auto; }
  .console-title { font-size: 10px; color: #555; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 8px; }
  .console-line { font-size: 11px; color: #7eb87e; font-family: monospace; padding: 2px 0; border-bottom: 1px solid #0d0f18; }
</style>
