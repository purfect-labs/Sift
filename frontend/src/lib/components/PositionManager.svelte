<script>
  import { onMount } from 'svelte';
  import { JobService } from '../../../bindings/jobdash';

  export let onUpdate;

  let positions = [];
  let newTitle = '';
  let recommending = false;

  onMount(() => {
    JobService.GetPositions().then(p => positions = p || []).catch(e => console.error(e));
  });

  function addPosition() {
    if (!newTitle.trim()) return;
    positions = [...positions, newTitle.trim()];
    JobService.SavePositions(positions).catch(e => console.error(e));
    newTitle = '';
  }

  function removePosition(title) {
    positions = positions.filter(t => t !== title);
    JobService.SavePositions(positions).catch(e => console.error(e));
  }

  function recommendPositions() {
    recommending = true;
    JobService.RecommendPositions()
      .then(titles => { positions = titles || []; recommending = false; })
      .catch(e => { console.error(e); recommending = false; });
  }
</script>

<div class="position-panel">
  <div class="panel-header">
    <span class="panel-title">&#127919; Target Positions</span>
    <button class="btn btn-ai" on:click={recommendPositions} disabled={recommending}>
      {recommending ? 'Asking Hermes...' : '🤖 AI: Suggest Job Titles from Resume'}
    </button>
  </div>
  {#if positions.length > 0}
    <div class="position-tags">
      {#each positions as title}
        <span class="pos-tag">
          {title}
          <button class="pos-remove" on:click={() => removePosition(title)}>&times;</button>
        </span>
      {/each}
    </div>
  {:else}
    <p class="hint">Add target job titles manually or click "AI Recommend" for Hermes suggestions. These become search queries.</p>
  {/if}
  <div class="add-row">
    <input type="text" placeholder="Add a job title..." bind:value={newTitle} class="add-input"
      on:keydown={(e) => { if (e.key === 'Enter') addPosition(); }} />
    <button class="btn btn-add" on:click={addPosition}>Add</button>
  </div>
</div>

<style>
  .position-panel { background: #0d0f1a; border: 1px solid #1a1d30; border-radius: 8px; padding: 16px; margin-bottom: 16px; }
  .panel-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
  .panel-title { font-size: 14px; font-weight: 600; color: #999; }
  .btn { padding: 6px 14px; border: none; border-radius: 5px; font-size: 12px; cursor: pointer; }
  .btn-ai { background: #1a3a2a; color: #66bb6a; }
  .btn-ai:hover:not(:disabled) { background: #244a34; }
  .btn-ai:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn-add { background: #1e2240; color: #8b9dff; }
  .btn-add:hover { background: #2a3058; }
  .position-tags { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 10px; }
  .pos-tag { background: #1a1f35; color: #8b9dff; font-size: 12px; padding: 4px 10px; border-radius: 14px; display: flex; align-items: center; gap: 6px; }
  .pos-remove { background: none; border: none; color: #ef5350; font-size: 14px; cursor: pointer; padding: 0; line-height: 1; }
  .pos-remove:hover { color: #ff6e6e; }
  .hint { font-size: 11px; color: #555; margin-bottom: 10px; }
  .add-row { display: flex; gap: 8px; }
  .add-input { flex: 1; background: #151828; border: 1px solid #1e2240; color: #ccc; padding: 7px 12px; border-radius: 5px; font-size: 12px; outline: none; }
  .add-input:focus { border-color: #3a4480; }
</style>
