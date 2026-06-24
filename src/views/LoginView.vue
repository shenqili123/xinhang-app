<template>
  <div ref="root">
    <section class="page-hero" :style="{ '--page-hero-image': 'url(/assets/hero-seminar.png)' }">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('Parent Portal', '家长门户') }}</p>
        <h1>{{ t('Sign in to your account', '登录您的账户') }}</h1>
        <p>{{ t('Access your application status, entrance permits, and admission results.', '查看报名状态、准考证和录取结果。') }}</p>
      </div>
    </section>

    <section class="auth-workspace">
      <form class="auth-form reveal" @submit.prevent="handleLogin">
        <div class="panel-heading">
          <p class="eyebrow">{{ t('Sign In', '登录') }}</p>
          <h2>{{ t('Welcome back', '欢迎回来') }}</h2>
        </div>

        <div v-if="msg" :class="['auth-msg', msgType]">{{ msg }}</div>

        <div class="form-section">
          <div class="form-grid">
            <label class="field wide"><span>{{ t('Email', '邮箱') }}</span><input v-model="form.email" type="email" :placeholder="t('Email address', '请输入邮箱')" required /></label>
            <label class="field wide"><span>{{ t('Password', '密码') }}</span><input v-model="form.password" type="password" :placeholder="t('Your password', '请输入密码')" required /></label>
          </div>
        </div>

        <div class="form-actions">
          <button class="btn btn-primary" type="submit" :disabled="loading">{{ loading ? t('Signing in...', '登录中...') : t('Sign In', '登录') }}</button>
        </div>
        <p class="auth-links">{{ t("Don't have an account?", '还没有账户？') }} <router-link to="/register">{{ t('Create one', '去注册') }}</router-link></p>
      </form>

      <aside class="auth-side reveal">
        <h3>{{ t('New to Xinhang?', '初次了解新航？') }}</h3>
        <p>{{ t('Create an account to track applications, download entrance permits, and check admission results.', '注册账户以跟踪报名、下载准考证和查询录取结果。') }}</p>
        <router-link class="btn btn-primary" to="/register">{{ t('Register Now', '立即注册') }}</router-link>
        <router-link class="btn btn-light" to="/apply" style="margin-top:8px">{{ t('Apply Without Account', '无需账户直接报名') }}</router-link>
      </aside>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useLanguage } from '../composables/useLanguage.js'
import { useReveal } from '../composables/useReveal.js'
import { useAuth } from '../composables/useAuth.js'

const { t } = useLanguage()
const { setAuth } = useAuth()
const router = useRouter()
const root = ref(null)
useReveal(root)

const form = ref({ email: '', password: '' })
const msg = ref('')
const msgType = ref('error')
const loading = ref(false)

async function handleLogin() {
  loading.value = true
  msg.value = ''
  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    })
    const data = await res.json()
    if (res.ok) {
      setAuth(data.token, data.user)
      router.push('/apply')
    } else {
      msgType.value = 'error'
      msg.value = data.error || t('Login failed', '登录失败')
    }
  } catch {
    msgType.value = 'error'
    msg.value = t('Network error', '网络错误')
  }
  loading.value = false
}
</script>
