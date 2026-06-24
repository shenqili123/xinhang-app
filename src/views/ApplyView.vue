<template>
  <div ref="root">
    <section class="application-hero">
      <div class="application-hero-copy reveal">
        <p class="school-label">{{ t('Admissions Application', '招生报名') }}</p>
        <h1>{{ t('Apply online. Keep every next step clear.', '在线填报，后续安排一目了然。') }}</h1>
        <p>{{ t('Families submit student information once, then use the same entrance to receive exam details, generate an entrance permit, and check admission results.', '家庭一次性提交学生基础信息，并通过同一入口接收考试安排、生成电子准考证、查询录取相关结果。') }}</p>
      </div>
      <dl class="application-quickfacts reveal">
        <div><dt>{{ t('Enrollment', '招生年级') }}</dt><dd>{{ t('K-12 pathway', '小学至高中') }}</dd></div>
        <div><dt>{{ t('Format', '报名方式') }}</dt><dd>{{ t('Online first', '线上预报名') }}</dd></div>
        <div><dt>{{ t('Next Step', '后续安排') }}</dt><dd>{{ t('Exam notice', '考试通知') }}</dd></div>
      </dl>
    </section>

    <!-- ============ 提交成功 ============ -->
    <section v-if="submitted" class="success-section">
      <h2 class="success-title">{{ t('Application Submitted Successfully!', '报名信息提交成功！') }}</h2>
      <p class="success-subtitle">{{ t('Your application number is', '您的报名号为') }}：<strong>{{ receiptPermitNo }}</strong></p>

      <div class="permit-card-wrapper" ref="permitCard">
        <div class="permit-card">
          <div class="permit-head">
            <img src="/assets/school-name-blue.png" alt="山东新航实验国际学校" />
            <span>{{ t('Entrance Permit', '电子准考证') }}</span>
          </div>
          <div class="permit-title">
            <p>{{ t('2026 Admission Assessment', '2026 入学综合评估') }}</p>
            <strong>{{ receiptPermitNo }}</strong>
          </div>
          <dl class="permit-details">
            <div><dt>{{ t('Student', '学生') }}</dt><dd>{{ submitSnapshot.studentName }}</dd></div>
            <div><dt>{{ t('Grade', '报考年级') }}</dt><dd>{{ submitSnapshot.gradeLabel }}</dd></div>
            <div><dt>{{ t('Current School', '现就读学校') }}</dt><dd>{{ submitSnapshot.currentSchool }}</dd></div>
            <div><dt>{{ t('Exam Time', '考试时间') }}</dt><dd>{{ t('To be announced by SMS', '以短信通知为准') }}</dd></div>
            <div><dt>{{ t('Venue', '考试地点') }}</dt><dd>{{ t('Xinhang Campus', '新航校园') }}</dd></div>
          </dl>
          <div class="permit-bottom">
            <img v-if="qrSrc" :src="qrSrc" alt="QR Code" class="permit-qr" crossorigin="anonymous" />
            <p>{{ t('Bring this permit and valid ID on assessment day.', '考试当天请携带准考证及有效身份证件。') }}</p>
          </div>
        </div>
      </div>

      <div class="permit-actions">
        <button class="btn btn-primary" @click="downloadPDF">{{ t('Download PDF', '下载 PDF') }}</button>
        <button class="btn btn-secondary" @click="resetForm">{{ t('OK, Back to Form', '确定，返回报名') }}</button>
      </div>
    </section>

    <!-- ============ 表单 ============ -->
    <section v-else class="application-workspace" id="application">
      <aside class="application-rail reveal">
        <p class="eyebrow">{{ t('Process', '流程') }}</p>
        <ol>
          <li><span>01</span><strong>{{ t('Fill Information', '信息填报') }}</strong><p>{{ t('Student and guardian details.', '填写学生与监护人基础信息。') }}</p></li>
          <li><span>02</span><strong>{{ t('Confirm Exam', '确认考试') }}</strong><p>{{ t('Receive date, venue, and reminders.', '获取考试时间、地点与注意事项。') }}</p></li>
          <li><span>03</span><strong>{{ t('Entrance Permit', '电子准考证') }}</strong><p>{{ t('Generate and download before exam day.', '考前生成并下载电子准考证。') }}</p></li>
          <li><span>04</span><strong>{{ t('Results', '成绩查询') }}</strong><p>{{ t('Return after score release.', '成绩发布后回到入口查询。') }}</p></li>
        </ol>
      </aside>

      <form class="application-form reveal" @submit.prevent="submitForm">
        <div class="panel-heading">
          <p class="eyebrow">{{ t('Student Application', '学生报名信息') }}</p>
          <h2>{{ t('Basic information', '基础信息填报') }}</h2>
          <p>{{ t('All fields except Family Notes are required.', '除家庭补充说明外，所有字段均为必填。') }}</p>
        </div>

        <div v-if="msg" :class="['auth-msg', msgType]" ref="msgBox">{{ msg }}</div>

        <div class="form-section">
          <h3>{{ t('Student', '学生信息') }}</h3>
          <div class="form-grid">
            <label class="field"><span>{{ t('Student Name', '学生姓名') }} *</span><input v-model="form.studentName" type="text" :placeholder="t('Chen Ming', '请输入学生姓名')" required /></label>
            <label class="field"><span>{{ t('Applying Grade', '报考年级') }} *</span>
              <select v-model="form.grade" required>
                <option value="1">{{ t('Grade 1', '一年级') }}</option>
                <option value="7">{{ t('Grade 7', '七年级') }}</option>
                <option value="10">{{ t('Grade 10', '高一') }}</option>
                <option value="0">{{ t('Transfer Student', '插班生') }}</option>
              </select>
            </label>
            <label class="field"><span>{{ t('Current School', '现就读学校') }} *</span><input v-model="form.currentSchool" type="text" :placeholder="t('Current school', '请输入现就读学校')" required /></label>
            <label class="field"><span>{{ t('ID / Passport No.', '证件号码') }} *</span><input v-model="form.idNumber" type="text" :placeholder="t('ID or passport number', '身份证号或护照号')" required /></label>
            <label class="field"><span>{{ t('Gender', '性别') }} *</span>
              <select v-model="form.gender" required>
                <option value="">{{ t('-- Select --', '-- 请选择 --') }}</option>
                <option value="Female">{{ t('Female', '女') }}</option>
                <option value="Male">{{ t('Male', '男') }}</option>
              </select>
            </label>
            <label class="field"><span>{{ t('Boarding Need', '寄宿需求') }} *</span>
              <select v-model="form.boardingNeed" required>
                <option value="boarding">{{ t('Need boarding', '需要寄宿') }}</option>
                <option value="day">{{ t('Day student', '走读') }}</option>
                <option value="tbd">{{ t('To be discussed', '待沟通') }}</option>
              </select>
            </label>
          </div>
        </div>

        <div class="form-section">
          <h3>{{ t('Guardian', '监护人信息') }}</h3>
          <div class="form-grid">
            <label class="field"><span>{{ t('Guardian Name', '监护人姓名') }} *</span><input v-model="form.parentName" type="text" :placeholder="t('Parent or guardian', '请输入监护人姓名')" required /></label>
            <label class="field"><span>{{ t('Mobile Phone', '联系电话') }} *</span><input v-model="form.phone" type="tel" :placeholder="t('Mobile phone', '请输入联系电话')" required /></label>
            <label class="field"><span>{{ t('Relationship', '与学生关系') }} *</span>
              <select v-model="form.relationship" required>
                <option value="Parent">{{ t('Parent', '父母') }}</option>
                <option value="Guardian">{{ t('Guardian', '监护人') }}</option>
                <option value="Other">{{ t('Other Family', '其他亲属') }}</option>
              </select>
            </label>
            <label class="field"><span>{{ t('Email / WeChat', '邮箱 / 微信') }} *</span><input v-model="form.email" type="text" :placeholder="t('Email or WeChat', '邮箱或微信号')" required /></label>
          </div>
        </div>

        <div class="form-section">
          <h3>{{ t('Pathway', '意向方向') }}</h3>
          <div class="form-grid">
            <label class="field"><span>{{ t('Preferred Track', '意向方向') }} *</span>
              <select v-model="form.track" required>
                <option value="integrated">{{ t('Domestic & international pathway', '国内与国际融合方向') }}</option>
                <option value="academic">{{ t('Academic foundation', '学术基础方向') }}</option>
                <option value="boarding">{{ t('Boarding growth', '寄宿成长方向') }}</option>
              </select>
            </label>
            <label class="field"><span>{{ t('Preferred Visit Date', '意向探校日期') }} *</span><input v-model="form.visitDate" type="date" required /></label>
            <label class="field wide"><span>{{ t('Family Notes (optional)', '家庭补充说明（选填）') }}</span><textarea v-model="form.notes" rows="4" :placeholder="t('Learning background, boarding needs, or questions for admission office', '可填写学习背景、寄宿需求或想咨询的问题')"></textarea></label>
          </div>
        </div>

        <label class="consent-line">
          <input v-model="consent" type="checkbox" required />
          <span>{{ t('I confirm the information above is accurate and agree to be contacted by the admission office.', '我确认以上信息真实有效，并同意招生办公室与我联系。') }}</span>
        </label>

        <div class="form-actions">
          <button class="btn btn-primary" type="submit" :disabled="submitting">{{ submitting ? t('Submitting...', '提交中...') : t('Generate Application Receipt', '生成报名回执') }}</button>
          <router-link class="btn btn-light" to="/admission">{{ t('Back to Admission', '返回招生页') }}</router-link>
        </div>
      </form>
    </section>

    <!-- ============ 查询区域 ============ -->
    <section class="query-section">
      <article class="query-card reveal">
        <p class="eyebrow">{{ t('Entrance Permit', '准考证查询') }}</p>
        <h2>{{ t('Find or download the exam permit.', '查询或下载电子准考证。') }}</h2>
        <div class="query-row">
          <input v-model="queryPermitInput" type="text" :placeholder="t('Application number / mobile phone', '报名号 / 手机号')" />
          <button class="btn btn-secondary" type="button" @click="queryPermit" :disabled="querying">{{ t('Query Permit', '查询准考证') }}</button>
        </div>
        <div v-if="queryResult" :class="['query-result', queryResult.found ? 'ok' : 'fail']">
          <template v-if="queryResult.found">
            <p><strong>{{ t('Permit No.', '准考证号') }}:</strong> {{ queryResult.permitNo }}</p>
            <p><strong>{{ t('Student', '学生') }}:</strong> {{ queryResult.student }}</p>
            <p><strong>{{ t('Status', '状态') }}:</strong> {{ statusLabel(queryResult.status) }}</p>
          </template>
          <p v-else>{{ queryResult.message }}</p>
        </div>
      </article>
      <article class="query-card reveal">
        <p class="eyebrow">{{ t('Results', '成绩查询') }}</p>
        <h2>{{ t('Check assessment and admission status.', '查询考试与录取状态。') }}</h2>
        <div class="query-row">
          <input v-model="queryResultInput" type="text" :placeholder="t('Application number / mobile phone', '报名号 / 手机号')" />
          <button class="btn btn-primary" type="button" @click="queryResults" :disabled="queryingResults">{{ t('Check Results', '查询成绩') }}</button>
        </div>
        <div v-if="resultQueryResult" :class="['query-result', resultQueryResult.found ? 'ok' : 'fail']">
          <template v-if="resultQueryResult.found">
            <p><strong>{{ t('Student', '学生') }}:</strong> {{ resultQueryResult.student }}</p>
            <p><strong>{{ t('Status', '状态') }}:</strong> {{ statusLabel(resultQueryResult.status) }}</p>
          </template>
          <p v-else>{{ resultQueryResult.message }}</p>
        </div>
      </article>
    </section>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useLanguage } from '../composables/useLanguage.js'
import { useReveal } from '../composables/useReveal.js'
import { useAuth } from '../composables/useAuth.js'
import html2canvas from 'html2canvas'
import { jsPDF } from 'jspdf'

const { t } = useLanguage()
const { authHeader } = useAuth()
const router = useRouter()
const root = ref(null)
useReveal(root)

const form = ref({
  studentName: '', grade: '7', currentSchool: '', gender: '',
  idNumber: '', boardingNeed: 'boarding',
  parentName: '', phone: '', relationship: 'Parent', email: '',
  track: 'integrated', visitDate: '', notes: ''
})
const consent = ref(false)
const submitting = ref(false)
const submitted = ref(false)
const msg = ref('')
const msgType = ref('success')
const receiptPermitNo = ref('')
const qrSrc = ref('')
const submitSnapshot = ref({})
const permitCard = ref(null)
const msgBox = ref(null)

const gradeMap = { '1': 'Grade 1 / 一年级', '7': 'Grade 7 / 七年级', '10': 'Grade 10 / 高一', '0': 'Transfer / 插班' }

function showError(text) {
  msgType.value = 'error'
  msg.value = text
  setTimeout(() => msgBox.value?.scrollIntoView({ behavior: 'smooth', block: 'center' }), 50)
}

const requiredFields = [
  { key: 'studentName', label: () => t('Student Name', '学生姓名') },
  { key: 'gender', label: () => t('Gender', '性别') },
  { key: 'currentSchool', label: () => t('Current School', '现就读学校') },
  { key: 'idNumber', label: () => t('ID / Passport No.', '证件号码') },
  { key: 'parentName', label: () => t('Guardian Name', '监护人姓名') },
  { key: 'phone', label: () => t('Mobile Phone', '联系电话'), min: 8 },
  { key: 'email', label: () => t('Email / WeChat', '邮箱 / 微信') },
  { key: 'visitDate', label: () => t('Preferred Visit Date', '意向探校日期') },
]

async function submitForm() {
  const f = form.value
  for (const r of requiredFields) {
    const val = (f[r.key] || '').trim()
    if (!val) { showError(t(`Please fill in: ${r.label()}`, `请填写：${r.label()}`)); return }
    if (r.min && val.length < r.min) { showError(t(`${r.label()} is too short`, `${r.label()} 长度不足`)); return }
  }
  if (!consent.value) { showError(t('Please check the consent box', '请勾选确认信息真实有效')); return }

  submitting.value = true
  msg.value = ''
  try {
    const payload = {
      studentName: f.studentName.trim(),
      gender: f.gender,
      grade: parseInt(f.grade) || 1,
      idNumber: f.idNumber.trim(),
      boardingNeed: f.boardingNeed,
      parentName: f.parentName.trim(),
      phone: f.phone.trim(),
      relationship: f.relationship,
      email: f.email.trim(),
      currentSchool: f.currentSchool.trim(),
      track: f.track,
      visitDate: f.visitDate,
      notes: f.notes.trim()
    }
    const res = await fetch('/api/apply', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...authHeader() },
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (res.ok) {
      receiptPermitNo.value = data.permitNo
      qrSrc.value = `/api/permit-qr?no=${encodeURIComponent(data.permitNo)}&student=${encodeURIComponent(f.studentName.trim())}`
      submitSnapshot.value = {
        studentName: f.studentName.trim(),
        gradeLabel: gradeMap[f.grade] || f.grade,
        currentSchool: f.currentSchool.trim()
      }
      submitted.value = true
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } else {
      showError(data.message || t('Submission failed', '提交失败'))
    }
  } catch {
    showError(t('Network error, please retry', '网络错误，请重试'))
  }
  submitting.value = false
}

function resetForm() {
  router.go(0)
}

async function downloadPDF() {
  const el = permitCard.value
  if (!el) return

  const imgs = el.querySelectorAll('img')
  await Promise.all([...imgs].map(img => {
    if (img.complete) return Promise.resolve()
    return new Promise(resolve => {
      img.onload = resolve
      img.onerror = resolve
    })
  }))
  await new Promise(r => setTimeout(r, 300))

  const canvas = await html2canvas(el, { scale: 2, useCORS: true, allowTaint: true })
  const pdf = new jsPDF('p', 'mm', 'a4')
  const imgW = 170
  const imgH = (canvas.height * imgW) / canvas.width
  pdf.addImage(canvas.toDataURL('image/png'), 'PNG', 20, 20, imgW, imgH)
  pdf.save(`准考证_${receiptPermitNo.value}.pdf`)
}

// 查询
const queryPermitInput = ref('')
const querying = ref(false)
const queryResult = ref(null)
async function queryPermit() {
  if (!queryPermitInput.value.trim()) return
  querying.value = true; queryResult.value = null
  try { const r = await fetch(`/api/query-permit?q=${encodeURIComponent(queryPermitInput.value.trim())}`); queryResult.value = await r.json() }
  catch { queryResult.value = { found: false, message: '网络错误' } }
  querying.value = false
}
const queryResultInput = ref('')
const queryingResults = ref(false)
const resultQueryResult = ref(null)
async function queryResults() {
  if (!queryResultInput.value.trim()) return
  queryingResults.value = true; resultQueryResult.value = null
  try { const r = await fetch(`/api/query-permit?q=${encodeURIComponent(queryResultInput.value.trim())}`); resultQueryResult.value = await r.json() }
  catch { resultQueryResult.value = { found: false, message: '网络错误' } }
  queryingResults.value = false
}
function statusLabel(s) {
  const map = { pending: t('Submitted, awaiting review', '已提交，待审核'), reviewed: t('Reviewed', '已审核'), admitted: t('Admitted', '已录取'), rejected: t('Not admitted', '未录取') }
  return map[s] || s
}
</script>

<style scoped>
.success-section {
  max-width: 700px;
  margin: 3rem auto;
  padding: 0 1.5rem;
  text-align: center;
}
.success-title {
  font-size: 1.6rem;
  color: #22c55e;
  margin: 0 0 .5rem;
}
.success-subtitle {
  font-size: 1.1rem;
  color: #334155;
  margin: 0 0 2rem;
}
.success-subtitle strong {
  font-size: 1.3rem;
  color: #1a73e8;
}
.permit-card-wrapper {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0,0,0,.12);
  overflow: hidden;
}
.permit-card-wrapper .permit-card {
  text-align: left;
}
.permit-qr {
  width: 110px;
  height: 110px;
}
.permit-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-top: 2rem;
}
.query-result {
  margin-top: 1rem; padding: 1rem; border-radius: 8px; font-size: .95rem;
}
.query-result.ok { background: #f0fdf4; border: 1px solid #bbf7d0; }
.query-result.fail { background: #fef2f2; border: 1px solid #fecaca; color: #991b1b; }
.query-result p { margin: .3rem 0; }
</style>
