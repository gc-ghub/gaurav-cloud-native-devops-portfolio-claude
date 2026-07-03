<script>
  import { onMount } from 'svelte'
  import ProjectCard from './ProjectCard.svelte'
  import GithubStatsBar from './GithubStatsBar.svelte'
  import { fetchProjects } from '../lib/api.js'
  import { reveal } from '../lib/reveal.js'

  let projects = $state([])
  let loading = $state(true)
  let error = $state(null)

  onMount(async () => {
    try {
      projects = await fetchProjects()
    } catch (e) {
      error = e.message
    } finally {
      loading = false
    }
  })
</script>

<section id="projects" use:reveal class="mx-auto max-w-6xl px-6 py-20">
  <h2 class="text-center text-3xl font-bold text-text">Featured Projects</h2>
  <p class="mx-auto mt-3 max-w-2xl text-center text-text-muted">
    A selection of platform engineering and DevOps automation projects, synced live from GitHub.
  </p>
  <GithubStatsBar />

  {#if loading}
    <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {#each Array(6) as _}
        <div class="h-48 animate-pulse rounded-xl border border-border bg-surface"></div>
      {/each}
    </div>
  {:else if error}
    <p class="mt-12 text-center text-sm text-text-muted">
      Couldn't load projects right now ({error}) — check them out directly on GitHub instead.
    </p>
  {:else}
    <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {#each projects as project}
        <ProjectCard {...project} />
      {/each}
    </div>
  {/if}

  <div class="mt-10 text-center">
    <a
      href="https://github.com/gc-ghub"
      target="_blank"
      rel="noreferrer"
      class="text-sm font-semibold text-primary hover:underline"
    >
      View full GitHub profile →
    </a>
  </div>
</section>
