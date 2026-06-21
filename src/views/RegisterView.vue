<template>
  <section class="auth-page">
    <div class="auth-card">
      <h1>{{ t('Create Account', '注册账号') }}</h1>
      <p class="auth-subtitle">{{ t('Register to explore campus and apply for admission.', '注册后可浏览校园信息并提交报名申请。') }}</p>

      <form @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <label>{{ t('Full Name', '姓名') }}</label>
          <input v-model="form.name" type="text" required :placeholder="t('Enter your name', '请输入姓名')" />
        </div>
        <div class="form-group">
          <label>{{ t('Email', '邮箱') }}</label>
          <input v-model="form.email" type="email" required :placeholder="t('Enter your email', '请输入邮箱')" />
        </div>
        <div class="form-group">
          <label>{{ t('Phone', '手机号') }}</label>
          <input v-model="form.phone" type="tel" required :placeholder="t('Enter your phone number', '请输入手机号')" />
        </div>
        <div class="form-group">
          <label>{{ t('Password', '密码') }}</label>
          <input v-model="form.password" type="password" required minlength="6" :placeholder="t('At least 6 characters', '至少6位字符')" />
        </div>
        <div class="form-group">
          <label>{{ t('Confirm Password', '确认密码') }}</label>
          <input v-model="form.confirmPassword" type="password" required :placeholder="t('Re-enter password', '请再次输入密码')" />
        </div>

        <p v-if="error" class="form-error">{{ error }}</p>
        <p v-if="success" class="form-success">{{ success }}</p>

        <button type="submit" class="btn btn-primary form-submit" :disabled="loading">
          {{ loading ? t('Registering...', '注册中...') : t('Register', '注册') }}
        </button>
      </form>

      <p class="auth-switch">
        {{ t('Already have an account?', '已有账号？') }}
        <router-link to="/login">{{ t('Log in', '登录') }}</router-link>
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

const form = reactive({ name: '', email: '', phone: '', password: '', confirmPassword: '' })
const error = ref('')
const success = ref('')
const loading = ref(false)

async function handleRegister() {
  error.value = ''
  success.value = ''

  if (form.password !== form.confirmPassword) {
    error.value = t('Passwords do not match.', '两次密码不一致。')
    return
  }

  loading.value = true
  try {
    const res = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: form.name, email: form.email, phone: form.phone, password: form.password }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Registration failed')
    success.value = t('Registration successful! Redirecting...', '注册成功！正在跳转...')
    setTimeout(() => router.push('/login'), 1500)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>
