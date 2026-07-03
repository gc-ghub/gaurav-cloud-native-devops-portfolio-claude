<script>
  import { reveal } from '../lib/reveal.js'

  const values = [
    'Automate everything that can be automated',
    'Treat infrastructure as a product',
    'Design for reliability, scalability, and observability',
    'Shift security left through DevSecOps practices',
    'Build platforms that empower engineering teams',
    'Solve operational problems with automation and AI',
  ]
</script>

<section use:reveal class="mx-auto max-w-6xl px-6 py-20">
  <h2 class="text-center text-3xl font-bold text-text">Engineering Philosophy</h2>

  <div class="mt-12 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    {#each values as value}
      <div class="rounded-xl border border-border bg-surface p-6 text-sm font-medium text-text-muted">
        {value}
      </div>
    {/each}
  </div>
</section>
