<script>
  import { onMount } from 'svelte';
  import { JobService } from '../../../bindings/jobdash';

  let keywords = [];
  let analysis = '';
  let loading = true;
  let analyzing = false;
  let filtering = false;
  let filterResult = null;

  onMount(() => {
    JobService.GetKeywordStats()
      .then(k => { keywords = k || []; loading = false; })
      .catch(e => { console.error(e); loading = false; });
  });

  function runMarketAnalysis() {
    analyzing = true;
    analysis = '';
    JobService.GetMarketAnalysis()
      .then(a => { analysis = a; analyzing = false; })
      .catch(e => { analysis = `Error: ${e}`; analyzing = false; });
  }

  function runRemoteFilter() {
    filtering = true;
    filterResult = null;
    JobService.FilterRemoteJobs()
      .then(r => { filterResult = r; filtering = false; JobService.GetKeywordStats().then(k => keywords = k || []); })
      .catch(e => { filterResult = { errors: [String(e)] }; filtering = false; });
  }

  function getBarWidth(count, max) {
    return max > 0 ? (count / max) * 100 : 0;
  }

  $: maxCount = keywords.length > 0 ? keywords[0].count : 1;
</script>

<div class="insights">
  <h2>Market Insights</h2>

  <div class="section">
    <div class="section-header">
      <h3>Top 100 Keywords Across All Jobs</h3>
      <span class="hint">Frequency of skills/terms in scraped job listings</span>
    </div>

    <div class="action-row">
      <span class="hint">Keywords from Hermes auto-analysis of all scraped jobs</span>
    </div>

    {#if loading}
      <div class="loading">Computing...</div>
    {:else if keywords.length === 0}
      <div class="empty">No job data yet. Scrape some jobs first.</div>
    {:else}
      <div class="keyword-list">
        {#each keywords as kw, i}
          <div class="kw-row">
            <span class="kw-rank">#{i + 1}</span>
            <span class="kw-name">{kw.keyword}</span>
            <div class="kw-bar-bg">
              <div class="kw-bar" style="width: {getBarWidth(kw.count, maxCount)}%"></div>
            </div>
            <span class="kw-count">{kw.count}</span>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .insights { max-width: 900px; }
  h2 { font-size: 22px; font-weight: 600; margin-bottom: 24px; color: #e8e8f0; }
  .section { margin-bottom: 32px; }
  .section-header { margin-bottom: 14px; }
  h3 { font-size: 16px; color: #999; font-weight: 600; }
  .hint { font-size: 11px; color: #555; margin-left: 8px; }
  .loading, .empty { text-align: center; padding: 40px 0; color: #555; }
  .keyword-list { max-height: 500px; overflow-y: auto; background: #0d0f1a; border: 1px solid #1a1d30; border-radius: 8px; padding: 12px; }
  .kw-row { display: flex; align-items: center; gap: 10px; padding: 3px 0; font-size: 12px; }
  .kw-rank { color: #444; width: 30px; text-align: right; flex-shrink: 0; }
  .kw-name { color: #aaa; width: 180px; flex-shrink: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .kw-bar-bg { flex: 1; height: 6px; background: #1a1d30; border-radius: 3px; overflow: hidden; }
  .kw-bar { height: 100%; background: #3a5ccc; border-radius: 3px; transition: width 0.3s; }
  .kw-count { color: #666; width: 35px; text-align: right; flex-shrink: 0; }
  .btn-analyze { background: #1a3a2a; color: #66bb6a; padding: 10px 20px; border: none; border-radius: 6px; font-size: 13px; cursor: pointer; }
  .btn-analyze:hover:not(:disabled) { background: #244a34; }
  .btn-analyze:disabled { opacity: 0.5; }
  .btn-filter { background: #2a3a8a; color: #ccd5ff; padding: 10px 20px; border: none; border-radius: 6px; font-size: 13px; cursor: pointer; }
  .btn-filter:hover:not(:disabled) { background: #3450b0; }
  .btn-filter:disabled { opacity: 0.5; }
  .action-row { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
  .filter-result { font-size: 12px; color: #66bb6a; }
  .analysis-output { background: #0d0f1a; border: 1px solid #1a2040; border-radius: 8px; padding: 16px; margin-top: 12px; font-size: 13px; color: #aaa; line-height: 1.7; }
</style>
