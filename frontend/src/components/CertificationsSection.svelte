<script>
  import { onMount } from 'svelte'
  import { fetchCertifications } from '../lib/api.js'
  import { reveal } from '../lib/reveal.js'

  let certifications = $state([])
  let loading = $state(true)
  let error = $state(null)

  onMount(async () => {
    try {
      certifications = await fetchCertifications()
    } catch (e) {
      error = e.message
    } finally {
      loading = false
    }
  })
</script>

<section id="achievements" use:reveal class="mx-auto max-w-6xl px-6 py-20">
  <h2 class="text-center text-3xl font-bold text-text">Certifications &amp; Achievements</h2>
  <p class="mx-auto mt-3 max-w-2xl text-center text-text-muted">
    Synced live from Credly.
  </p>

  {#if loading}
    <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
      {#each Array(4) as _}
        <div class="h-32 animate-pulse rounded-xl border border-border bg-surface"></div>
      {/each}
    </div>
  {:else if error || certifications.length === 0}
    <p class="mt-12 text-center text-sm text-text-muted">
      Certifications are temporarily unavailable — view them directly on Credly instead.
    </p>
  {:else}
    <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
      {#each certifications as cert}
        <a
          href={cert.url}
          target="_blank"
          rel="noreferrer"
          class="flex flex-col items-center gap-3 rounded-xl border border-border bg-surface p-6 text-center transition-transform hover:-translate-y-1 hover:border-primary"
        >
          {#if cert.image_url}
            <img
              src={cert.image_url}
              alt={cert.name}
              loading="lazy"
              decoding="async"
              class="h-14 w-14 object-contain"
            />
          {:else}
            <div class="flex h-14 w-14 items-center justify-center rounded-full bg-primary/10 text-2xl">🏅</div>
          {/if}
          <p class="text-sm font-medium text-text">{cert.name}</p>
          <p class="text-xs text-text-muted">{cert.issuer}</p>
        </a>
      {/each}
    </div>
  {/if}

  <div class="mt-10 text-center">
    <a
      href="https://www.credly.com/users/gaurav-chaurasia.25908149/badges"
      target="_blank"
      rel="noreferrer"
      class="text-sm font-semibold text-primary hover:underline"
    >
      View full Credly profile →
    </a>
  </div>
</section>
