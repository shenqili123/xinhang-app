<template>
  <div ref="root" class="about-page">
    <section class="page-hero" style="--page-hero-image: url('/assets/campus-blossom.jpg');" aria-label="About Xinhang">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('About Xinhang', '关于新航') }}</p>
        <h1>{{ t('Purposeful, global, and quietly ambitious.', '有目标、有格局、有底气。') }}</h1>
        <p>{{ t('Xinhang is designed for families who want more than a school schedule: rigorous academics, steady mentorship, residential culture, and a pathway that prepares students to step confidently into the wider world.', '新航面向期待严谨学术、稳定陪伴、寄宿文化与成长路径的家庭，让学生在走向更宽广世界时更有信心。') }}</p>
      </div>
    </section>

    <section class="page-section" id="character">
      <div class="page-heading reveal">
        <p class="eyebrow">{{ t('School Character', '学校气质') }}</p>
        <h2>{{ t('An international school with boarding-school polish.', '一所拥有寄宿学校气质的国际化学校。') }}</h2>
      </div>
      <div class="lead-columns">
        <p class="reveal">{{ t('The school experience is intentionally structured: students learn in classrooms, practice independence in daily life, and build confidence through ceremony, public expression, athletics, and service.', '学校的体验是有意设计的：学生在课堂学习，在生活中练习独立，在仪式、表达、运动和服务中建立自信。') }}</p>
        <p class="reveal">{{ t('Its identity is rooted in Jinan and open to the world. The result should feel calm, selective, and serious: a campus where students are known, challenged, and expected to grow.', '它的特质根植济南、面向世界。校园给人的感受是沉稳、有要求的：每个学生都被看见、被挑战、被期待成长。') }}</p>
      </div>
    </section>

    <section class="image-split" id="faculty">
      <figure class="visual-panel reveal">
        <img src="/assets/campus-academic-building.jpg" alt="Xinhang academic building and campus landscape" />
      </figure>
      <div class="text-stack reveal">
        <p class="eyebrow">{{ t('Our Approach', '我们的方式') }}</p>
        <h2>{{ t('Growth is not left to chance.', '成长不靠偶然。') }}</h2>
        <p>{{ t('Academic planning, boarding routines, faculty guidance, and student activities work as one system. Each student is expected to develop habits of preparation, expression, collaboration, and reflection.', '学业规划、寄宿规律、教师指导和学生活动作为一套完整系统运行。每一位学生都被期待建立准备、表达、协作和反思的习惯。') }}</p>
        <dl class="detail-list">
          <div><dt>{{ t('Mentorship', '导师制') }}</dt><dd>{{ t('Faculty relationships that help students set goals and follow through.', '师生关系帮助学生设定目标并坚持到底。') }}</dd></div>
          <div><dt>{{ t('Belonging', '归属感') }}</dt><dd>{{ t('A school culture built through shared routines, houses, ceremonies, and teams.', '通过日常惯例、学院制、典礼和团队建设构建共同体文化。') }}</dd></div>
          <div><dt>{{ t('Standards', '高标准') }}</dt><dd>{{ t('Clear expectations for academic effort, conduct, and personal responsibility.', '对学术投入、行为规范和个人责任有清晰的要求。') }}</dd></div>
        </dl>
      </div>
    </section>

    <section class="signature-section">
      <div class="signature-copy reveal">
        <p class="eyebrow">{{ t('School Philosophy', '办学理念') }}</p>
        <h2>{{ t('Make Each Life Bloom Wonderfully.', '让每一个生命精彩绽放。') }}</h2>
        <p>{{ t('The school motto is treated as a practical promise: to notice each student clearly, to create room for excellence, and to help different forms of talent become visible.', '校训是一个实际的承诺：清晰地看见每一个学生，为优秀留出空间，让不同形式的才华被看见。') }}</p>
      </div>
      <img class="signature-mark reveal" src="/assets/motto-gold-wide.png" alt="School motto" />
    </section>

    <section class="page-section ivory" id="exchange">
      <div class="feature-grid">
        <article class="feature-card reveal">
          <span>01</span>
          <h3>{{ t('Academic Ambition', '学术追求') }}</h3>
          <p>{{ t('Students are expected to write, speak, reason, and build the discipline needed for demanding high school and university pathways.', '学生被期待具备写作、表达、推理能力，并形成面对高中和大学挑战所需的纪律性。') }}</p>
        </article>
        <article class="feature-card reveal">
          <span>02</span>
          <h3>{{ t('Residential Culture', '寄宿文化') }}</h3>
          <p>{{ t('Boarding life gives students a daily environment for independence, time management, and mature peer relationships.', '寄宿生活为学生提供练习独立、时间管理和成熟同伴关系的日常环境。') }}</p>
        </article>
        <article class="feature-card reveal">
          <span>03</span>
          <h3>{{ t('Global Outlook', '全球视野') }}</h3>
          <p>{{ t('English learning, cross-cultural confidence, and college planning are woven into the student\'s long-term development.', '英语学习、跨文化自信和升学规划贯穿学生的长期发展。') }}</p>
        </article>
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
