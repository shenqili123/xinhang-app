<template>
  <section class="page-hero" :style="{ '--page-hero-image': `url('/images/hero-commencement.jpg')` }">
    <div class="page-hero-inner reveal visible">
      <p class="school-label">{{ t('Admission', '招生') }}</p>
      <h1>{{ t('See the campus. Meet the people. Feel the culture.', '看见校园，遇见师生，感受文化。') }}</h1>
      <p>{{ t("The right school visit should show more than facilities. It should reveal how students live, how teachers guide, and whether the culture has the standards your family is looking for.", '一次真正有效的探校不只看设施，更要看学生如何生活、教师如何引导，以及这所学校的文化是否符合家庭期待的标准。') }}</p>
    </div>
  </section>

  <section class="page-section" id="visit">
    <div class="page-heading reveal visible">
      <p class="eyebrow">{{ t('Start Here', '从这里开始') }}</p>
      <h2>{{ t('A visit is the best first conversation.', '探校，是最好的第一次对话。') }}</h2>
      <p>{{ t("Families are invited to speak with the admission office, tour the campus, and discuss the student's academic background, boarding needs, and future goals.", '欢迎家庭与招生办公室沟通，参观校园，并讨论学生的学术背景、寄宿需求和未来目标。') }}</p>
    </div>
    <div class="admission-card-grid">
      <article class="feature-card reveal visible" v-for="(c, i) in cards" :key="i">
        <span>{{ String(i + 1).padStart(2, '0') }}</span>
        <h3>{{ t(c.en.title, c.zh.title) }}</h3>
        <p>{{ t(c.en.desc, c.zh.desc) }}</p>
      </article>
    </div>
    <div class="visit-actions reveal visible">
      <a class="btn btn-primary" href="tel:0531-00000000">{{ t('Call Admission', '联系招生办') }}</a>
      <router-link class="btn btn-secondary" to="/apply">{{ t('Apply Now', '立即报名') }}</router-link>
    </div>
  </section>

  <section class="image-split reverse">
    <figure class="visual-panel reveal visible">
      <img src="/images/commencement-class.jpg" alt="Xinhang graduating class" />
    </figure>
    <div class="text-stack reveal visible">
      <p class="eyebrow">{{ t('What Families Should Look For', '家庭应关注什么') }}</p>
      <h2>{{ t('Not just where students study, but who they become.', '不只是在哪里学习，而是会成为怎样的人。') }}</h2>
      <p>{{ t("Admission conversations should help families understand the school's expectations, daily rhythm, boarding culture, and how academic planning is connected to long-term growth.", '招生沟通应该帮助家庭理解学校的期待、日常节奏、寄宿文化，以及学业规划如何连接长期成长。') }}</p>
      <dl class="detail-list">
        <div><dt>{{ t('Academic Fit', '学术匹配') }}</dt><dd>{{ t('Does the student have the foundation and mindset to grow here?', '学生是否具备在这里成长所需的基础和心态？') }}</dd></div>
        <div><dt>{{ t('Residential Fit', '寄宿匹配') }}</dt><dd>{{ t('Will the daily structure help the student become more independent?', '日常结构是否能帮助学生变得更加独立？') }}</dd></div>
        <div><dt>{{ t('Future Fit', '未来匹配') }}</dt><dd>{{ t('Does the pathway match the family\'s ambitions and timeline?', '学校路径是否符合家庭的目标和时间规划？') }}</dd></div>
      </dl>
    </div>
  </section>

  <section class="steps" id="steps">
    <article class="reveal visible"><span>01</span><h3>{{ t('Inquire', '咨询') }}</h3><p>{{ t("Share the student's grade level, background, and areas of interest with the admission office.", '向招生办公室说明学生学段、背景和兴趣方向。') }}</p></article>
    <article class="reveal visible"><span>02</span><h3>{{ t('Visit', '探校') }}</h3><p>{{ t('Join an open day or schedule a family tour and consultation on campus.', '参加开放日，或预约家庭校园参观与咨询。') }}</p></article>
    <article class="reveal visible"><span>03</span><h3>{{ t('Apply', '申请') }}</h3><p>{{ t('Complete the admission conversation, required materials, and academic planning process.', '完成入学沟通、必要材料提交与学业规划流程。') }}</p></article>
  </section>
</template>

<script setup>
import { useLanguage } from '../composables/useLanguage'
const { t } = useLanguage()

const cards = [
  { en: { title: 'Family Consultation', desc: 'Talk through grade level, curriculum goals, language background, and possible pathway.' }, zh: { title: '家庭咨询', desc: '沟通学段、课程目标、语言基础与可能路径。' } },
  { en: { title: 'Campus Visit', desc: 'Walk the learning spaces, residential environment, activity venues, and signature campus areas.' }, zh: { title: '校园参观', desc: '走进学习空间、寄宿环境、活动场所和校园代表性区域。' } },
  { en: { title: 'Student Fit', desc: 'Review academic readiness, interests, habits, and the support a student may need to thrive.' }, zh: { title: '学生匹配', desc: '了解学生的学术准备、兴趣、习惯，以及成长所需的支持。' } },
]
</script>
