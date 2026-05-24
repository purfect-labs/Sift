<script>
  import { JobService } from '../../../bindings/jobdash';
  import { Browser } from '@wailsio/runtime';

  export let job;
  export let onUpdate;

  let expanded = false;
  let editing = false;
  let notes = job.notes || '';

  const statuses = ['new', 'saved', 'applied', 'interviewing', 'offer', 'rejected', 'not_remote', 'archived'];
  const statusLabels = { new:'New',saved:'Saved',applied:'Applied',interviewing:'Interview',offer:'Offer',rejected:'Rejected',not_remote:'Filtered',archived:'Archived' };
  const statusColors = { new:'#555',saved:'#4fc3f7',applied:'#ffb74d',interviewing:'#ab47bc',offer:'#4caf50',rejected:'#ef5350',not_remote:'#444',archived:'#666' };

  function us(s) { JobService.UpdateStatus(job.id,s).then(()=>{job.status=s;onUpdate()}).catch(e=>console.error(e)) }
  function sn() { JobService.UpdateNotes(job.id,notes).then(()=>{job.notes=notes;editing=false}).catch(e=>console.error(e)) }
  function ol() { if(job.url) Browser.OpenURL(job.url) }
  function dl() { JobService.DeleteJob(job.id).then(()=>onUpdate()).catch(e=>console.error(e)) }

  let ofs=job.offer_salary||'', ofb=job.offer_benefits||'', ofe=job.offer_equity||'', eo=false;
  function so() { JobService.UpdateOffer(job.id,ofs,ofb,ofe).then(()=>{job.offer_salary=ofs;job.offer_benefits=ofb;job.offer_equity=ofe;eo=false}).catch(e=>console.error(e)) }

  function fd(d) { return d?new Date(d).toLocaleDateString():'-' }
</script>

<tr class="jr" class:or={job.status==='offer'} on:click={()=>expanded=!expanded}>
  <td><div class="ms" style="background:hsl({Math.max(0,Math.min(120,job.match_score||0))},50%,30%)">{job.match_score||0}%</div></td>
  <td class="tc"><div class="tl">{job.title||'-'}</div><div class="cl">{job.company||'-'}{#if job.salary} <span class="sa">{job.salary}</span>{/if}</div></td>
  <td><select class="ss" value={job.status} on:change={e=>us(e.target.value)} on:click|stopPropagation style="color:{statusColors[job.status]};border-color:{statusColors[job.status]}">
    {#each statuses as s}<option value={s}>{statusLabels[s]}</option>{/each}</select></td>
  <td class="dc">{fd(job.applied_date)}</td>
  <td class="dc">{fd(job.interview_date)}</td>
  <td class="dc">{fd(job.offer_date)}</td>
  <td class="oc">
    {#if job.status==='offer'}
      {#if eo}
        <div class="oe" on:click|stopPropagation><input placeholder="$" bind:value={ofs} class="oi"/><input placeholder="Benefits" bind:value={ofb} class="oi"/><input placeholder="Equity" bind:value={ofe} class="oi"/><button class="bs" on:click={so}>✓</button></div>
      {:else}
        <div class="os" on:dblclick={e=>{eo=true;e.stopPropagation()}}>
          {#if job.offer_salary}<span class="ov">{job.offer_salary}</span>{/if}
          {#if job.offer_benefits}<span class="ov">{job.offer_benefits}</span>{/if}
          {#if job.offer_equity}<span class="ov">{job.offer_equity}</span>{/if}
          {#if !job.offer_salary&&!job.offer_benefits&&!job.offer_equity}<span class="oh">dbl-click</span>{/if}
        </div>
      {/if}
    {:else}<span class="na">-</span>{/if}
  </td>
  <td class="ac">
    {#if job.url}<button class="ab" on:click|stopPropagation={ol}>&#8599;</button>{/if}
    <button class="ab ad" on:click|stopPropagation={dl}>&#10005;</button>
  </td>
</tr>

<tr class="dr" style="display:{expanded?'table-row':'none'}">
  <td colspan="8">
    <div class="dc2">
      {#if job.skills}
        <div class="ds"><span class="dl2">Skills:</span><div class="st">{#each (job.skills||'').split(',').slice(0,15) as s}{#if s.trim()}<span class="tg">{s.trim()}</span>{/if}{/each}</div></div>
      {/if}
      {#if job.description}
        <div class="ds"><span class="dl2">Description:</span><p class="dt">{job.description}</p></div>
      {/if}
      {#if job.gap_analysis}
        <div class="ds"><span class="dl2">⚠ Gaps:</span><div class="st">{#each (job.gap_analysis||'').split(',') as g}{#if g.trim()}<span class="gt">{g.trim()}</span>{/if}{/each}</div></div>
      {/if}
      <div class="ds"><span class="dl2">Notes:</span>
        {#if editing}
          <textarea bind:value={notes} rows="2" class="ni"></textarea>
          <div class="na2"><button class="bs" on:click={sn}>Save</button><button class="bc" on:click={()=>{editing=false;notes=job.notes||''}}>Cancel</button></div>
        {:else}
          <div class="nd" on:dblclick={()=>editing=true}>{job.notes||'Double-click...'}</div>
        {/if}
      </div>
    </div>
  </td>
</tr>

<style>
  .jr{border-bottom:1px solid #111422;cursor:pointer}.jr:hover{background:#0d0f18}.or,.or:hover{background:#0d1a0d}
  td{padding:6px 8px;font-size:11px;white-space:nowrap;vertical-align:middle}
  .ms{width:36px;height:22px;border-radius:4px;display:flex;align-items:center;justify-content:center;font-size:10px;font-weight:700;color:#fff}
  .tc{max-width:280px}.tl{font-size:12px;font-weight:600;color:#ddd;overflow:hidden;text-overflow:ellipsis}.cl{font-size:10px;color:#777}.sa{color:#66bb6a;font-size:10px}
  .ss{background:transparent;border:1px solid;padding:2px 5px;border-radius:3px;font-size:10px;cursor:pointer;outline:none;font-weight:600}
  .dc{font-size:10px;color:#666;text-align:center}
  .oc{min-width:140px}.oe{display:flex;gap:3px;align-items:center}.oi{width:58px;background:#0d0f1a;border:1px solid #1e2240;color:#ccc;padding:2px 5px;border-radius:3px;font-size:10px}.os{display:flex;gap:3px;flex-wrap:wrap}.ov{background:#1a2a1a;color:#66bb6a;font-size:10px;padding:1px 5px;border-radius:3px}.oh{color:#444;font-size:10px;font-style:italic}.na{color:#333}
  .ac{text-align:center}.ab{background:#1e2240;border:none;color:#8b9dff;padding:2px 6px;border-radius:3px;cursor:pointer;font-size:11px;margin:0 1px}.ab:hover{background:#2a3058}.ad{color:#ef5350}.ad:hover{background:#3a1a1a}
  .dr td{padding:0;background:#090b12}.dc2{padding:10px 16px;display:flex;flex-direction:column;gap:8px}
  .ds{}.dl2{font-size:10px;color:#555;text-transform:uppercase;font-weight:600;display:block;margin-bottom:3px}
  .st{display:flex;flex-wrap:wrap;gap:3px}.tg{background:#1a1f35;color:#8b9dff;font-size:9px;padding:1px 6px;border-radius:8px}.gt{background:#2a1a1a;color:#ef5350;font-size:9px;padding:1px 6px;border-radius:8px}
	.dt{font-size:10px;color:#888;line-height:1.5;margin:0;white-space:pre-wrap}
  .nd{font-size:10px;color:#777;cursor:pointer;font-style:italic}.ni{width:100%;background:#0d0f1a;border:1px solid #1e2240;color:#ccc;padding:5px;border-radius:3px;font-size:10px;resize:vertical;font-family:inherit}
  .na2{display:flex;gap:3px;margin-top:3px}.bs{background:#2a3a8a;color:#ccd5ff;padding:2px 7px;border:none;border-radius:3px;font-size:9px;cursor:pointer}.bc{background:#222;color:#888;padding:2px 7px;border:none;border-radius:3px;font-size:9px;cursor:pointer}
</style>
