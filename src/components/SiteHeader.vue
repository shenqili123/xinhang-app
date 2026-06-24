<template>
  <header class="site-header" :class="{ scrolled, 'menu-open': menuOpen }">
    <router-link class="brand-block" to="/" aria-label="山东新航实验国际学校首页">
      <img src="/assets/school-name-white.png" alt="山东新航实验国际学校" />
    </router-link>

    <div class="utility-bar">
      <nav class="utility-nav" aria-label="辅助导航">
        <router-link to="/admission">{{ t('Summer', '夏校') }}</router-link>
        <router-link to="/admission">{{ t('Giving', '支持新航') }}</router-link>
        <router-link to="/student-life">{{ t('Alumni', '校友') }}</router-link>
        <router-link to="/admission#visit">{{ t('Contact', '联系') }}</router-link>
        <template v-if="isLoggedIn">
          <router-link to="/profile">{{ user?.name || t('Profile', '个人中心') }}</router-link>
          <a href="#" class="logout-link" @click.prevent="handleLogout">{{ t('Sign Out', '退出') }}</a>
        </template>
        <template v-else>
          <router-link to="/login">{{ t('Sign In', '登录') }}</router-link>
          <router-link to="/register">{{ t('Register', '注册') }}</router-link>
        </template>
      </nav>
    </div>

    <div class="primary-bar">
      <button class="lang-toggle" type="button" @click="toggle" :aria-label="t('Switch to Chinese', '切换到英文')">
        {{ t('中文', 'EN') }}
      </button>
      <button class="nav-toggle" :aria-label="t('Open navigation', '打开导航')" :aria-expanded="String(menuOpen)" @click="menuOpen = !menuOpen">☰</button>
      <nav class="main-nav" :class="{ open: menuOpen }" aria-label="主导航">
        <router-link to="/admission" @click="menuOpen = false">{{ t('Admission', '招生') }}</router-link>
        <router-link to="/about" @click="menuOpen = false">{{ t('About', '关于新航') }}</router-link>
        <router-link to="/academics" @click="menuOpen = false">{{ t('Academics', '学术课程') }}</router-link>
        <router-link to="/campus" @click="menuOpen = false">{{ t('Campus', '校园') }}</router-link>
        <router-link to="/student-life" @click="menuOpen = false">{{ t('Student Life', '学生生活') }}</router-link>
        <template v-if="isLoggedIn">
          <router-link class="nav-user-mobile" to="/profile" @click="menuOpen = false">{{ t('My Profile', '个人中心') }}</router-link>
          <a class="nav-logout-mobile" href="#" @click.prevent="handleLogout(); menuOpen = false">{{ t('Sign Out', '退出') }}</a>
        </template>
        <template v-else>
          <router-link class="nav-auth-mobile" to="/login" @click="menuOpen = false">{{ t('Sign In', '登录') }}</router-link>
        </template>
        <router-link class="nav-apply" to="/apply" @click="menuOpen = false">{{ t('Apply Now', '立即报名') }}</router-link>
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
