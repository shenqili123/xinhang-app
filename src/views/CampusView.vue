<template>
  <div ref="root">
    <section class="page-hero campus-hero" :style="{ '--page-hero-image': 'url(/assets/campus-overview.jpg)' }">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('Campus', '校园') }}</p>
        <h1>{{ t('A place with room for ambition.', '一座让志向有空间生长的校园。') }}</h1>
        <p>{{ t('The campus gathers classrooms, performance spaces, athletics, residence, and ceremony into one daily environment. It feels open, ordered, and built for students who are growing fast.', '校园把教室、表演空间、运动场、寄宿生活与仪式场景放入同一套日常环境。它开阔、有秩序，也适合快速成长的学生。') }}</p>
      </div>
    </section>

    <section class="page-section">
      <div class="metric-strip reveal">
        <div><strong>K-12</strong><span>{{ t('Continuous pathway', '一贯制成长路径') }}</span></div>
        <div><strong>Boarding</strong><span>{{ t('Structured student life', '有秩序的学生生活') }}</span></div>
        <div><strong>{{ t('Arts', '艺术') }}</strong><span>{{ t('Stage and assembly spaces', '舞台与集会空间') }}</span></div>
        <div><strong>{{ t('Athletics', '运动') }}</strong><span>{{ t('Fields, courts, and daily movement', '场地、训练与日常活力') }}</span></div>
      </div>
    </section>

    <section class="image-split">
      <figure class="visual-panel reveal"><img src="/assets/campus-blossom.jpg" alt="Campus with blossoms" /></figure>
      <div class="text-stack reveal">
        <p class="eyebrow">{{ t('Campus Life', '校园生活') }}</p>
        <h2>{{ t('Elegant without feeling distant.', '高级，但不疏离。') }}</h2>
        <p>{{ t('The school grounds balance strong architecture with seasonal landscape. It gives the site a private-school presence while still feeling warm, walkable, and student-centered.', '校园用有力量的建筑与四季景观保持平衡，既有私校的体面感，也保留温度、可达性与学生中心的日常感。') }}</p>
        <dl class="detail-list">
          <div><dt>{{ t('Arrival', '抵达感') }}</dt><dd>{{ t('A clear first impression for visiting families and prospective students.', '给来访家庭和学生留下清晰而体面的第一印象。') }}</dd></div>
          <div><dt>{{ t('Daily Rhythm', '日常动线') }}</dt><dd>{{ t('Academic, athletic, and residential spaces are connected by a calm campus flow.', '学习、运动和生活空间由平静清晰的校园动线连接。') }}</dd></div>
          <div><dt>{{ t('Occasion', '仪式场景') }}</dt><dd>{{ t('Formal ceremonies and public events have settings that feel memorable.', '正式仪式与公共活动拥有足够令人记住的空间。') }}</dd></div>
        </dl>
      </div>
    </section>

    <section class="page-section ivory">
      <div class="page-heading reveal">
        <p class="eyebrow">{{ t('Scenes of Xinhang', '新航场景') }}</p>
        <h2>{{ t('Campus as a living portfolio.', '校园本身就是一份鲜活的作品集。') }}</h2>
      </div>
      <div class="media-grid">
        <figure class="media-card wide reveal"><img src="/assets/commencement-atrium.jpg" alt="Graduation atrium" /><figcaption>{{ t('Formal spaces for celebration and public memory.', '用于庆祝与共同记忆的正式空间。') }}</figcaption></figure>
        <figure class="media-card reveal"><img src="/assets/campus-academic-building.jpg" alt="Academic landscape" /><figcaption>{{ t('Academic buildings with mature landscape.', '学术建筑与成熟景观彼此呼应。') }}</figcaption></figure>
        <figure class="media-card reveal"><img src="/assets/commencement-walk.jpg" alt="Graduates walking" /><figcaption>{{ t('Campus pathways become part of school tradition.', '校园路径也成为学校传统的一部分。') }}</figcaption></figure>
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
