# Custom Features Backup

This document preserves all custom functionality added to the frontend that must be re-integrated after the design migration.

## 1. useAuth.js

```js
import { ref, computed } from 'vue'

const token = ref(localStorage.getItem('token') || '')
const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

export function useAuth() {
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setAuth(t, u) {
    if (t !== undefined) { token.value = t; localStorage.setItem('token', t) }
    if (u !== undefined) { user.value = u; localStorage.setItem('user', JSON.stringify(u)) }
  }

  function logout() {
    token.value = ''; user.value = null
    localStorage.removeItem('token'); localStorage.removeItem('user')
  }

  function authHeader() {
    return token.value ? { Authorization: `Bearer ${token.value}` } : {}
  }

  return { token, user, isLoggedIn, isAdmin, setAuth, logout, authHeader }
}
```

## 2. Router (index.js)

```js
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Home', component: () => import('../views/HomeView.vue') },
  { path: '/about', name: 'About', component: () => import('../views/AboutView.vue') },
  { path: '/academics', name: 'Academics', component: () => import('../views/AcademicsView.vue') },
  { path: '/admission', name: 'Admission', component: () => import('../views/AdmissionView.vue') },
  { path: '/campus', name: 'Campus', component: () => import('../views/CampusView.vue') },
  { path: '/student-life', name: 'StudentLife', component: () => import('../views/StudentLifeView.vue') },
  { path: '/apply', name: 'Apply', component: () => import('../views/ApplyView.vue'), meta: { requiresAuth: true } },
  { path: '/profile', name: 'Profile', component: () => import('../views/ProfileView.vue'), meta: { requiresAuth: true } },
  { path: '/register', name: 'Register', component: () => import('../views/RegisterView.vue'), meta: { guestOnly: true } },
  { path: '/login', name: 'Login', component: () => import('../views/LoginView.vue'), meta: { guestOnly: true } },
  { path: '/verify', name: 'Verify', component: () => import('../views/VerifyView.vue') },
  { path: '/news', name: 'News', component: () => import('../views/NewsView.vue') },
  { path: '/news/:id', name: 'NewsDetail', component: () => import('../views/NewsDetailView.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to) {
    if (to.hash) return { el: to.hash, behavior: 'smooth' }
    return { top: 0, behavior: 'smooth' }
  },
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.guestOnly && token) {
    next({ name: 'Profile' })
  } else {
    next()
  }
})

export default router
```

## 3. API Endpoints Used

- POST /api/send-code { email }
- POST /api/register { name, email, phone, password, code }
- POST /api/login { email, password }
- GET /api/profile (JWT)
- PUT /api/profile { name, phone } (JWT)
- GET /api/my-applications (JWT)
- POST /api/apply { ...fields } (JWT)
- GET /api/query-permit?q=xxx
- GET /api/permit-qr?no=xxx&student=xxx
- GET /api/news?page=N&pageSize=N&category=xxx
- GET /api/news/:id
- GET /api/news-categories

## 4. Header Auth UI

In SiteHeader template:
```vue
<template v-if="isLoggedIn">
  <router-link to="/profile">{{ user?.name || t('Profile', '个人中心') }}</router-link>
  <a href="#" class="logout-link" @click.prevent="handleLogout">{{ t('Sign Out', '退出') }}</a>
</template>
<template v-else>
  <router-link to="/login">{{ t('Sign In', '登录') }}</router-link>
  <router-link to="/register">{{ t('Register', '注册') }}</router-link>
</template>
```

## 5. Complete View Files

All view files (LoginView, RegisterView, ProfileView, ApplyView, NewsView, NewsDetailView) are preserved in git history and documented above in the conversation context.
