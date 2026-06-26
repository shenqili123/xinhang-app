<template>
  <header class="site-header" :class="{ scrolled, 'menu-open': menuOpen }">
    <router-link class="brand-block" to="/" aria-label="山东新航实验国际学校首页">
      <img src="/assets/school-name-white.png" alt="山东新航实验国际学校" />
    </router-link>

    <div class="utility-bar">
      <nav class="utility-nav" aria-label="辅助导航">
        <router-link to="/news">{{ t('News', '新闻动态') }}</router-link>
        <template v-if="isLoggedIn">
          <router-link to="/profile">{{ t('My Account', '个人中心') }}</router-link>
          <a href="#" @click.prevent="handleLogout">{{ t('Sign Out', '退出') }}</a>
        </template>
        <template v-else>
          <router-link to="/login">{{ t('Sign In', '登录') }}</router-link>
          <router-link to="/register">{{ t('Register', '注册') }}</router-link>
        </template>
      </nav>
    </div>

    <div class="primary-bar">
      <button class="lang-toggle" type="button" :aria-label="t('切换到中文', 'Switch to English')" @click="toggle">{{ t('中文', 'EN') }}</button>

      <!-- Auth quick links (visible on desktop) -->
      <div class="auth-links-bar">
        <template v-if="isLoggedIn">
          <router-link to="/profile" class="auth-link">{{ user?.name || t('Account', '账户') }}</router-link>
          <a href="#" class="auth-link" @click.prevent="handleLogout">{{ t('Sign Out', '退出') }}</a>
        </template>
        <template v-else>
          <router-link to="/login" class="auth-link">{{ t('Sign In', '登录') }}</router-link>
          <router-link to="/register" class="auth-link">{{ t('Register', '注册') }}</router-link>
        </template>
      </div>

      <button class="nav-toggle" :aria-label="t('Open navigation', '打开导航')" :aria-expanded="String(menuOpen)" @click="menuOpen = !menuOpen">☰</button>
      <nav class="main-nav" :class="{ open: menuOpen }" aria-label="主导航">
        <div class="nav-item">
          <router-link to="/about" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/about' }" @click="navClick('/about')">{{ t('About', '关于新航') }}</router-link>
          <div class="nav-submenu">
            <router-link to="/about" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/about:sub0' }" @click="navClick('/about:sub0')">{{ t('School Overview', '学校概况') }}</router-link>
            <router-link to="/about#character" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/about#character' }" @click="navClick('/about#character')">{{ t('Introduction', '学校简介') }}</router-link>
            <router-link to="/campus" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/campus:walk' }" @click="navClick('/campus:walk')">{{ t('Campus Walk', '漫步校园') }}</router-link>
            <router-link to="/about#exchange" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/about#exchange' }" @click="navClick('/about#exchange')">{{ t('International Exchange', '对外交流') }}</router-link>
          </div>
        </div>
        <div class="nav-item">
          <router-link to="/academics" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/academics' }" @click="navClick('/academics')">{{ t('Academics', '学术课程') }}</router-link>
          <div class="nav-submenu">
            <router-link to="/academics" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/academics:sub0' }" @click="navClick('/academics:sub0')">{{ t('Teaching & Research', '教学教研') }}</router-link>
            <router-link to="/academics#curriculum" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/academics#curriculum' }" @click="navClick('/academics#curriculum')">{{ t('School-based Curriculum', '校本课程') }}</router-link>
            <router-link to="/academics#global" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/academics#global' }" @click="navClick('/academics#global')">{{ t('International Academics', '国际部学术') }}</router-link>
            <router-link to="/academics#pathways" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/academics#pathways' }" @click="navClick('/academics#pathways')">{{ t('Curriculum System', '课程体系') }}</router-link>
          </div>
        </div>
        <div class="nav-item">
          <router-link to="/campus" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/campus' }" @click="navClick('/campus')">{{ t('Campus', '校园') }}</router-link>
          <div class="nav-submenu">
            <router-link to="/news" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/news' }" @click="navClick('/news')">{{ t('Campus News', '校园动态') }}</router-link>
            <router-link to="/campus#culture" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/campus#culture' }" @click="navClick('/campus#culture')">{{ t('Civilized Campus', '文明校园风采') }}</router-link>
            <router-link to="/campus" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/campus:library' }" @click="navClick('/campus:library')">{{ t('Library', '图书馆') }}</router-link>
            <router-link to="/campus#environment" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/campus#environment' }" @click="navClick('/campus#environment')">{{ t('Campus Environment', '校园环境') }}</router-link>
          </div>
        </div>
        <div class="nav-item">
          <router-link to="/student-life" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/student-life' }" @click="navClick('/student-life')">{{ t('Student Life', '学生生活') }}</router-link>
          <div class="nav-submenu">
            <router-link to="/student-life" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/student-life:sub0' }" @click="navClick('/student-life:sub0')">{{ t('Student Activities', '学生活动') }}</router-link>
            <router-link to="/student-life#character" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/student-life#character' }" @click="navClick('/student-life#character')">{{ t('Character Education', '德育天地') }}</router-link>
            <router-link to="/student-life#clubs" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/student-life#clubs' }" @click="navClick('/student-life#clubs')">{{ t('Student Clubs', '学生社团') }}</router-link>
            <router-link to="/student-life#ambassadors" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/student-life#ambassadors' }" @click="navClick('/student-life#ambassadors')">{{ t('Student Ambassadors', '学生大使') }}</router-link>
            <router-link to="/student-life#clubs2" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/student-life#clubs2' }" @click="navClick('/student-life#clubs2')">{{ t('Clubs', '社团') }}</router-link>
            <router-link to="/student-life#arts" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/student-life#arts' }" @click="navClick('/student-life#arts')">{{ t('Arts', '艺术') }}</router-link>
            <router-link to="/student-life#athletics" active-class="" exact-active-class="" :class="{ 'nav-link-active': clickedHref === '/student-life#athletics' }" @click="navClick('/student-life#athletics')">{{ t('Athletics', '运动') }}</router-link>
          </div>
        </div>
        <router-link class="nav-apply" to="/apply" @click="closeSubmenu">{{ t('Apply Now', '立即报名') }}</router-link>
        <div class="nav-mobile-extra">
          <router-link to="/news" @click="closeSubmenu">{{ t('News', '新闻动态') }}</router-link>
          <template v-if="isLoggedIn">
            <router-link to="/profile" @click="closeSubmenu">{{ t('My Account', '个人中心') }}</router-link>
            <a href="#" @click.prevent="handleLogout(); closeSubmenu()">{{ t('Sign Out', '退出') }}</a>
          </template>
          <template v-else>
            <router-link to="/login" @click="closeSubmenu">{{ t('Sign In', '登录') }}</router-link>
            <router-link to="/register" @click="closeSubmenu">{{ t('Register', '注册') }}</router-link>
          </template>
        </div>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLanguage } from '../composables/useLanguage.js'
import { useAuth } from '../composables/useAuth.js'

const { t, toggle } = useLanguage()
const { isLoggedIn, user, logout } = useAuth()
const router = useRouter()
const scrolled = ref(false)
const menuOpen = ref(false)
const clickedHref = ref('')

function navClick(href) {
  clickedHref.value = href
  menuOpen.value = false
  document.activeElement?.blur()
}

function closeSubmenu() {
  menuOpen.value = false
  document.activeElement?.blur()
}

function handleLogout() {
  logout()
  router.push('/')
}

const onScroll = () => { scrolled.value = window.scrollY > 24 }
const onEsc = (e) => { if (e.key === 'Escape') menuOpen.value = false }

onMounted(() => {
  onScroll()
  window.addEventListener('scroll', onScroll, { passive: true })
  window.addEventListener('keydown', onEsc)
})
onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  window.removeEventListener('keydown', onEsc)
})
</script>
