<template>
  <div ref="root">
    <section class="page-hero" :style="{ '--page-hero-image': 'url(/assets/hero-commencement.jpg)' }">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('Admission', '招生') }}</p>
        <h1>{{ t('See the campus. Meet the people. Feel the culture.', '看见校园，遇见师生，感受文化。') }}</h1>
        <p>{{ t("The right school visit should show more than facilities. It should reveal how students live, how teachers guide, and whether the culture has the standards your family is looking for.", '一次真正有效的探校不只看设施，更要看学生如何生活、教师如何引导，以及这所学校的文化是否符合家庭期待的标准。') }}</p>
      </div>
    </section>

    <section class="portal-banner" id="apply">
      <div class="reveal">
        <p class="eyebrow">{{ t('Admissions Portal', '招生入口') }}</p>
        <h2>{{ t('A clear path from first inquiry to exam results.', '从首次咨询到成绩查询，一条清晰路径。') }}</h2>
        <p>{{ t('Families can begin with a simple online application, receive exam information, download the entrance permit, and return later to check results.', '家庭可在线提交报名信息，接收考试安排，下载电子准考证，并在成绩发布后回到系统查询。') }}</p>
      </div>
      <ol class="portal-steps reveal">
        <li><span>01</span><strong>{{ t('Apply Online', '在线报名') }}</strong></li>
        <li><span>02</span><strong>{{ t('Exam Notice', '考试通知') }}</strong></li>
        <li><span>03</span><strong>{{ t('Entrance Permit', '准考证') }}</strong></li>
        <li><span>04</span><strong>{{ t('Results', '成绩查询') }}</strong></li>
      </ol>
      <router-link class="btn btn-primary reveal" to="/apply">{{ t('Start Here', '从这里开始') }}</router-link>
    </section>

    <section class="page-section" id="visit">
      <div class="page-heading reveal">
        <p class="eyebrow">{{ t('Start Here', '从这里开始') }}</p>
        <h2>{{ t('A visit is the best first conversation.', '探校，是最好的第一次对话。') }}</h2>
        <p>{{ t("Families are invited to speak with the admission office, tour the campus, and discuss the student's academic background, boarding needs, and future goals.", '欢迎家庭与招生办公室沟通，参观校园，并讨论学生的学术背景、寄宿需求和未来目标。') }}</p>
      </div>
      <div class="admission-card-grid">
        <article class="feature-card reveal"><span>01</span><h3>{{ t('Family Consultation', '家庭咨询') }}</h3><p>{{ t('Talk through grade level, curriculum goals, language background, and possible pathway.', '沟通学段、课程目标、语言基础与可能路径。') }}</p></article>
        <article class="feature-card reveal"><span>02</span><h3>{{ t('Campus Visit', '校园参观') }}</h3><p>{{ t('Walk the learning spaces, residential environment, activity venues, and signature campus areas.', '走进学习空间、寄宿环境、活动场所和校园代表性区域。') }}</p></article>
        <article class="feature-card reveal"><span>03</span><h3>{{ t('Student Fit', '学生匹配') }}</h3><p>{{ t('Review academic readiness, interests, habits, and the support a student may need to thrive.', '了解学生的学术准备、兴趣、习惯，以及成长所需的支持。') }}</p></article>
      </div>
      <div class="visit-actions reveal">
        <a class="btn btn-primary" href="tel:0531-00000000">{{ t('Call Admission', '联系招生办') }}</a>
        <a class="btn btn-secondary" href="mailto:admissions@xinhang.school?subject=Schedule a Visit">{{ t('Schedule a Visit', '预约探校') }}</a>
        <router-link class="btn btn-light" to="/apply">{{ t('Apply Online', '在线报名') }}</router-link>
      </div>
    </section>

    <section class="image-split reverse">
      <figure class="visual-panel reveal"><img src="/assets/commencement-class.jpg" alt="Graduating class" /></figure>
      <div class="text-stack reveal">
        <p class="eyebrow">{{ t('What Families Should Look For', '家庭应关注什么') }}</p>
        <h2>{{ t('Not just where students study, but who they become.', '不只是在哪里学习，而是会成为怎样的人。') }}</h2>
        <p>{{ t("Admission conversations should help families understand the school's expectations, daily rhythm, boarding culture, and how academic planning is connected to long-term growth.", '招生沟通应该帮助家庭理解学校的期待、日常节奏、寄宿文化，以及学业规划如何连接长期成长。') }}</p>
        <dl class="detail-list">
          <div><dt>{{ t('Academic Fit', '学术匹配') }}</dt><dd>{{ t('Does the student have the foundation and mindset to grow here?', '学生是否具备在这里成长所需的基础和心态？') }}</dd></div>
          <div><dt>{{ t('Residential Fit', '寄宿匹配') }}</dt><dd>{{ t('Will the daily structure help the student become more independent?', '日常结构是否能帮助学生变得更加独立？') }}</dd></div>
          <div><dt>{{ t('Future Fit', '未来匹配') }}</dt><dd>{{ t("Does the pathway match the family's ambitions and timeline?", '学校路径是否符合家庭的目标和时间规划？') }}</dd></div>
        </dl>
      </div>
    </section>

    <section class="steps" id="steps">
      <article class="reveal"><span>01</span><h3>{{ t('Apply', '报名') }}</h3><p>{{ t('Submit basic student and family information through the online admission entrance.', '通过在线招生入口提交学生与家庭基础信息。') }}</p></article>
      <article class="reveal"><span>02</span><h3>{{ t('Prepare', '备考') }}</h3><p>{{ t('Receive exam details, appointment information, and entrance permit instructions.', '接收考试安排、到校时间与准考证相关说明。') }}</p></article>
      <article class="reveal"><span>03</span><h3>{{ t('Results', '查询') }}</h3><p>{{ t('Return to the system after score release to check admission-related exam information.', '成绩发布后回到系统查询考试与录取相关信息。') }}</p></article>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useLanguage } from '../composables/useLanguage.js'
import { useReveal } from '../composables/useReveal.js'
const { t } = useLanguage()
const root = ref(null)
useReveal(root)
</script>
