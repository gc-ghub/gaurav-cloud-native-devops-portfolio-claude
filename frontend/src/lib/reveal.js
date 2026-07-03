// Svelte action: fades/slides an element in when it scrolls into view.
// Skips the animation entirely for prefers-reduced-motion (WCAG 2.1 AA).
export function reveal(node) {
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return {}
  }

  node.classList.add('opacity-0', 'translate-y-4')
  node.style.transition = 'opacity 0.6s ease, transform 0.6s ease'

  const observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          node.classList.remove('opacity-0', 'translate-y-4')
          observer.unobserve(node)
        }
      }
    },
    { threshold: 0.15 },
  )

  observer.observe(node)

  return {
    destroy() {
      observer.disconnect()
    },
  }
}
