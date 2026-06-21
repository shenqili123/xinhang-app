<template>
  <section class="auth-page">
    <div class="auth-card">
      <h1>{{ t('Log In', '登录') }}</h1>
      <p class="auth-subtitle">{{ t('Sign in to manage your application.', '登录后管理您的报名申请。') }}</p>

      <form @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <label>{{ t('Email', '邮箱') }}</label>
          <input v-model="form.email" type="email" required :placeholder="t('Enter your email', '请输入邮箱')" />
        </div>
        <div class="form-group">
          <label>{{ t('Password', '密码') }}</label>
          <input v-model="form.password" type="password" required :placeholder="t('Enter your password', '请输入密码')" />
        </div>

        <p v-if="error" class="form-error">{{ error }}</p>
        <p v-if="success" class="form-success">{{ success }}</p>

        <button type="submit" class="btn btn-primary form-submit" :disabled="loading">
          {{ loading ? t('Logging in...', '登录中...') : t('Log In', '登录') }}
        </button>
      </form>

      <p class="auth-switch">
        {{ t("Don't have an account?", '还没有账号？') }}
        <router-link to="/register">{{ t('Register', '注册') }}</router-link>
      </p>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useLanguage } from '../composables/useLanguage'

const { t } = useLanguage()
const router = useRouter()

const form = reactive({ email: '', password: '' })
const error = ref('')
const success = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Login failed')
    localStorage.setItem('xinhang-token', data.token)
    success.value = t('Login successful!', '登录成功！')
    setTimeout(() => router.push('/'), 1000)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>
