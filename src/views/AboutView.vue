<template>
  <section class="page-hero" :style="{ '--page-hero-image': `url('/images/campus-blossom.jpg')` }" aria-label="About">
    <div class="page-hero-inner reveal visible">
      <p class="school-label">{{ t('About Xinhang', '关于新航') }}</p>
      <h1>{{ t('Purposeful, global, and quietly ambitious.', '有目标，有世界感，也有安静的雄心。') }}</h1>
      <p>{{ t('Xinhang is designed for families who want more than a school schedule: rigorous academics, steady mentorship, residential culture, and a pathway that prepares students to step confidently into the wider world.', '新航面向期待更完整成长体验的家庭：扎实学术、稳定导师陪伴、寄宿文化，以及帮助学生自信走向更广阔世界的成长路径。') }}</p>
    </div>
  </section>

  <section class="page-section">
    <div class="page-heading reveal visible">
      <p class="eyebrow">{{ t('School Character', '学校气质') }}</p>
      <h2>{{ t('An international school with boarding-school polish.', '一所拥有寄宿学校气质的国际化学校。') }}</h2>
    </div>
    <div class="lead-columns">
      <p class="reveal visible">{{ t('The school experience is intentionally structured: students learn in classrooms, practice independence in daily life, and build confidence through ceremony, public expression, athletics, and service.', '学校体验被有意识地组织起来：学生在课堂中学习，在日常生活中练习独立，在仪式、表达、运动与服务中建立自信。') }}</p>
      <p class="reveal visible">{{ t('Its identity is rooted in Jinan and open to the world. The result should feel calm, selective, and serious: a campus where students are known, challenged, and expected to grow.', '它根植济南，也面向世界。我们想呈现的是一种克制、精选、认真而有分量的校园气质：学生被看见、被挑战，也被期待成长。') }}</p>
    </div>
  </section>

  <section class="image-split">
    <figure class="visual-panel reveal visible">
      <img src="/images/campus-academic-building.jpg" alt="Xinhang academic building" />
    </figure>
    <div class="text-stack reveal visible">
      <p class="eyebrow">{{ t('Our Approach', '我们的方式') }}</p>
      <h2>{{ t('Growth is not left to chance.', '成长不是交给偶然。') }}</h2>
      <p>{{ t('Academic planning, boarding routines, faculty guidance, and student activities work as one system. Each student is expected to develop habits of preparation, expression, collaboration, and reflection.', '学业规划、寄宿节奏、教师陪伴与学生活动彼此配合，形成一套完整系统。学生在其中建立准备、表达、协作与反思的习惯。') }}</p>
      <dl class="detail-list">
        <div><dt>{{ t('Mentorship', '导师陪伴') }}</dt><dd>{{ t('Faculty relationships that help students set goals and follow through.', '教师关系帮助学生设定目标，并持续推进。') }}</dd></div>
        <div><dt>{{ t('Belonging', '归属感') }}</dt><dd>{{ t('A school culture built through shared routines, houses, ceremonies, and teams.', '通过共同作息、社群、仪式与团队建立校园文化。') }}</dd></div>
        <div><dt>{{ t('Standards', '高标准') }}</dt><dd>{{ t('Clear expectations for academic effort, conduct, and personal responsibility.', '对学术努力、行为习惯和个人责任有清晰期待。') }}</dd></div>
      </dl>
    </div>
  </section>

  <section class="signature-section">
    <div class="signature-copy reveal visible">
      <p class="eyebrow">{{ t('School Philosophy', '教育理念') }}</p>
      <h2>{{ t('Make Each Life Bloom Wonderfully.', '让每一个生命精彩绽放。') }}</h2>
      <p>{{ t('The school motto is treated as a practical promise: to notice each student clearly, to create room for excellence, and to help different forms of talent become visible.', '这不是一句装饰性的口号，而是一种实践承诺：清楚看见每个学生，为优秀创造空间，让不同类型的天赋都能被看见。') }}</p>
    </div>
    <img class="signature-mark reveal visible" src="/images/motto-gold-wide.png" alt="School motto" />
  </section>

  <section class="page-section ivory">
    <div class="feature-grid">
      <article class="feature-card reveal visible" v-for="(f, i) in features" :key="i">
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

const features = [
  { en: { title: 'Academic Ambition', desc: 'Students are expected to write, speak, reason, and build the discipline needed for demanding high school and university pathways.' }, zh: { title: '学术志向', desc: '学生被期待具备写作、表达、推理能力，并建立面向高要求高中与大学路径的自律。' } },
  { en: { title: 'Residential Culture', desc: 'Boarding life gives students a daily environment for independence, time management, and mature peer relationships.' }, zh: { title: '寄宿文化', desc: '寄宿生活为学生提供练习独立、时间管理和成熟同伴关系的日常环境。' } },
  { en: { title: 'Global Outlook', desc: "English learning, cross-cultural confidence, and college planning are woven into the student's long-term development." }, zh: { title: '全球视野', desc: '英语学习、跨文化自信与升学规划被融入学生的长期发展。' } },
]
</script>
