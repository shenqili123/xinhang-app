<template>
  <section class="page-hero" :style="{ '--page-hero-image': `url('/images/campus-academic-building.jpg')` }">
    <div class="page-hero-inner reveal visible">
      <p class="school-label">{{ t('Academics', '学术课程') }}</p>
      <h1>{{ t('Built for students who are ready for more.', '为准备走得更远的学生而设。') }}</h1>
      <p>{{ t('Courses are shaped around deep literacy, confident English, disciplined inquiry, and the kind of academic habits students need for competitive secondary and university pathways.', '课程围绕深度阅读、自信英文、严谨探究和长期学术习惯展开，为学生面向更具竞争力的高中与大学路径做准备。') }}</p>
    </div>
  </section>

  <section class="page-section">
    <div class="page-heading reveal visible">
      <p class="eyebrow">{{ t('Academic Pathway', '学术路径') }}</p>
      <h2>{{ t('A coherent journey from foundation to future.', '从基础到未来的一条连贯路径。') }}</h2>
      <p>{{ t('Each division has its own developmental focus, but the through-line is consistent: strong fundamentals, thoughtful challenge, and visible growth.', '每个学段都有不同的发展重点，但主线一致：扎实基础、有分寸的挑战，以及看得见的成长。') }}</p>
    </div>
    <div class="program-grid">
      <article class="program-card reveal visible" v-for="(p, i) in programs" :key="i" :class="{ accent: p.accent }">
        <span>{{ String(i + 1).padStart(2, '0') }}</span>
        <h3>{{ t(p.en.title, p.zh.title) }}</h3>
        <p>{{ t(p.en.desc, p.zh.desc) }}</p>
      </article>
    </div>
  </section>

  <section class="image-split reverse">
    <figure class="visual-panel reveal visible">
      <img src="/images/commencement-atrium.jpg" alt="Graduation atrium" />
    </figure>
    <div class="text-stack reveal visible">
      <p class="eyebrow">{{ t('Learning Culture', '学习文化') }}</p>
      <h2>{{ t('Serious work, public confidence.', '认真学习，也自信表达。') }}</h2>
      <p>{{ t('Students are asked to do more than complete assignments. They learn to present, defend ideas, collaborate with peers, and carry themselves with maturity in formal school settings.', '学生不只是完成作业，也要学会展示、论证、协作，并在正式场合中以成熟姿态表达自己。') }}</p>
      <div class="mini-stat-grid">
        <div><strong>{{ t('Read', '阅读') }}</strong><span>{{ t('for precision', '追求准确') }}</span></div>
        <div><strong>{{ t('Write', '写作') }}</strong><span>{{ t('with structure', '有结构') }}</span></div>
        <div><strong>{{ t('Speak', '表达') }}</strong><span>{{ t('with confidence', '有自信') }}</span></div>
        <div><strong>{{ t('Lead', '领导') }}</strong><span>{{ t('with purpose', '有目标') }}</span></div>
      </div>
    </div>
  </section>

  <section class="pathway-band">
    <div class="page-heading reveal visible">
      <p class="eyebrow">{{ t('Preparation', '升学准备') }}</p>
      <h2>{{ t('College readiness starts before the application season.', '升学准备不从申请季才开始。') }}</h2>
    </div>
    <div class="feature-grid">
      <article class="feature-card dark reveal visible" v-for="(f, i) in prep" :key="i">
        <span>{{ String(i + 1).padStart(2, '0') }}</span>
        <h3>{{ t(f.en.title, f.zh.title) }}</h3>
        <p>{{ t(f.en.desc, f.zh.desc) }}</p>
      </article>
    </div>
  </section>
</template>

<script setup>
import { useLanguage } from '../composables/useLanguage'
const { t } = useLanguage()

const programs = [
  { en: { title: 'Lower School', desc: 'Reading, expression, curiosity, and healthy learning routines become the base for long-term academic confidence.' }, zh: { title: '小学部', desc: '阅读、表达、好奇心与健康学习习惯，构成长期学术自信的基础。' } },
  { en: { title: 'Middle School', desc: 'Students strengthen bilingual ability, scientific thinking, project habits, and independent study skills.' }, zh: { title: '初中部', desc: '学生强化双语能力、科学思维、项目习惯和独立学习能力。' } },
  { en: { title: 'Upper School', desc: 'Academic planning, subject depth, assessment readiness, and college counseling help students make purposeful choices.' }, zh: { title: '高中部', desc: '学业规划、学科深度、考试准备与升学指导，帮助学生做出更有目标感的选择。' } },
  { accent: true, en: { title: 'Global Pathways', desc: 'International curriculum, cross-cultural readiness, and application guidance open broader routes after graduation.' }, zh: { title: '国际方向', desc: '国际课程、跨文化准备和申请指导，为毕业后的更多路径打开可能。' } },
]

const prep = [
  { en: { title: 'Academic Advising', desc: 'Students receive guidance on course choices, study rhythms, and long-term goals.' }, zh: { title: '学业规划', desc: '学生在课程选择、学习节奏和长期目标上获得指导。' } },
  { en: { title: 'Portfolio Building', desc: 'Projects, activities, and leadership experiences become part of a coherent student story.' }, zh: { title: '成长档案', desc: '项目、活动与领导力经历共同构成清晰的学生成长叙事。' } },
  { en: { title: 'Application Guidance', desc: 'Families gain a clearer view of pathways, requirements, timelines, and fit.' }, zh: { title: '申请指导', desc: '家庭能够更清楚地理解路径、要求、时间线与匹配度。' } },
]
</script>
