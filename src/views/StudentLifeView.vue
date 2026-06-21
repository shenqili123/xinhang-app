<template>
  <section class="page-hero" :style="{ '--page-hero-image': `url('/images/commencement-class.jpg')` }">
    <div class="page-hero-inner reveal visible">
      <p class="school-label">{{ t('Student Life', '学生生活') }}</p>
      <h1>{{ t('Confidence grows outside the classroom, too.', '自信也在课堂之外生长。') }}</h1>
      <p>{{ t('Boarding routines, public ceremonies, arts, athletics, and student leadership give young people a place to practice maturity before the world asks for it.', '寄宿节奏、公共仪式、艺术、运动与学生领导力，为年轻人在真正面对世界之前提供练习成熟的场域。') }}</p>
    </div>
  </section>

  <section class="page-section">
    <div class="page-heading reveal visible">
      <p class="eyebrow">{{ t('Daily Experience', '日常体验') }}</p>
      <h2>{{ t('A school culture students can feel.', '一种学生能感受到的学校文化。') }}</h2>
      <p>{{ t('Life at Xinhang is designed to make students visible: in the classroom, in the residence, on stage, on the field, and in the rituals that mark progress.', '新航的学生生活让学生被看见：在课堂、宿舍、舞台、运动场，以及标记成长的仪式里。') }}</p>
    </div>
    <div class="life-grid">
      <article class="life-card reveal visible" v-for="(c, i) in lifeCards" :key="i">
        <div :class="['life-img', c.cls]"></div>
        <div>
          <h3>{{ t(c.en.title, c.zh.title) }}</h3>
          <p>{{ t(c.en.desc, c.zh.desc) }}</p>
        </div>
      </article>
    </div>
  </section>

  <section class="image-split reverse">
    <figure class="visual-panel reveal visible">
      <img src="/images/commencement-walk.jpg" alt="Graduates walking" />
    </figure>
    <div class="text-stack reveal visible">
      <p class="eyebrow">{{ t('Traditions', '传统') }}</p>
      <h2>{{ t('Moments that students remember.', '让学生记住自己的时刻。') }}</h2>
      <p>{{ t('High-end private schools understand the power of ritual. At Xinhang, annual ceremonies and campus events help students see themselves as part of a larger story.', '高端私校懂得仪式的力量。在新航，年度典礼和校园活动帮助学生把自己放进一段更大的成长故事里。') }}</p>
      <dl class="detail-list">
        <div><dt>{{ t('Leadership', '领导力') }}</dt><dd>{{ t('Students take responsibility in assemblies, clubs, teams, and school events.', '学生在集会、社团、团队和学校活动中承担责任。') }}</dd></div>
        <div><dt>{{ t('Community', '共同体') }}</dt><dd>{{ t('Shared milestones create a stronger sense of class identity and school pride.', '共同里程碑带来更强的班级认同与学校荣誉感。') }}</dd></div>
        <div><dt>{{ t('Expression', '表达力') }}</dt><dd>{{ t('Formal settings teach students how to speak, perform, and carry themselves.', '正式场景训练学生如何发言、展示和呈现自己。') }}</dd></div>
      </dl>
    </div>
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

const lifeCards = [
  { cls: 'art', en: { title: 'Commencement & Leadership', desc: 'Ceremony gives students a language for responsibility, honor, and belonging.' }, zh: { title: '仪式与领导力', desc: '仪式让学生理解责任、荣誉与归属。' } },
  { cls: 'science', en: { title: 'Arts & Assembly', desc: 'Students practice presence through performance, hosting, music, and public expression.' }, zh: { title: '艺术与集会', desc: '学生通过表演、主持、音乐与公共表达练习稳定的在场感。' } },
  { cls: 'sport', en: { title: 'Boarding & Athletics', desc: 'Routines, teams, and residential life build independence and steady habits.' }, zh: { title: '寄宿与运动', desc: '作息、团队与寄宿生活共同建立独立性和稳定习惯。' } },
]

const features = [
  { en: { title: 'Boarding', desc: 'Structured evenings, faculty support, and peer life help students practice independence with care.' }, zh: { title: '寄宿', desc: '有结构的晚间生活、教师支持与同伴关系，帮助学生在被照看的环境中练习独立。' } },
  { en: { title: 'Clubs & Arts', desc: 'Creative and academic activities give students a wider way to discover interests and strengths.' }, zh: { title: '社团与艺术', desc: '创意与学术活动让学生以更多方式发现兴趣与优势。' } },
  { en: { title: 'Athletics', desc: 'Sports keep the day balanced and turn persistence, teamwork, and resilience into habit.' }, zh: { title: '运动', desc: '运动让一天保持平衡，也把坚持、协作与韧性变成习惯。' } },
]
</script>
