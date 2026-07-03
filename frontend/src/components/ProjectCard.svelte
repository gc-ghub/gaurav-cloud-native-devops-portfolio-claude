<script>
  let { name, description, tech, href, stars, forks } = $props()
</script>

<a
  {href}
  target="_blank"
  rel="noreferrer"
  class="flex flex-col rounded-xl border border-border bg-surface p-6 transition-transform hover:-translate-y-1 hover:border-primary"
>
  <h3 class="font-semibold text-text">{name}</h3>
  <p class="mt-2 flex-1 text-sm text-text-muted">{description}</p>

  <div class="mt-4 flex flex-wrap gap-2">
    {#each tech as badge}
      <span class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary">
        {badge}
      </span>
    {/each}
  </div>

  <div class="mt-4 flex items-center gap-4 text-xs text-text-muted">
    <span>⭐ {stars}</span>
    <span>🍴 {forks}</span>
  </div>
</a>
