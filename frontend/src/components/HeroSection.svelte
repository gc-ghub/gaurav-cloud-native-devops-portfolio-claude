<script>
  import SocialIcon from './SocialIcon.svelte'
  import { socials } from '../lib/socials.js'
</script>

<section id="top" class="mx-auto flex max-w-6xl flex-col-reverse items-center gap-12 px-6 py-20 md:flex-row md:py-32">
  <div class="flex-1 text-center md:text-left">
    <h1 class="text-4xl font-bold tracking-tight text-text md:text-5xl">
      Hi 👋, I'm Gaurav Chaurasia
    </h1>
    <p class="mt-4 text-lg font-medium text-primary">
      Platform Engineer | Kubernetes | GitOps | SRE | Cloud-Native Observability | AI-Powered Incident Response
    </p>
    <p class="mt-6 max-w-xl text-base leading-relaxed text-text-muted">
      I design, build, and operate production-grade cloud-native platforms focused on
      Kubernetes, GitOps, observability, security, and reliability engineering.
    </p>

    <div class="mt-8 flex justify-center gap-3 md:justify-start">
      {#each socials as link}
        <a
          href={link.href}
          target="_blank"
          rel="noreferrer"
          aria-label={link.label}
          class="flex h-10 w-10 items-center justify-center rounded-full border border-border text-text-muted transition-colors hover:border-primary hover:text-primary"
        >
          <SocialIcon name={link.icon} />
        </a>
      {/each}
    </div>
  </div>

  <div class="flex flex-1 justify-center">
    <img
      src="https://avatars.githubusercontent.com/u/161034127?v=4"
      alt="Gaurav Chaurasia"
      class="h-64 w-64 rounded-full border border-border object-cover md:h-80 md:w-80"
    />
  </div>
</section>
