import { ref, watch } from 'vue'

const language = ref(readPreferred())

function readPreferred() {
  try { return localStorage.getItem('xinhang-language') === 'zh' ? 'zh' : 'en' }
  catch { return 'en' }
}

function writePreferred(lang) {
  try { localStorage.setItem('xinhang-language', lang) } catch {}
}

watch(language, (val) => {
  writePreferred(val)
  document.documentElement.lang = val === 'zh' ? 'zh-CN' : 'en'
}, { immediate: true })

export function useLanguage() {
  const toggle = () => { language.value = language.value === 'zh' ? 'en' : 'zh' }
  const isChinese = () => language.value === 'zh'
  const t = (en, zh) => language.value === 'zh' ? zh : en
  return { language, toggle, isChinese, t }
}
