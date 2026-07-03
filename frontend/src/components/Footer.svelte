<script>
  import SocialIcon from './SocialIcon.svelte'
  import { socials } from '../lib/socials.js'

  const year = new Date().getFullYear()
</script>

<footer class="border-t border-border">
  <div class="mx-auto flex max-w-6xl flex-col items-center gap-4 px-6 py-10 text-center">
    <div class="flex gap-6 text-sm text-text-muted">
      {#each socials as link}
        <a
          href={link.href}
          target="_blank"
          rel="noreferrer"
          class="flex items-center gap-2 hover:text-primary"
        >
          <SocialIcon name={link.icon} />
          {link.label}
        </a>
      {/each}
    </div>
    <p class="text-xs text-text-muted">
      © {year} Gaurav Chaurasia. Built with Svelte + Go.
    </p>
  </div>
</footer>
