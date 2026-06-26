<template>
  <div ref="root" class="news-detail-page">
    <section v-if="loading" class="news-detail-loading">
      <p>{{ t('Loading...', '加载中...') }}</p>
    </section>

    <section v-else-if="!article" class="news-detail-empty">
      <h2>{{ t('Article not found', '文章不存在') }}</h2>
      <router-link to="/news" class="btn btn-secondary">{{ t('Back to News', '返回新闻列表') }}</router-link>
    </section>

    <template v-else>
      <section class="news-detail-hero" :style="heroBg">
        <div class="news-detail-hero-inner reveal">
          <router-link to="/news" class="back-link">&larr; {{ t('All News', '返回列表') }}</router-link>
          <span class="news-detail-category">{{ categoryLabel(article.category) }}</span>
          <h1>{{ article.title }}</h1>
          <div class="news-detail-meta">
            <span v-if="article.authorName">{{ article.authorName }}</span>
            <span v-if="article.source">{{ t('Source', '来源') }}：{{ article.source }}</span>
            <time>{{ formatDate(article.publishedAt) }}</time>
          </div>
        </div>
      </section>

      <article class="news-detail-content reveal">
        <div class="article-body" v-html="article.content"></div>
        <div v-if="article.keywords" class="article-keywords">
          <span v-for="kw in keywords" :key="kw" class="keyword-tag">{{ kw }}</span>
        </div>
      </article>

      <nav class="news-detail-nav">
        <router-link to="/news" class="btn btn-secondary">{{ t('Back to News', '返回新闻列表') }}</router-link>
      </nav>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useLanguage } from '../composables/useLanguage.js'

const { t } = useLanguage()
const route = useRoute()
const root = ref(null)

const article = ref(null)
const loading = ref(true)

function revealElements() {
  nextTick(() => {
    if (!root.value) return
    root.value.querySelectorAll('.reveal').forEach(el => el.classList.add('visible'))
  })
}

const categoryLabels = {
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

function categoryLabel(key) {
  return categoryLabels[key] || key
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

const keywords = computed(() => {
  if (!article.value?.keywords) return []
  return article.value.keywords.split(/[,，、\s]+/).filter(Boolean)
})

const heroBg = computed(() => {
  if (article.value?.coverImage) {
    return { '--page-hero-image': `url(${article.value.coverImage})` }
  }
  return { '--page-hero-image': 'url(/assets/campus-blossom.jpg)' }
})

async function fetchArticle() {
  loading.value = true
  try {
    const res = await fetch(`/api/news/${route.params.id}`)
    if (res.ok) {
      const data = await res.json()
      article.value = data.data
    }
  } catch {
    article.value = null
  }
  loading.value = false
  revealElements()
}

onMounted(fetchArticle)
</script>
