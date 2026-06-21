<template>
  <section class="auth-page">
    <div class="auth-card wide">
      <h1>{{ t('Apply for Admission', '报名申请') }}</h1>
      <p class="auth-subtitle">{{ t('Fill out the form below to begin the application process.', '请填写以下信息提交报名申请。') }}</p>

      <form @submit.prevent="handleApply" class="auth-form">
        <div class="form-row">
          <div class="form-group">
            <label>{{ t('Student Name', '学生姓名') }}</label>
            <input v-model="form.studentName" type="text" required :placeholder="t('Enter student name', '请输入学生姓名')" />
          </div>
          <div class="form-group">
            <label>{{ t('Date of Birth', '出生日期') }}</label>
            <input v-model="form.birthDate" type="date" required />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>{{ t('Gender', '性别') }}</label>
            <select v-model="form.gender" required>
              <option value="">{{ t('Select', '请选择') }}</option>
              <option value="male">{{ t('Male', '男') }}</option>
              <option value="female">{{ t('Female', '女') }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>{{ t('Applying for Grade', '申请年级') }}</label>
            <select v-model="form.grade" required>
              <option value="">{{ t('Select grade', '请选择年级') }}</option>
              <option v-for="g in 12" :key="g" :value="g">{{ t(`Grade ${g}`, `${g}年级`) }}</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>{{ t('Parent/Guardian Name', '家长姓名') }}</label>
            <input v-model="form.parentName" type="text" required :placeholder="t('Enter parent name', '请输入家长姓名')" />
          </div>
          <div class="form-group">
            <label>{{ t('Contact Phone', '联系电话') }}</label>
            <input v-model="form.phone" type="tel" required :placeholder="t('Enter phone number', '请输入联系电话')" />
          </div>
        </div>

        <div class="form-group">
          <label>{{ t('Contact Email', '联系邮箱') }}</label>
          <input v-model="form.email" type="email" required :placeholder="t('Enter email', '请输入邮箱')" />
        </div>

        <div class="form-group">
          <label>{{ t('Current School', '目前就读学校') }}</label>
          <input v-model="form.currentSchool" type="text" :placeholder="t('Enter current school name', '请输入目前就读学校')" />
        </div>

        <div class="form-group">
          <label>{{ t('Additional Notes', '备注信息') }}</label>
          <textarea v-model="form.notes" rows="4" :placeholder="t('Any additional information you would like to share', '请填写其他需要补充的信息')"></textarea>
        </div>

        <p v-if="error" class="form-error">{{ error }}</p>
        <p v-if="success" class="form-success">{{ success }}</p>

        <button type="submit" class="btn btn-primary form-submit" :disabled="loading">
          {{ loading ? t('Submitting...', '提交中...') : t('Submit Application', '提交报名') }}
        </button>
      </form>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useLanguage } from '../composables/useLanguage'

const { t } = useLanguage()

const form = reactive({
  studentName: '', birthDate: '', gender: '', grade: '',
  parentName: '', phone: '', email: '', currentSchool: '', notes: '',
})
const error = ref('')
const success = ref('')
const loading = ref(false)

async function handleApply() {
  error.value = ''
  loading.value = true
  try {
    const token = localStorage.getItem('xinhang-token')
    const res = await fetch('/api/apply', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: JSON.stringify(form),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Submission failed')
    success.value = t('Application submitted successfully! We will contact you soon.', '报名申请已提交成功！我们将尽快与您联系。')
    Object.keys(form).forEach(k => form[k] = '')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>
