<template>
  <div ref="root">
    <section class="page-hero campus-hero" style="--page-hero-image: url('/assets/campus-overview.jpg');" aria-label="Campus">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('Campus', '校园') }}</p>
        <h1>{{ t('A place with room for ambition.', '一座让志向有空间生长的校园。') }}</h1>
        <p>{{ t('The campus gathers classrooms, performance spaces, athletics, residence, and ceremony into one daily environment. It feels open, ordered, and built for students who are growing fast.', '校园把教室、演出空间、运动场、宿舍和典礼聚合在一个日常环境里。它开阔、有序，为快速成长的学生而建。') }}</p>
      </div>
    </section>

    <section class="page-section">
      <div class="metric-strip reveal" aria-label="Campus highlights">
        <div><strong>K-12</strong><span>{{ t('Continuous pathway', '贯通培养') }}</span></div>
        <div><strong>{{ t('Boarding', '寄宿') }}</strong><span>{{ t('Structured student life', '有序的学生生活') }}</span></div>
        <div><strong>{{ t('Arts', '艺术') }}</strong><span>{{ t('Stage and assembly spaces', '舞台与集会空间') }}</span></div>
        <div><strong>{{ t('Athletics', '运动') }}</strong><span>{{ t('Fields, courts, and daily movement', '球场、跑道与日常运动') }}</span></div>
      </div>
    </section>

    <section class="image-split" id="environment">
      <figure class="visual-panel reveal">
        <img src="/assets/campus-blossom.jpg" alt="Xinhang campus building with blossoms" />
      </figure>
      <div class="text-stack reveal">
        <p class="eyebrow">{{ t('Campus Life', '校园日常') }}</p>
        <h2>{{ t('Elegant without feeling distant.', '典雅而不疏远。') }}</h2>
        <p>{{ t('The school grounds balance strong architecture with seasonal landscape. It gives the site a private-school presence while still feeling warm, walkable, and student-centered.', '校园在厚重建筑与四季景致之间取得平衡，既有寄宿学校的质感，又温暖、宜步行、以学生为中心。') }}</p>
        <dl class="detail-list">
          <div><dt>{{ t('Arrival', '初见') }}</dt><dd>{{ t('A clear first impression for visiting families and prospective students.', '为来访家庭和意向学生留下清晰的第一印象。') }}</dd></div>
          <div><dt>{{ t('Daily Rhythm', '日常节奏') }}</dt><dd>{{ t('Academic, athletic, and residential spaces are connected by a calm campus flow.', '教学、运动和生活空间以从容的动线相互连接。') }}</dd></div>
          <div><dt>{{ t('Occasion', '仪式场合') }}</dt><dd>{{ t('Formal ceremonies and public events have settings that feel memorable.', '正式典礼和公共活动有令人印象深刻的场地。') }}</dd></div>
        </dl>
      </div>
    </section>

    <section class="page-section ivory" id="culture">
      <div class="page-heading reveal">
        <p class="eyebrow">{{ t('Scenes of Xinhang', '新航场景') }}</p>
        <h2>{{ t('Campus as a living portfolio.', '校园是一本会生长的画册。') }}</h2>
      </div>
      <div class="media-grid">
        <figure class="media-card wide reveal">
          <img src="/assets/commencement-atrium.jpg" alt="Graduation atrium installation" />
          <figcaption>{{ t('Formal spaces for celebration and public memory.', '用于庆典和公共记忆的正式空间。') }}</figcaption>
        </figure>
        <figure class="media-card reveal">
          <img src="/assets/campus-academic-building.jpg" alt="Xinhang academic landscape" />
          <figcaption>{{ t('Academic buildings with mature landscape.', '学术建筑与成熟景观。') }}</figcaption>
        </figure>
        <figure class="media-card reveal">
          <img src="/assets/commencement-walk.jpg" alt="Graduates walking through campus" />
          <figcaption>{{ t('Campus pathways become part of school tradition.', '校园步道成为学校传统的一部分。') }}</figcaption>
        </figure>
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
