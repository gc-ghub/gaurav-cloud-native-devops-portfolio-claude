<script>
  import { onMount } from 'svelte'
  import { fetchGithubStats } from '../lib/api.js'

  let stats = $state(null)

  onMount(async () => {
    try {
      stats = await fetchGithubStats()
    } catch {
      // Non-critical widget — fail silently, no error UI.
    }
  })
</script>

{#if stats}
  <div class="mt-4 flex justify-center gap-6 text-sm text-text-muted">
    <span>{stats.public_repos} public repos</span>
    <span>⭐ {stats.total_stars} total stars</span>
    <span>🍴 {stats.total_forks} total forks</span>
  </div>
{/if}
