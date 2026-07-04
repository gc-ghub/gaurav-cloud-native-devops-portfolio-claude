<script>
  import { brandIcons, genericIcons, skillIcons } from '../lib/techIcons.js'

  let { name } = $props()

  const entry = $derived(skillIcons[name])
  const brand = $derived(entry?.slug ? brandIcons[entry.slug] : null)
  const generic = $derived(entry?.generic ? genericIcons[entry.generic] : null)
</script>

{#if brand}
  <svg width="14" height="14" viewBox="0 0 24 24" fill={brand.hex} aria-hidden="true">
    <path d={brand.d} />
  </svg>
{:else if generic}
  <svg
    width="14"
    height="14"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    {@html generic}
  </svg>
{/if}
