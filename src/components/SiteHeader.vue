<template>
  <header class="site-header" :class="{ scrolled, 'menu-open': menuOpen }">
    <router-link class="brand-block" to="/" aria-label="山东新航实验国际学校首页">
      <img src="/images/school-name-white.png" alt="山东新航实验国际学校" />
    </router-link>

    <div class="utility-bar">
      <nav class="utility-nav" aria-label="辅助导航">
        <router-link to="/admission">{{ t('Summer', '夏校') }}</router-link>
        <router-link to="/admission">{{ t('Giving', '支持新航') }}</router-link>
        <router-link to="/student-life">{{ t('Alumni', '校友') }}</router-link>
        <router-link to="/admission#visit">{{ t('Contact', '联系') }}</router-link>
      </nav>
    </div>

    <div class="primary-bar">
      <button class="lang-toggle" type="button" @click="toggle">
        {{ isChinese() ? 'EN' : '中文' }}
      </button>
      <router-link class="btn btn-primary header-register-btn" to="/register">
        {{ t('Register', '注册') }}
      </router-link>
      <button class="nav-toggle" :aria-label="t('Open navigation', '打开导航')" :aria-expanded="menuOpen ? 'true' : 'false'" @click="menuOpen = !menuOpen">☰</button>
      <nav class="main-nav" :class="{ open: menuOpen }" aria-label="主导航">
        <router-link to="/admission" @click="menuOpen = false">{{ t('Admission', '招生') }}</router-link>
        <router-link to="/about" @click="menuOpen = false">{{ t('About', '关于新航') }}</router-link>
        <router-link to="/academics" @click="menuOpen = false">{{ t('Academics', '学术课程') }}</router-link>
        <router-link to="/campus" @click="menuOpen = false">{{ t('Campus', '校园') }}</router-link>
        <router-link to="/student-life" @click="menuOpen = false">{{ t('Student Life', '学生生活') }}</router-link>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useLanguage } from '../composables/useLanguage'

const { toggle, isChinese, t } = useLanguage()
const scrolled = ref(false)
const menuOpen = ref(false)

function onScroll() { scrolled.value = window.scrollY > 24 }
onMounted(() => { onScroll(); window.addEventListener('scroll', onScroll, { passive: true }) })
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>

<style scoped>
.header-register-btn {
  margin-right: 14px;
  font-size: 13px;
  min-height: 34px;
  padding: 0 16px;
}
@media (max-width: 820px) {
  .header-register-btn { display: none; }
}
</style>
