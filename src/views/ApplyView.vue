<template>
  <div ref="root">
    <section class="application-hero">
      <div class="application-hero-copy reveal">
        <p class="school-label">{{ t('Admissions Application', '招生报名') }}</p>
        <h1>{{ t('Apply to Xinhang.', '新航在线报名。') }}</h1>
        <p>{{ t('Submit once, then use the same entrance for exam notices, permits and results.', '一次填报，同一入口查询考试通知、准考证与成绩。') }}</p>
      </div>
      <dl class="application-quickfacts reveal">
        <div><dt>{{ t('Enrollment', '招生年级') }}</dt><dd>{{ t('K-12 pathway', '小学至高中') }}</dd></div>
        <div><dt>{{ t('Format', '报名方式') }}</dt><dd>{{ t('Online first', '线上预报名') }}</dd></div>
        <div><dt>{{ t('Next Step', '后续安排') }}</dt><dd>{{ t('Exam notice', '考试通知') }}</dd></div>
      </dl>
    </section>

    <nav class="application-subnav" aria-label="Application sections">
      <a href="#application">{{ t('Application Form', '报名信息') }}</a>
      <a href="#permit">{{ t('Entrance Permit', '准考证') }}</a>
      <a href="#queries">{{ t('Results Query', '成绩查询') }}</a>
      <a href="#guide">{{ t('Application Guide', '报考指南') }}</a>
    </nav>

    <!-- Success: Permit Display -->
    <section v-if="submitted" class="permit-section">
      <h2 class="permit-section-title">{{ t('Application Submitted Successfully!', '报名信息提交成功！') }}</h2>
      <p class="permit-section-sub">{{ t('Your permit number is', '您的准考证号为') }}：<strong>{{ receiptPermitNo }}</strong></p>

      <div class="permit-wrapper" ref="permitCard">
        <div class="exam-permit">
          <div class="permit-header">
            <h2 class="permit-school-name">山东新航实验外国语学校</h2>
            <p class="permit-subtitle">{{ gradeExamTitle }}</p>
            <h1 class="permit-main-title">准 考 证</h1>
          </div>

          <div class="permit-body">
            <div class="permit-info-section">
              <div class="permit-photo-area">
                <img v-if="submitSnapshot.photo" :src="submitSnapshot.photo" alt="考生照片" />
                <div v-else class="photo-placeholder">照片</div>
              </div>
              <dl class="permit-info-fields">
                <div><dt>准考证号：</dt><dd>{{ receiptPermitNo }}</dd></div>
                <div><dt>姓&emsp;&emsp;名：</dt><dd>{{ submitSnapshot.studentName }}</dd></div>
                <div><dt>性&emsp;&emsp;别：</dt><dd>{{ submitSnapshot.gender === 'Male' ? '男' : '女' }}</dd></div>
                <div><dt>证件号码：</dt><dd>{{ maskIDNumber(submitSnapshot.idNumber) }}</dd></div>
                <div><dt>报考年级：</dt><dd>{{ submitSnapshot.division }} {{ submitSnapshot.gradeLabel }}</dd></div>
                <div><dt>考&emsp;&emsp;点：</dt><dd>山东新航实验外国语学校</dd></div>
                <div><dt>考场号：</dt><dd>{{ submitSnapshot.examRoom }}</dd></div>
                <div><dt>座位号：</dt><dd>{{ submitSnapshot.seatNumber }}</dd></div>
              </dl>
            </div>

            <div class="permit-barcode">
              <img v-if="qrSrc" :src="qrSrc" alt="QR Code" crossorigin="anonymous" />
            </div>

            <div class="permit-schedule">
              <table class="schedule-table">
                <thead>
                  <tr>
                    <th class="corner-cell">
                      <span class="top-label">日期</span>
                      <span class="bottom-label">科目</span>
                    </th>
                    <th>考试日</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td class="time-label">上午</td>
                    <td>
                      <div>笔试一</div>
                      <div class="time-hint">(以短信通知为准)</div>
                    </td>
                  </tr>
                  <tr>
                    <td class="time-label">上午</td>
                    <td>
                      <div>笔试二</div>
                      <div class="time-hint">(以短信通知为准)</div>
                    </td>
                  </tr>
                  <tr>
                    <td class="time-label">下午</td>
                    <td>
                      <div>面试一</div>
                      <div class="time-hint">(以短信通知为准)</div>
                    </td>
                  </tr>
                  <tr>
                    <td class="time-label">下午</td>
                    <td>
                      <div>面试二</div>
                      <div class="time-hint">(以短信通知为准)</div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <p class="permit-notice">注：考试当天请携带本准考证及有效身份证件，提前15分钟到达考场。</p>
          </div>
        </div>
      </div>

      <div class="permit-actions">
        <button class="btn btn-primary" @click="downloadPDF">{{ t('Download PDF', '下载 PDF') }}</button>
        <button class="btn btn-secondary" @click="resetForm">{{ t('OK, Back to Form', '确定，返回报名') }}</button>
      </div>
    </section>

    <!-- Form -->
    <section v-else class="application-workspace" id="application">
      <form class="application-form reveal" @submit.prevent="submitForm">
        <div class="panel-heading">
          <p class="eyebrow">{{ t('Student Application', '学生报名信息') }}</p>
          <h2>{{ t('Basic information', '基础信息填报') }}</h2>
          <p>{{ t('Please keep names and phone numbers consistent with future exam and result queries.', '请确保姓名与手机号准确，后续准考证和成绩查询将使用同一信息。') }}</p>
        </div>

        <div v-if="msg" :class="['auth-msg', msgType]" ref="msgBox">{{ msg }}</div>

        <div class="form-section">
          <h3>{{ t('Student', '学生信息') }}</h3>
          <div class="form-grid">
            <label class="field"><span>{{ t('Student Name', '学生姓名') }} *</span><input v-model="form.studentName" type="text" :placeholder="t('Chen Ming', '请输入学生姓名')" required /></label>
            <label class="field"><span>{{ t('Gender', '性别') }} *</span>
              <select v-model="form.gender" required>
                <option value="">{{ t('-- Select --', '-- 请选择 --') }}</option>
                <option value="Male">{{ t('Male', '男') }}</option>
                <option value="Female">{{ t('Female', '女') }}</option>
              </select>
            </label>
            <label class="field"><span>{{ t('Applying Division', '报考学部') }} *</span>
              <select v-model="form.division" required @change="onDivisionChange">
                <option value="">{{ t('-- Select --', '-- 请选择 --') }}</option>
                <option value="小学部">{{ t('Primary School', '小学部') }}</option>
                <option value="初中部">{{ t('Junior High', '初中部') }}</option>
                <option value="高中部">{{ t('Senior High', '高中部') }}</option>
              </select>
            </label>
            <label class="field"><span>{{ t('Applying Grade', '报考年级') }} *</span>
              <select v-model="form.grade" required :disabled="!form.division">
                <option value="">{{ t('-- Select division first --', '-- 请先选择学部 --') }}</option>
                <option v-for="g in gradeOptions" :key="g.value" :value="g.value">{{ g.label }}</option>
              </select>
            </label>
            <label class="field"><span>{{ t('Current School', '现就读学校') }} *</span><input v-model="form.currentSchool" type="text" :placeholder="t('Current school', '请输入现就读学校')" required /></label>
            <label class="field"><span>{{ t('ID Number (18 digits)', '身份证号（18位）') }} *</span>
              <input v-model="form.idNumber" type="text" maxlength="18" :placeholder="t('18-digit ID number', '请输入18位身份证号')" required @blur="validateID" />
              <span v-if="idError" class="field-error">{{ idError }}</span>
            </label>
            <label class="field"><span>{{ t('Boarding Need', '寄宿需求') }} *</span>
              <select v-model="form.boardingNeed" required>
                <option value="boarding">{{ t('Need boarding', '需要寄宿') }}</option>
                <option value="day">{{ t('Day student', '走读') }}</option>
                <option value="tbd">{{ t('To be discussed', '待沟通') }}</option>
              </select>
            </label>
            <div class="field wide photo-upload-field">
              <span>{{ t('Student Photo', '学生照片') }} *</span>
              <div class="photo-upload-area" @click="triggerPhotoInput" @dragover.prevent @drop.prevent="handleDrop">
                <img v-if="photoPreview" :src="photoPreview" class="photo-preview" />
                <div v-else class="photo-upload-placeholder">
                  <span class="upload-icon">📷</span>
                  <span>{{ t('Click or drag to upload', '点击或拖拽上传照片') }}</span>
                  <span class="upload-hint">{{ t('JPG/PNG, max 5MB', 'JPG/PNG格式，不超过5MB') }}</span>
                </div>
              </div>
              <input ref="photoInput" type="file" accept="image/jpeg,image/png" style="display:none" @change="handlePhotoSelect" />
              <span v-if="photoUploading" class="field-hint">{{ t('Uploading...', '上传中...') }}</span>
            </div>
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
          <button class="btn btn-primary" type="submit" :disabled="submitting || photoUploading">{{ submitting ? t('Submitting...', '提交中...') : t('Generate Application Receipt', '生成报名回执') }}</button>
        </div>
      </form>

    </section>

    <!-- Query section -->
    <section class="query-section" id="queries">
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

    <!-- Guide section -->
    <section class="portal-banner application-guide" id="guide">
      <div class="reveal">
        <p class="eyebrow">{{ t('Application Guide', '报考指南') }}</p>
        <h2>{{ t('Four steps from application to results.', '从报名到查询，四步完成。') }}</h2>
        <p>{{ t('For grade selection, boarding needs or curriculum questions, families may contact the admission office before submitting.', '如需确认报考年级、寄宿需求或课程方向，可在提交前联系招生老师。') }}</p>
      </div>
      <ol class="portal-steps reveal">
        <li><span>01</span><strong>{{ t('Apply Online', '在线报名') }}</strong></li>
        <li><span>02</span><strong>{{ t('Campus Visit', '预约探校') }}</strong></li>
        <li><span>03</span><strong>{{ t('Exam Notice', '考试通知') }}</strong></li>
        <li><span>04</span><strong>{{ t('Results', '成绩查询') }}</strong></li>
      </ol>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLanguage } from '../composables/useLanguage.js'
import { useReveal } from '../composables/useReveal.js'
import { useAuth } from '../composables/useAuth.js'
import html2canvas from 'html2canvas'
import { jsPDF } from 'jspdf'

const { t } = useLanguage()
const { user, authHeader, logout } = useAuth()
const router = useRouter()
const root = ref(null)
useReveal(root)

const form = ref({
  studentName: '', division: '', grade: '', currentSchool: '', gender: '',
  idNumber: '', boardingNeed: 'boarding', photo: '',
  parentName: '', phone: '', relationship: 'Parent', email: '',
  track: 'integrated', visitDate: '', notes: ''
})

const gradeOptions = computed(() => {
  switch (form.value.division) {
    case '小学部': return [
      { value: '1', label: t('Grade 1', '一年级') },
      { value: '2', label: t('Grade 2', '二年级') },
      { value: '3', label: t('Grade 3', '三年级') },
      { value: '4', label: t('Grade 4', '四年级') },
      { value: '5', label: t('Grade 5', '五年级') },
      { value: '6', label: t('Grade 6', '六年级') }
    ]
    case '初中部': return [
      { value: '7', label: t('Grade 7', '七年级') },
      { value: '8', label: t('Grade 8', '八年级') },
      { value: '9', label: t('Grade 9', '九年级') }
    ]
    case '高中部': return [
      { value: '10', label: t('Grade 10', '高一') },
      { value: '11', label: t('Grade 11', '高二') },
      { value: '12', label: t('Grade 12', '高三') }
    ]
    default: return []
  }
})

function onDivisionChange() {
  form.value.grade = ''
}

onMounted(() => {
  if (user.value) {
    form.value.parentName = user.value.name || ''
    form.value.phone = user.value.phone || ''
    form.value.email = user.value.email || ''
  }
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
const photoInput = ref(null)
const photoPreview = ref('')
const photoUploading = ref(false)
const idError = ref('')

const gradeMap = {
  '1': '一年级', '2': '二年级', '3': '三年级', '4': '四年级', '5': '五年级', '6': '六年级',
  '7': '七年级', '8': '八年级', '9': '九年级',
  '10': '高一', '11': '高二', '12': '高三',
  '0': '插班'
}

const gradeExamTitle = computed(() => {
  const g = submitSnapshot.value.grade
  const div = submitSnapshot.value.division || ''
  if (g >= 1 && g <= 6) return `2026年${div}${gradeMap[String(g)]}入学综合评估`
  if (g >= 7 && g <= 9) return `2026年${div}${gradeMap[String(g)]}入学综合评估`
  if (g >= 10 && g <= 12) return `2026年${div}${gradeMap[String(g)]}入学综合评估`
  return '2026年入学综合评估'
})

function validateIDNumber(id) {
  if (id.length !== 18) return '身份证号必须为18位'
  if (!/^\d{17}[\dXx]$/.test(id)) return '身份证号格式不正确'
  const weights = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
  const checkCodes = '10X98765432'
  let sum = 0
  for (let i = 0; i < 17; i++) sum += parseInt(id[i]) * weights[i]
  const expected = checkCodes[sum % 11]
  const last = id[17].toUpperCase()
  if (last !== expected) return '身份证号校验位不正确'
  return ''
}

function validateID() {
  const val = form.value.idNumber.trim()
  if (val.length > 0) {
    idError.value = validateIDNumber(val)
  } else {
    idError.value = ''
  }
}

function maskIDNumber(id) {
  if (!id || id.length < 10) return id
  return id.substring(0, 6) + '********' + id.substring(14)
}

function triggerPhotoInput() {
  photoInput.value?.click()
}

function handleDrop(e) {
  const file = e.dataTransfer?.files?.[0]
  if (file) uploadPhoto(file)
}

function handlePhotoSelect(e) {
  const file = e.target.files?.[0]
  if (file) uploadPhoto(file)
}

async function uploadPhoto(file) {
  if (!file.type.match(/^image\/(jpeg|png)$/)) {
    showError(t('Only JPG/PNG photos are supported', '仅支持JPG/PNG格式照片'))
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    showError(t('Photo must be under 5MB', '照片不能超过5MB'))
    return
  }

  photoUploading.value = true
  const formData = new FormData()
  formData.append('photo', file)

  try {
    const res = await fetch('/api/upload-photo', {
      method: 'POST',
      headers: { ...authHeader() },
      body: formData
    })
    const data = await res.json()
    if (res.ok) {
      form.value.photo = data.url
      photoPreview.value = data.url
    } else {
      showError(data.message || t('Upload failed', '上传失败'))
    }
  } catch {
    showError(t('Network error during upload', '上传网络错误'))
  }
  photoUploading.value = false
}

function showError(text) {
  msgType.value = 'error'
  msg.value = text
  setTimeout(() => msgBox.value?.scrollIntoView({ behavior: 'smooth', block: 'center' }), 50)
}

const requiredFields = [
  { key: 'studentName', label: () => t('Student Name', '学生姓名') },
  { key: 'gender', label: () => t('Gender', '性别') },
  { key: 'division', label: () => t('Applying Division', '报考学部') },
  { key: 'grade', label: () => t('Applying Grade', '报考年级') },
  { key: 'currentSchool', label: () => t('Current School', '现就读学校') },
  { key: 'idNumber', label: () => t('ID Number', '身份证号') },
  { key: 'photo', label: () => t('Student Photo', '学生照片') },
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

  const idErr = validateIDNumber(f.idNumber.trim())
  if (idErr) { showError(idErr); return }

  if (!consent.value) { showError(t('Please check the consent box', '请勾选确认信息真实有效')); return }

  submitting.value = true
  msg.value = ''
  try {
    const payload = {
      studentName: f.studentName.trim(),
      gender: f.gender,
      division: f.division,
      grade: parseInt(f.grade) || 1,
      idNumber: f.idNumber.trim(),
      photo: f.photo,
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
        gender: f.gender,
        division: f.division,
        grade: parseInt(f.grade) || 1,
        gradeLabel: gradeMap[f.grade] || f.grade,
        currentSchool: f.currentSchool.trim(),
        idNumber: f.idNumber.trim(),
        photo: f.photo,
        examRoom: data.examRoom,
        seatNumber: data.seatNumber,
      }
      submitted.value = true
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } else if (res.status === 401) {
      logout()
      router.push({ name: 'Login', query: { redirect: '/apply' } })
    } else {
      showError(data.message || t('Submission failed', '提交失败'))
    }
  } catch {
    showError(t('Network error, please retry', '网络错误，请重试'))
  }
  submitting.value = false
}

function resetForm() { router.go(0) }

async function downloadPDF() {
  const el = permitCard.value
  if (!el) return
  const imgs = el.querySelectorAll('img')
  await Promise.all([...imgs].map(img => {
    if (img.complete) return Promise.resolve()
    return new Promise(resolve => { img.onload = resolve; img.onerror = resolve })
  }))
  await new Promise(r => setTimeout(r, 300))
  const canvas = await html2canvas(el, { scale: 2, useCORS: true, allowTaint: true })
  const pdf = new jsPDF('p', 'mm', 'a4')
  const imgW = 170
  const imgH = (canvas.height * imgW) / canvas.width
  pdf.addImage(canvas.toDataURL('image/png'), 'PNG', 20, 20, imgW, imgH)
  pdf.save(`准考证_${receiptPermitNo.value}.pdf`)
}

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
.permit-section {
  max-width: 800px;
  margin: 3rem auto;
  padding: 0 1.5rem;
  text-align: center;
}
.permit-section-title { font-size: 1.6rem; color: #22c55e; margin: 0 0 .5rem; }
.permit-section-sub { font-size: 1.1rem; color: #334155; margin: 0 0 2rem; }
.permit-section-sub strong { font-size: 1.3rem; color: #1a73e8; }
.permit-wrapper { background: #fff; border: 2px solid #333; overflow: hidden; }
.permit-actions { display: flex; gap: 1rem; justify-content: center; margin-top: 2rem; }

/* Exam permit - Traditional style */
.exam-permit {
  padding: 2rem 2.5rem;
  font-family: "SimSun", "宋体", serif;
  color: #000;
  text-align: center;
}
.permit-header {
  border-bottom: 2px solid #000;
  padding-bottom: 1rem;
  margin-bottom: 1.5rem;
}
.permit-school-name {
  font-size: 1.4rem;
  font-weight: bold;
  letter-spacing: .3em;
  margin: 0 0 .5rem;
}
.permit-subtitle {
  font-size: .9rem;
  color: #333;
  margin: 0 0 .8rem;
}
.permit-main-title {
  font-size: 2.2rem;
  font-weight: bold;
  letter-spacing: 1em;
  margin: 0;
  color: #c41e3a;
}
.permit-body { text-align: left; }
.permit-info-section {
  display: flex;
  gap: 2rem;
  margin-bottom: 1.5rem;
}
.permit-photo-area {
  flex-shrink: 0;
  width: 120px;
  height: 160px;
  border: 1px solid #999;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  order: 2;
}
.permit-photo-area img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.photo-placeholder {
  color: #999;
  font-size: .9rem;
}
.permit-info-fields {
  flex: 1;
  margin: 0;
  order: 1;
}
.permit-info-fields div {
  display: flex;
  padding: .5rem 0;
  border-bottom: 1px dashed #ccc;
}
.permit-info-fields dt {
  white-space: nowrap;
  color: #333;
  font-weight: normal;
  min-width: 5.5em;
}
.permit-info-fields dd {
  margin: 0;
  font-weight: bold;
}
.permit-barcode {
  text-align: center;
  margin: 1rem 0;
}
.permit-barcode img {
  height: 60px;
  width: auto;
}
.permit-schedule {
  margin: 1.5rem 0;
}
.schedule-table {
  width: 100%;
  border-collapse: collapse;
  font-size: .85rem;
}
.schedule-table th,
.schedule-table td {
  border: 1px solid #333;
  padding: .6rem .8rem;
  text-align: center;
}
.schedule-table th {
  background: #f5f5f5;
  font-weight: bold;
}
.corner-cell {
  position: relative;
  width: 5rem;
}
.corner-cell .top-label {
  position: absolute;
  top: 4px;
  right: 8px;
  font-size: .75rem;
}
.corner-cell .bottom-label {
  position: absolute;
  bottom: 4px;
  left: 8px;
  font-size: .75rem;
}
.time-label {
  font-weight: bold;
  background: #f9f9f9;
}
.time-hint {
  font-size: .75rem;
  color: #666;
}
.permit-notice {
  margin-top: 1.2rem;
  padding-top: .8rem;
  border-top: 1px solid #ccc;
  font-size: .8rem;
  color: #c41e3a;
}

/* Photo upload */
.photo-upload-field > span {
  display: block;
  font-size: .85rem;
  font-weight: 600;
  color: #334155;
  margin-bottom: .4rem;
}
.photo-upload-area {
  width: 150px;
  height: 190px;
  border: 2px dashed #cbd5e1;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  transition: border-color .2s;
}
.photo-upload-area:hover { border-color: var(--signal); }
.photo-preview { width: 100%; height: 100%; object-fit: cover; }
.photo-upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: .4rem;
  color: #94a3b8;
  font-size: .8rem;
  text-align: center;
  padding: .5rem;
}
.upload-icon { font-size: 2rem; }
.upload-hint { font-size: .7rem; }
.field-error { color: #dc2626; font-size: .8rem; margin-top: .3rem; display: block; }
.field-hint { color: var(--signal); font-size: .8rem; margin-top: .3rem; display: block; }

/* Query */
.query-result { margin-top: 1rem; padding: 1rem; border-radius: 8px; font-size: .95rem; }
.query-result.ok { background: #f0fdf4; border: 1px solid #bbf7d0; }
.query-result.fail { background: #fef2f2; border: 1px solid #fecaca; color: #991b1b; }
.query-result p { margin: .3rem 0; }

@media (max-width: 600px) {
  .permit-info-section {
    flex-direction: column;
    align-items: center;
  }
  .permit-photo-area { order: 0; }
  .permit-info-fields { order: 1; }
  .exam-permit { padding: 1.5rem 1rem; }
}
</style>
