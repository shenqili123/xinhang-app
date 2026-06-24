import { ref, computed, watchEffect } from 'vue'

const lang = ref(readStorage())

function readStorage() {
  try {
    return localStorage.getItem('xinhang-language') === 'zh' ? 'zh' : 'en'
  } catch { return 'en' }
}

watchEffect(() => {
  document.documentElement.lang = lang.value === 'zh' ? 'zh-CN' : 'en'
  try { localStorage.setItem('xinhang-language', lang.value) } catch {}
})

export function useLanguage() {
  const isChinese = computed(() => lang.value === 'zh')
  const t = (en, zh) => lang.value === 'zh' ? zh : en
  const toggle = () => { lang.value = lang.value === 'zh' ? 'en' : 'zh' }
  return { lang, isChinese, t, toggle }
}
