<template>
  <div ref="root">
    <section class="page-hero" :style="{ '--page-hero-image': 'url(/assets/campus-blossom.jpg)' }">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('Account Registration', '账户注册') }}</p>
        <h1>{{ t('Create your account', '创建您的账户') }}</h1>
        <p>{{ t('Register to manage applications, track admission status, and access the parent portal.', '注册账户以管理报名信息、跟踪录取状态，并访问家长门户。') }}</p>
      </div>
    </section>

    <section class="auth-workspace">
      <form class="auth-form reveal" @submit.prevent="handleRegister">
        <div class="panel-heading">
          <p class="eyebrow">{{ t('New Account', '新账户') }}</p>
          <h2>{{ t('Registration', '注册') }}</h2>
        </div>

        <div v-if="msg" :class="['auth-msg', msgType]">{{ msg }}</div>

        <div class="form-section">
          <div class="form-grid">
            <label class="field"><span>{{ t('Full Name', '姓名') }}</span><input v-model="form.name" type="text" :placeholder="t('Your full name', '请输入姓名')" required /></label>
            <label class="field"><span>{{ t('Phone', '手机号') }}</span><input v-model="form.phone" type="tel" :placeholder="t('Mobile phone', '请输入手机号')" required /></label>
            <label class="field wide">
              <span>{{ t('Email', '邮箱') }}</span>
              <div class="code-row">
                <input v-model="form.email" type="email" :placeholder="t('Email address', '请输入邮箱')" required />
                <button class="btn btn-secondary" type="button" @click="sendCode" :disabled="codeCd > 0">
                  {{ codeCd > 0 ? `${codeCd}s` : t('Send Code', '发送验证码') }}
                </button>
              </div>
            </label>
            <label class="field"><span>{{ t('Verification Code', '验证码') }}</span><input v-model="form.code" type="text" :placeholder="t('6-digit code', '6位验证码')" required maxlength="6" /></label>
            <label class="field"><span>{{ t('Password', '密码') }}</span><input v-model="form.password" type="password" :placeholder="t('At least 6 characters', '至少6位字符')" required minlength="6" /></label>
          </div>
        </div>

        <div class="form-actions">
          <button class="btn btn-primary" type="submit" :disabled="loading">{{ loading ? t('Registering...', '注册中...') : t('Create Account', '注册') }}</button>
        </div>
        <p class="auth-links">{{ t('Already have an account?', '已有账户？') }} <router-link to="/login">{{ t('Sign in', '去登录') }}</router-link></p>
      </form>

      <aside class="auth-side reveal">
        <h3>{{ t('Why Register?', '为什么要注册？') }}</h3>
        <p>{{ t('Registration is required before applying. Register once and manage all applications from your account.', '报名前需要先注册账户。注册后可在一个账户中管理所有报名信息。') }}</p>
        <p>{{ t('Track your application status and receive updates directly.', '跟踪报名状态，直接接收最新通知。') }}</p>
        <p>{{ t('Access admission results and entrance permit downloads.', '查询录取结果，下载电子准考证。') }}</p>
      </aside>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useLanguage } from '../composables/useLanguage.js'
import { useReveal } from '../composables/useReveal.js'
import { useAuth } from '../composables/useAuth.js'

const { t } = useLanguage()
const { setAuth } = useAuth()
const router = useRouter()
const route = useRoute()
const root = ref(null)
useReveal(root)

const form = ref({ name: '', email: '', phone: '', password: '', code: '' })
const msg = ref('')
const msgType = ref('success')
const loading = ref(false)
const codeCd = ref(0)
let cdTimer = null

async function sendCode() {
  if (!form.value.email) { msg.value = t('Please enter email first', '请先输入邮箱'); msgType.value = 'error'; return }
  try {
    const res = await fetch('/api/send-code', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: form.value.email })
    })
    const data = await res.json()
    if (res.ok) {
      msgType.value = 'success'
      msg.value = t('Verification code sent to your email', '验证码已发送到您的邮箱')
      codeCd.value = 60
      cdTimer = setInterval(() => { codeCd.value--; if (codeCd.value <= 0) clearInterval(cdTimer) }, 1000)
    } else {
      msgType.value = 'error'
      msg.value = data.message || t('Failed to send code', '发送失败')
    }
  } catch {
    msgType.value = 'error'
    msg.value = t('Network error', '网络错误')
  }
}

async function handleRegister() {
  loading.value = true
  msg.value = ''
  try {
    const res = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    })
    const data = await res.json()
    if (res.ok) {
      setAuth(data.token, data.user)
      const redirect = route.query.redirect || '/profile'
      router.push(redirect)
    } else {
      msgType.value = 'error'
      msg.value = data.message || t('Registration failed', '注册失败')
    }
  } catch {
    msgType.value = 'error'
    msg.value = t('Network error', '网络错误')
  }
  loading.value = false
}
</script>
