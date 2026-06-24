<template>
  <main class="verify-page">
    <section class="verify-card">
      <img src="/assets/school-name-blue.png" alt="新航" class="verify-logo" />
      <h2 class="verify-title">{{ t('Permit Verification', '准考证验证') }}</h2>
      <p class="verify-desc">{{ t('Staff only. Enter the verification PIN and the QR code content to check a permit.', '仅限教职工使用。请输入验证密码和二维码内容来核查准考证。') }}</p>

      <form v-if="!result" @submit.prevent="handleVerify" class="verify-form">
        <label class="field">
          <span>{{ t('Verification PIN', '验证密码') }}</span>
          <input v-model="pin" type="password" :placeholder="t('Enter staff PIN', '请输入教职工验证密码')" required />
        </label>
        <label class="field">
          <span>{{ t('QR Code Content', '二维码内容') }}</span>
          <textarea v-model="code" rows="3" :placeholder="t('Paste the scanned QR code text here', '将扫描到的二维码文字粘贴到此处')" required></textarea>
        </label>
        <button class="btn btn-primary" type="submit" :disabled="loading">
          {{ loading ? t('Verifying...', '验证中...') : t('Verify', '验证') }}
        </button>
      </form>

      <div v-if="result && result.valid" class="verify-result success">
        <div class="result-icon ok">&#10003;</div>
        <h3>{{ t('Valid Permit', '准考证有效') }}</h3>
        <p class="result-sub">{{ t('Issued by Xinhang International School', '由新航实验国际学校官方签发') }}</p>
        <dl class="result-details">
          <div><dt>{{ t('Permit No.', '准考证号') }}</dt><dd>{{ result.appNo }}</dd></div>
          <div><dt>{{ t('Student', '学生姓名') }}</dt><dd>{{ result.student }}</dd></div>
        </dl>
        <button class="btn btn-secondary" @click="reset">{{ t('Verify Another', '验证下一个') }}</button>
      </div>

      <div v-if="result && !result.valid" class="verify-result fail">
        <div class="result-icon bad">&#10007;</div>
        <h3>{{ t('Verification Failed', '验证失败') }}</h3>
        <p class="result-sub">{{ result.message }}</p>
        <button class="btn btn-secondary" @click="reset">{{ t('Try Again', '重试') }}</button>
      </div>
    </section>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import { useLanguage } from '../composables/useLanguage.js'

const { t } = useLanguage()

const pin = ref('')
const code = ref('')
const loading = ref(false)
const result = ref(null)

async function handleVerify() {
  loading.value = true
  try {
    const res = await fetch('/api/verify-qr', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: code.value.trim(), pin: pin.value })
    })
    result.value = await res.json()
  } catch {
    result.value = { valid: false, message: '网络错误，请稍后重试' }
  }
  loading.value = false
}

function reset() {
  result.value = null
  code.value = ''
}
</script>

<style scoped>
.verify-page {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: linear-gradient(135deg, #f0f4f8 0%, #e2e8f0 100%);
}
.verify-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0,0,0,.1);
  padding: 2.5rem 2rem;
  max-width: 480px;
  width: 100%;
  text-align: center;
}
.verify-logo { height: 44px; margin-bottom: 1.2rem; }
.verify-title { margin: 0 0 .4rem; font-size: 1.4rem; color: #1e293b; }
.verify-desc { color: #64748b; font-size: .9rem; margin: 0 0 1.5rem; }
.verify-form { text-align: left; }
.field { display: block; margin-bottom: 1rem; }
.field span { display: block; font-size: .85rem; font-weight: 600; color: #334155; margin-bottom: .3rem; }
.field input, .field textarea {
  width: 100%; padding: .65rem .8rem;
  border: 1px solid #cbd5e1; border-radius: 8px;
  font-size: .95rem; font-family: inherit;
  box-sizing: border-box;
}
.field textarea { resize: vertical; font-family: monospace; font-size: .85rem; }
.btn { display: inline-block; padding: .7rem 1.6rem; border-radius: 8px; border: none; cursor: pointer; font-size: .95rem; font-weight: 600; }
.btn-primary { background: #1a73e8; color: #fff; width: 100%; }
.btn-primary:disabled { opacity: .6; cursor: not-allowed; }
.btn-secondary { background: #e2e8f0; color: #334155; margin-top: 1rem; }
.verify-result { margin-top: .5rem; }
.result-icon {
  width: 64px; height: 64px; border-radius: 50%;
  display: inline-flex; align-items: center; justify-content: center;
  font-size: 32px; font-weight: bold; color: #fff;
  margin-bottom: .8rem;
}
.ok { background: #22c55e; }
.bad { background: #ef4444; }
.verify-result h3 { margin: 0 0 .3rem; font-size: 1.3rem; color: #1e293b; }
.result-sub { color: #64748b; margin: 0 0 1rem; font-size: .9rem; }
.result-details {
  text-align: left; background: #f8fafc;
  border-radius: 10px; padding: 1rem 1.2rem; margin-bottom: .5rem;
}
.result-details div {
  display: flex; justify-content: space-between;
  padding: .5rem 0; border-bottom: 1px solid #e2e8f0;
}
.result-details div:last-child { border-bottom: none; }
.result-details dt { color: #64748b; font-size: .9rem; }
.result-details dd { margin: 0; font-weight: 600; color: #1e293b; }
</style>
