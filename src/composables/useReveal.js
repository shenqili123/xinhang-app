import { onMounted, onUnmounted } from 'vue'

export function useReveal(rootRef) {
  let io
  onMounted(() => {
    const root = rootRef?.value || document
    io = new IntersectionObserver(entries => {
      entries.forEach(e => { if (e.isIntersecting) e.target.classList.add('visible') })
    }, { threshold: 0.16 })
    root.querySelectorAll('.reveal').forEach(el => io.observe(el))
  })
  onUnmounted(() => { io?.disconnect() })
}
