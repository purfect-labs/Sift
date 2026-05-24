<script>
  import { onMount } from 'svelte';
  import JobCard from './JobCard.svelte';
  import ScraperPanel from './ScraperPanel.svelte';
  import PositionManager from './PositionManager.svelte';
  import { JobService } from '../../../bindings/jobdash';

  let jobs = [];
  let loading = true;
  let search = '';
  let sortBy = 'match_score';
  let statusFilter = '';

  function loadJobs() {
    loading = true;
    JobService.GetJobs(statusFilter)
      .then(j => { jobs = j || []; loading = false; })
      .catch(e => { console.error(e); jobs = []; loading = false; });
  }

  onMount(loadJobs);

  function clearAll() {
    JobService.ClearAllJobs().then(() => loadJobs()).catch(e => alert('Clear failed: ' + e));
  }

  $: filteredJobs = (jobs || []).filter(j => {
    if (!search) return true;
    const q = search.toLowerCase();
    return (j.title || '').toLowerCase().includes(q) ||
           (j.company || '').toLowerCase().includes(q) ||
           (j.skills || '').toLowerCase().includes(q);
  }).sort((a, b) => {
    if (sortBy === 'match_score') return (b.match_score || 0) - (a.match_score || 0);
    if (sortBy === 'recent') return (b.scraped_at || '').localeCompare(a.scraped_at || '');
    if (sortBy === 'company') return (a.company || '').localeCompare(b.company || '');
    return 0;
  });

  $: statusFilter, loadJobs();

  const filterOptions = [
    { value: '', label: 'All Jobs' },
    { value: 'new', label: 'New' },
    { value: 'saved', label: 'Saved' },
    { value: 'applied', label: 'Applied' },
    { value: 'interviewing', label: 'Interviewing' },
    { value: 'offer', label: 'Offers' },
    { value: 'not_remote', label: 'Filtered Out' },
    { value: 'rejected', label: 'Rejected' },
    { value: 'archived', label: 'Archived' },
  ];
</script>

<div class="dashboard">
  <div class="header">
    <div>
      <h2>Jobs</h2>
      <span class="count">{filteredJobs.length} jobs</span>
    </div>
    <div class="controls">
      <select bind:value={statusFilter} class="filter-select">
        {#each filterOptions as opt}
          <option value={opt.value}>{opt.label}</option>
        {/each}
      </select>
      <input type="text" placeholder="Search..." bind:value={search} class="search-input" />
      <select bind:value={sortBy} class="sort-select">
        <option value="match_score">Match %</option>
        <option value="recent">Recent</option>
        <option value="company">Company</option>
      </select>
      <button class="btn btn-primary" on:click={loadJobs}>&#8635;</button>
      <button class="btn btn-danger" on:click={clearAll}>Clear All</button>
    </div>
  </div>

  <ScraperPanel onScrape={loadJobs} />
  <PositionManager onUpdate={loadJobs} />

  {#if loading}
    <div class="loading">Loading...</div>
  {:else if filteredJobs.length === 0}
    <div class="empty">
      <div class="empty-icon">&#128269;</div>
      <p>No jobs. Scrape some or change the filter.</p>
    </div>
  {:else}
    <div class="table-wrap">
      <table class="job-table">
        <thead>
          <tr>
            <th>%</th>
            <th>Title / Company</th>
            <th>Status</th>
            <th>Applied</th>
            <th>Interview</th>
            <th>Offer</th>
            <th>Offer Details</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {#each filteredJobs as job (job.id)}
            <JobCard {job} onUpdate={loadJobs} />
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  .dashboard { max-width: 1200px; }
  .header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 20px; flex-wrap: wrap; gap: 12px; }
  h2 { font-size: 22px; font-weight: 600; color: #e8e8f0; }
  .count { font-size: 12px; color: #666; margin-left: 8px; }
  .controls { display: flex; gap: 6px; align-items: center; flex-wrap: wrap; }
  .filter-select { background: #151828; border: 1px solid #1e2240; color: #ccc; padding: 8px 10px; border-radius: 6px; font-size: 13px; outline: none; }
  .search-input { background: #151828; border: 1px solid #1e2240; color: #ccc; padding: 8px 14px; border-radius: 6px; font-size: 13px; width: 180px; outline: none; }
  .search-input:focus { border-color: #3a4480; }
  .search-input::placeholder { color: #555; }
  .sort-select { background: #151828; border: 1px solid #1e2240; color: #ccc; padding: 8px 10px; border-radius: 6px; font-size: 13px; outline: none; }
  .btn { padding: 8px 14px; border: none; border-radius: 6px; font-size: 13px; cursor: pointer; }
  .btn-primary { background: #2a3a8a; color: #ccd5ff; }
  .btn-primary:hover { background: #3450b0; }
  .btn-danger { background: #3a1a1a; color: #ef5350; }
  .btn-danger:hover { background: #502020; }
  .loading, .empty { text-align: center; padding: 60px 0; color: #555; }
  .empty-icon { font-size: 48px; margin-bottom: 16px; }
  .table-wrap { overflow-x: auto; }
  .job-table { width: 100%; border-collapse: collapse; }
  .job-table th { text-align: left; padding: 8px 10px; font-size: 10px; color: #555; text-transform: uppercase; font-weight: 600; border-bottom: 1px solid #1a1d30; position: sticky; top: 0; background: #0a0c14; }
  .job-table th:last-child { text-align: center; }
</style>
