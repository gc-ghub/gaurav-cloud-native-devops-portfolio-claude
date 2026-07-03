<script>
  import { reveal } from '../lib/reveal.js'

  const skillGroups = [
    { title: 'Container & Orchestration', items: ['Kubernetes', 'Docker', 'Helm', 'Istio', 'Kyverno'] },
    { title: 'CI/CD & Automation', items: ['GitHub Actions', 'Jenkins', 'Docker Compose'] },
    { title: 'Infrastructure as Code', items: ['Terraform', 'Ansible', 'Helm'] },
    { title: 'Cloud Platforms', items: ['AWS', 'Azure', 'GCP'] },
    { title: 'Observability & Monitoring', items: ['Prometheus', 'Grafana', 'Loki', 'Tempo', 'Kiali'] },
    { title: 'Programming & Scripting', items: ['Python', 'Go', 'Bash', 'FastAPI'] },
    { title: 'Cloud-Native Tools', items: ['Argo CD', 'Kafka', 'Gitleaks', 'Trivy', 'CodeQL'] },
  ]
</script>

<section id="skills" use:reveal class="mx-auto max-w-6xl px-6 py-20">
  <h2 class="text-center text-3xl font-bold text-text">Skills &amp; Expertise</h2>
  <p class="mx-auto mt-3 max-w-2xl text-center text-text-muted">
    Core DevOps and Platform Engineering capabilities, grouped by focus area.
  </p>

  <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
    {#each skillGroups as group}
      <div class="rounded-xl border border-border bg-surface p-6 backdrop-blur">
        <h3 class="font-semibold text-primary">{group.title}</h3>
        <div class="mt-4 flex flex-wrap gap-2">
          {#each group.items as item}
            <span class="rounded-full border border-border px-3 py-1 text-xs font-medium text-text-muted">
              {item}
            </span>
          {/each}
        </div>
      </div>
    {/each}
  </div>
</section>
