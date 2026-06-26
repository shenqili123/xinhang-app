<template>
  <div ref="root">
    <section class="page-hero" style="--page-hero-image: url('/assets/campus-academic-building.jpg');" aria-label="Academics">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('Academics', '学术课程') }}</p>
        <h1>{{ t('Built for students who are ready for more.', '为准备走得更远的学生而设计。') }}</h1>
        <p>{{ t('Courses are shaped around deep literacy, confident English, disciplined inquiry, and the kind of academic habits students need for competitive secondary and university pathways.', '课程围绕深度读写、自信英语、学术探究和面向未来升学的学习习惯来设计。') }}</p>
      </div>
    </section>

    <section class="page-section" id="pathways">
      <div class="page-heading reveal">
        <p class="eyebrow">{{ t('Academic Pathway', '学段路径') }}</p>
        <h2>{{ t('A coherent journey from foundation to future.', '从基础到未来的连贯成长路径。') }}</h2>
        <p>{{ t('Each division has its own developmental focus, but the through-line is consistent: strong fundamentals, thoughtful challenge, and visible growth.', '每个学段有不同的成长重点，但主线始终一致：扎实基础、有思考的挑战、可见的进步。') }}</p>
      </div>

      <div class="program-grid">
        <article class="program-card reveal"><span>01</span><h3>{{ t('Lower School', '小学部') }}</h3><p>{{ t('Reading, expression, curiosity, and healthy learning routines become the base for long-term academic confidence.', '阅读、表达、好奇心和健康的学习习惯是长期学术自信的基础。') }}</p></article>
        <article class="program-card reveal"><span>02</span><h3>{{ t('Middle School', '初中部') }}</h3><p>{{ t('Students strengthen bilingual ability, scientific thinking, project habits, and independent study skills.', '学生加强双语能力、科学思维、项目学习和自主学习能力。') }}</p></article>
        <article class="program-card reveal"><span>03</span><h3>{{ t('Upper School', '高中部') }}</h3><p>{{ t('Academic planning, subject depth, assessment readiness, and college counseling help students make purposeful choices.', '学术规划、学科深度、评估准备和升学指导帮助学生做出有方向的选择。') }}</p></article>
        <article class="program-card reveal accent"><span>04</span><h3>{{ t('Global Pathways', '国际方向') }}</h3><p>{{ t('International curriculum, cross-cultural readiness, and application guidance open broader routes after graduation.', '国际课程、跨文化能力和申请指导，为毕业后打开更宽的路径。') }}</p></article>
      </div>
    </section>

    <section class="image-split reverse" id="curriculum">
      <figure class="visual-panel reveal">
        <img src="/assets/commencement-atrium.jpg" alt="Graduation displays in the Xinhang atrium" />
      </figure>
      <div class="text-stack reveal">
        <p class="eyebrow">{{ t('Learning Culture', '学习文化') }}</p>
        <h2>{{ t('Serious work, public confidence.', '严肃学习，从容表达。') }}</h2>
        <p>{{ t('Students are asked to do more than complete assignments. They learn to present, defend ideas, collaborate with peers, and carry themselves with maturity in formal school settings.', '学生不仅仅完成作业，更要学习展示、为观点辩护、与同伴协作，在正式场合展现成熟。') }}</p>
        <div class="mini-stat-grid">
          <div><strong>{{ t('Read', '阅读') }}</strong><span>{{ t('for precision', '精准理解') }}</span></div>
          <div><strong>{{ t('Write', '写作') }}</strong><span>{{ t('with structure', '结构清晰') }}</span></div>
          <div><strong>{{ t('Speak', '表达') }}</strong><span>{{ t('with confidence', '从容自信') }}</span></div>
          <div><strong>{{ t('Lead', '领导') }}</strong><span>{{ t('with purpose', '目标明确') }}</span></div>
        </div>
      </div>
    </section>

    <section class="pathway-band" id="global">
      <div class="page-heading reveal">
        <p class="eyebrow">{{ t('Preparation', '升学准备') }}</p>
        <h2>{{ t('College readiness starts before the application season.', '升学准备从申请季之前就开始了。') }}</h2>
      </div>
      <div class="feature-grid">
        <article class="feature-card dark reveal"><span>01</span><h3>{{ t('Academic Advising', '学术规划') }}</h3><p>{{ t('Students receive guidance on course choices, study rhythms, and long-term goals.', '学生获得选课、学习节奏和长期目标方面的指导。') }}</p></article>
        <article class="feature-card dark reveal"><span>02</span><h3>{{ t('Portfolio Building', '成长档案') }}</h3><p>{{ t('Projects, activities, and leadership experiences become part of a coherent student story.', '项目、活动和领导力经历构成连贯的学生成长叙事。') }}</p></article>
        <article class="feature-card dark reveal"><span>03</span><h3>{{ t('Application Guidance', '申请指导') }}</h3><p>{{ t('Families gain a clearer view of pathways, requirements, timelines, and fit.', '家庭对升学路径、要求、时间表和适配度有更清晰的了解。') }}</p></article>
      </div>
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
