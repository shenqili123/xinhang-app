<template>
  <div ref="root">
    <section class="page-hero" :style="{ '--page-hero-image': 'url(/assets/campus-overview.jpg)' }">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('Parent Portal', '家长门户') }}</p>
        <h1>{{ t('My Account', '个人中心') }}</h1>
        <p>{{ t('Manage your profile, track applications, and stay updated with the latest information.', '管理个人信息、跟踪报名状态、获取最新资讯。') }}</p>
      </div>
    </section>

    <section class="profile-workspace">
      <!-- 个人信息 -->
      <div class="profile-card reveal">
        <div class="panel-heading">
          <p class="eyebrow">{{ t('Account Information', '账户信息') }}</p>
          <h2>{{ t('My Profile', '个人资料') }}</h2>
        </div>

        <div v-if="profileMsg" :class="['auth-msg', profileMsgType]">{{ profileMsg }}</div>

        <div class="profile-info">
          <dl class="profile-fields">
            <div><dt>{{ t('Name', '姓名') }}</dt><dd>{{ user?.name || '-' }}</dd></div>
            <div><dt>{{ t('Email', '邮箱') }}</dt><dd>{{ user?.email || '-' }}</dd></div>
            <div><dt>{{ t('Phone', '手机号') }}</dt><dd>{{ user?.phone || '-' }}</dd></div>
            <div>
              <dt>{{ t('Email Status', '邮箱状态') }}</dt>
              <dd>
                <span :class="['profile-badge', user?.emailVerified ? 'verified' : 'unverified']">
                  {{ user?.emailVerified ? t('Verified', '已验证') : t('Not Verified', '未验证') }}
                </span>
              </dd>
            </div>
            <div><dt>{{ t('Registered', '注册时间') }}</dt><dd>{{ formatDate(user?.createdAt) }}</dd></div>
          </dl>
        </div>

        <form class="profile-edit" @submit.prevent="updateProfile" v-if="editing">
          <div class="form-section">
            <div class="form-grid">
              <label class="field"><span>{{ t('Name', '姓名') }}</span><input v-model="editForm.name" type="text" required /></label>
              <label class="field"><span>{{ t('Phone', '手机号') }}</span><input v-model="editForm.phone" type="tel" required /></label>
            </div>
          </div>
          <div class="form-actions">
            <button class="btn btn-primary" type="submit" :disabled="saving">{{ saving ? t('Saving...', '保存中...') : t('Save Changes', '保存修改') }}</button>
            <button class="btn btn-light" type="button" @click="editing = false">{{ t('Cancel', '取消') }}</button>
          </div>
        </form>

        <div class="form-actions" v-else>
          <button class="btn btn-secondary" @click="startEdit">{{ t('Edit Profile', '编辑资料') }}</button>
          <button class="btn btn-light" @click="handleLogout">{{ t('Sign Out', '退出登录') }}</button>
        </div>
      </div>

      <!-- 我的报名 -->
      <div class="applications-card reveal">
        <div class="panel-heading">
          <p class="eyebrow">{{ t('Application Records', '报名记录') }}</p>
          <h2>{{ t('My Applications', '我的报名') }}</h2>
        </div>

        <div v-if="loadingApps" class="loading-state">{{ t('Loading...', '加载中...') }}</div>

        <div v-else-if="applications.length === 0" class="empty-state">
          <p>{{ t('No applications yet. Start your journey with Xinhang!', '暂无报名记录，开启您的新航之旅吧！') }}</p>
          <router-link class="btn btn-primary" to="/apply">{{ t('Apply Now', '立即报名') }}</router-link>
        </div>

        <div v-else class="app-list">
          <div class="app-item" v-for="app in applications" :key="app.id">
            <div class="app-item-header">
              <div>
                <strong>{{ app.studentName }}</strong>
                <span class="app-permit">{{ app.permitNo }}</span>
              </div>
              <span :class="['app-status', 'status-' + app.status]">
                {{ statusLabel(app.status) }}
              </span>
            </div>
            <dl class="app-item-details">
              <div><dt>{{ t('Grade', '报名年级') }}</dt><dd>{{ gradeLabel(app.grade) }}</dd></div>
              <div><dt>{{ t('Track', '课程方向') }}</dt><dd>{{ app.track }}</dd></div>
              <div><dt>{{ t('School', '在读学校') }}</dt><dd>{{ app.currentSchool }}</dd></div>
              <div><dt>{{ t('Submitted', '提交时间') }}</dt><dd>{{ formatDate(app.createdAt) }}</dd></div>
            </dl>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLanguage } from '../composables/useLanguage.js'
import { useReveal } from '../composables/useReveal.js'
import { useAuth } from '../composables/useAuth.js'

const { t } = useLanguage()
const { user, authHeader, logout, setAuth } = useAuth()
const router = useRouter()
const root = ref(null)
useReveal(root)

const editing = ref(false)
const saving = ref(false)
const editForm = ref({ name: '', phone: '' })
const profileMsg = ref('')
const profileMsgType = ref('success')

const applications = ref([])
const loadingApps = ref(true)

async function fetchProfile() {
  try {
    const res = await fetch('/api/profile', { headers: authHeader() })
    if (res.ok) {
      const data = await res.json()
      if (data.user) setAuth(undefined, data.user)
    } else if (res.status === 401) {
      logout()
      router.push('/login')
    }
  } catch { /* use cached data */ }
}

function startEdit() {
  editForm.value = { name: user.value?.name || '', phone: user.value?.phone || '' }
  editing.value = true
  profileMsg.value = ''
}

async function updateProfile() {
  saving.value = true
  profileMsg.value = ''
  try {
    const res = await fetch('/api/profile', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...authHeader() },
      body: JSON.stringify(editForm.value)
    })
    const data = await res.json()
    if (res.ok) {
      setAuth(undefined, data.user)
      profileMsgType.value = 'success'
      profileMsg.value = t('Profile updated successfully', '资料更新成功')
      editing.value = false
    } else {
      profileMsgType.value = 'error'
      profileMsg.value = data.message || t('Update failed', '更新失败')
    }
  } catch {
    profileMsgType.value = 'error'
    profileMsg.value = t('Network error', '网络错误')
  }
  saving.value = false
}

async function fetchApplications() {
  loadingApps.value = true
  try {
    const res = await fetch('/api/my-applications', { headers: authHeader() })
    if (res.ok) {
      const data = await res.json()
      applications.value = data.data || []
    }
  } catch { /* ignore */ }
  loadingApps.value = false
}

function handleLogout() {
  logout()
  router.push('/')
}

function statusLabel(s) {
  const map = { pending: t('Pending Review', '待审核'), approved: t('Approved', '已录取'), rejected: t('Not Admitted', '未录取'), exam: t('Exam Scheduled', '待考试') }
  return map[s] || s
}

function gradeLabel(g) {
  if (g <= 6) return t(`Grade ${g}`, `${g}年级`)
  if (g <= 9) return t(`Grade ${g}`, `初${g - 6}`)
  return t(`Grade ${g}`, `高${g - 9}`)
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

onMounted(() => {
  fetchProfile()
  fetchApplications()
})
</script>

<style scoped>
.profile-workspace {
  display: grid;
  grid-template-columns: 420px minmax(0, 1fr);
  gap: 22px;
  align-items: start;
  padding: 72px;
  background: var(--ivory);
}

.profile-card,
.applications-card {
  background: #fff;
  box-shadow: 0 16px 46px rgba(7, 27, 55, .08);
  padding: 34px;
}

.profile-card {
  border-top: 5px solid var(--signal);
}

.applications-card {
  border-top: 5px solid var(--blue);
}

.profile-info {
  margin: 28px 0;
  padding: 0;
}

.profile-fields {
  display: grid;
  gap: 1px;
  margin: 0;
  background: rgba(7, 27, 55, .08);
}

.profile-fields div {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: 14px;
  padding: 16px 18px;
  background: var(--ivory);
}

.profile-fields dt {
  color: var(--muted);
  font-size: 13px;
  font-weight: 900;
  text-transform: uppercase;
}

.profile-fields dd {
  margin: 0;
  color: var(--navy);
  font-size: 15px;
  font-weight: 700;
  word-break: break-word;
}

.profile-badge {
  display: inline-block;
  padding: 3px 10px;
  font-size: 12px;
  font-weight: 900;
  text-transform: uppercase;
}

.profile-badge.verified {
  color: #0d6e3b;
  background: #e6f7ee;
}

.profile-badge.unverified {
  color: #8c1d18;
  background: #fde8e6;
}

.profile-edit .form-section {
  margin-top: 0;
  padding-top: 0;
  border-top: 0;
}

.loading-state,
.empty-state {
  padding: 48px 0;
  text-align: center;
  color: var(--muted);
  font-size: 16px;
}

.empty-state p {
  margin: 0 0 22px;
}

.app-list {
  display: grid;
  gap: 18px;
  margin-top: 22px;
}

.app-item {
  padding: 22px;
  border: 1px solid rgba(7, 27, 55, .1);
  transition: box-shadow .2s ease;
}

.app-item:hover {
  box-shadow: 0 8px 26px rgba(7, 27, 55, .08);
}

.app-item-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 16px;
  padding-bottom: 14px;
  border-bottom: 1px solid rgba(7, 27, 55, .08);
}

.app-item-header strong {
  display: block;
  color: var(--navy);
  font-size: 18px;
}

.app-permit {
  display: block;
  margin-top: 4px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 700;
  font-family: monospace;
}

.app-status {
  flex-shrink: 0;
  padding: 4px 12px;
  font-size: 12px;
  font-weight: 900;
  text-transform: uppercase;
}

.status-pending {
  color: #7c5c00;
  background: #fff8e1;
}

.status-approved {
  color: #0d6e3b;
  background: #e6f7ee;
}

.status-rejected {
  color: #8c1d18;
  background: #fde8e6;
}

.status-exam {
  color: #0e4a7a;
  background: #e3f2fd;
}

.app-item-details {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1px;
  margin: 0;
  background: rgba(7, 27, 55, .06);
}

.app-item-details div {
  padding: 10px 12px;
  background: #fff;
}

.app-item-details dt {
  margin-bottom: 4px;
  color: var(--muted);
  font-size: 11px;
  font-weight: 900;
  text-transform: uppercase;
}

.app-item-details dd {
  margin: 0;
  color: var(--navy);
  font-size: 14px;
  font-weight: 700;
}

@media (max-width: 1180px) {
  .profile-workspace {
    grid-template-columns: 1fr;
    padding: 52px;
  }
}

@media (max-width: 820px) {
  .profile-workspace {
    padding: 36px 24px;
  }

  .profile-card,
  .applications-card {
    padding: 26px;
  }

  .profile-fields div {
    grid-template-columns: 1fr;
    gap: 4px;
  }

  .app-item-details {
    grid-template-columns: 1fr;
  }
}
</style>
