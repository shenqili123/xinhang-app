<template>
  <div ref="root" class="news-page">
    <section class="page-hero" :style="{ '--page-hero-image': 'url(/assets/campus-blossom.jpg)' }" aria-label="News">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('News & Events', '新闻动态') }}</p>
        <h1>{{ t('Stories from Xinhang — campus, classrooms, and beyond.', '来自新航的故事——校园、课堂，以及更远的地方。') }}</h1>
        <p>{{ t('Stay connected with campus life, academic achievements, faculty highlights, and school-wide events.', '关注校园生活、学术成就、教师风采与学校活动。') }}</p>
      </div>
    </section>

    <section class="news-filters">
      <div class="news-filter-bar reveal">
        <button
          v-for="cat in displayCategories"
          :key="cat.key"
          :class="['filter-chip', { active: activeCategory === cat.key }]"
          @click="setCategory(cat.key)"
        >
          {{ cat.label }}
          <span class="chip-count">{{ cat.count }}</span>
        </button>
      </div>
    </section>

    <section class="news-grid-section">
      <div class="news-grid">
        <article v-for="item in newsList" :key="item.id" class="news-card" @click="goDetail(item.id)">
          <div class="news-card-image">
            <img v-if="item.coverImage" :src="item.coverImage" :alt="item.title" loading="lazy" />
            <div v-else class="news-card-placeholder">
              <span>{{ item.category }}</span>
            </div>
          </div>
          <div class="news-card-body">
            <span class="news-card-category">{{ categoryLabel(item.category) }}</span>
            <h3>{{ item.title }}</h3>
            <p>{{ item.summary }}</p>
            <div class="news-card-meta">
              <span v-if="item.authorName">{{ item.authorName }}</span>
              <time>{{ formatDate(item.publishedAt) }}</time>
            </div>
          </div>
        </article>
      </div>

      <div v-if="loading" class="news-loading">{{ t('Loading...', '加载中...') }}</div>

      <div v-if="!loading && newsList.length === 0" class="news-empty">
        {{ t('No articles found in this category.', '该分类暂无文章。') }}
      </div>

      <div v-if="totalPages > 1" class="news-pagination">
        <button :disabled="page <= 1" @click="changePage(page - 1)">&laquo;</button>
        <button
          v-for="p in visiblePages"
          :key="p"
          :class="{ active: p === page }"
          @click="changePage(p)"
        >{{ p }}</button>
        <button :disabled="page >= totalPages" @click="changePage(page + 1)">&raquo;</button>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useLanguage } from '../composables/useLanguage.js'
import { useReveal } from '../composables/useReveal.js'

const { t } = useLanguage()
const router = useRouter()
const root = ref(null)
useReveal(root)

let cardObserver
function observeCards() {
  nextTick(() => {
    if (cardObserver) cardObserver.disconnect()
    cardObserver = new IntersectionObserver(entries => {
      entries.forEach(e => { if (e.isIntersecting) e.target.classList.add('visible') })
    }, { threshold: 0.08 })
    const el = root.value
    if (el) el.querySelectorAll('.news-card').forEach(c => cardObserver.observe(c))
  })
}
onUnmounted(() => { if (cardObserver) cardObserver.disconnect() })

const newsList = ref([])
const categories = ref([])
const activeCategory = ref('')
const page = ref(1)
const total = ref(0)
const pageSize = 12
const loading = ref(false)

const totalPages = computed(() => Math.ceil(total.value / pageSize))

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, page.value - 2)
  const end = Math.min(totalPages.value, page.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const categoryLabels = {
  '': '全部',
  campus: '校园动态',
  news: '新闻报道',
  teachers: '名师风采',
  activities: '学生活动',
  highlights: '学子风采',
  media: '媒体聚焦',
  research: '教育科研',
  moral_education: '德育天地',
  partner_schools: '友好学校',
  home_school: '家校交流',
  party: '党建工作',
  admission: '招生专栏',
  international: '国际交流',
  notice: '通知公告',
  news_en: 'English News',
  other: '其他',
}

const displayCategories = computed(() => {
  const all = [{ key: '', label: t('All', '全部'), count: total.value }]
  for (const cat of categories.value) {
    all.push({
      key: cat.category,
      label: categoryLabels[cat.category] || cat.category,
      count: cat.count,
    })
  }
  return all
})

function categoryLabel(key) {
  return categoryLabels[key] || key
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

async function fetchNews() {
  loading.value = true
  try {
    let url = `/api/news?page=${page.value}&pageSize=${pageSize}`
    if (activeCategory.value) url += `&category=${activeCategory.value}`
    const res = await fetch(url)
    const data = await res.json()
    newsList.value = data.data || []
    total.value = data.total || 0
  } catch {
    newsList.value = []
  }
  loading.value = false
  observeCards()
}

async function fetchCategories() {
  try {
    const res = await fetch('/api/news-categories')
    const data = await res.json()
    categories.value = data.data || []
  } catch {
    categories.value = []
  }
}

function setCategory(key) {
  activeCategory.value = key
  page.value = 1
}

function changePage(p) {
  page.value = p
  window.scrollTo({ top: 300, behavior: 'smooth' })
}

function goDetail(id) {
  router.push(`/news/${id}`)
}

watch([activeCategory, page], fetchNews)

onMounted(() => {
  fetchCategories()
  fetchNews()
})
</script>
