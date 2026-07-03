<script>
  import { theme, toggleTheme } from '../stores/theme.js'

  const navLinks = [
    { label: 'Skills', href: '#skills' },
    { label: 'Featured Projects', href: '#projects' },
    { label: 'Achievements', href: '#achievements' },
  ]

  let menuOpen = $state(false)
</script>

<header
  class="sticky top-0 z-50 border-b border-border bg-bg/80 backdrop-blur"
>
  <div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
    <a href="#top" class="font-mono text-lg font-semibold text-text-muted">
      &lt; <span class="text-primary">Gaurav Chaurasia</span> /&gt;
    </a>

    <nav class="hidden items-center gap-6 md:flex">
      {#each navLinks as link}
        <a
          href={link.href}
          class="text-sm font-medium text-text-muted transition-colors hover:text-text"
        >
          {link.label}
        </a>
      {/each}
      <button
        type="button"
        onclick={toggleTheme}
        aria-label="Toggle theme"
        class="flex h-8 w-8 items-center justify-center rounded-full border border-border text-text-muted transition-colors hover:text-text"
      >
        {#if $theme === 'dark'}
          <!-- sun icon -->
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="4" />
            <path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
          </svg>
        {:else}
          <!-- moon icon -->
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
          </svg>
        {/if}
      </button>
    </nav>

    <button
      type="button"
      class="md:hidden text-text"
      aria-label="Toggle menu"
      onclick={() => (menuOpen = !menuOpen)}
    >
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M4 6h16M4 12h16M4 18h16" />
      </svg>
    </button>
  </div>

  {#if menuOpen}
    <nav class="flex flex-col gap-4 border-t border-border px-6 py-4 md:hidden">
      {#each navLinks as link}
        <a
          href={link.href}
          class="text-sm font-medium text-text-muted hover:text-text"
          onclick={() => (menuOpen = false)}
        >
          {link.label}
        </a>
      {/each}
      <button
        type="button"
        onclick={toggleTheme}
        class="self-start text-sm font-medium text-text-muted hover:text-text"
      >
        Switch to {$theme === 'dark' ? 'light' : 'dark'} mode
      </button>
    </nav>
  {/if}
</header>
