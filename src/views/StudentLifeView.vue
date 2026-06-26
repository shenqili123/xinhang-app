<template>
  <div ref="root">
    <section class="page-hero" style="--page-hero-image: url('/assets/commencement-class.jpg');" aria-label="Student Life">
      <div class="page-hero-inner reveal">
        <p class="school-label">{{ t('Student Life', '学生生活') }}</p>
        <h1>{{ t('Confidence grows outside the classroom, too.', '自信也在课堂之外生长。') }}</h1>
        <p>{{ t('Boarding routines, public ceremonies, arts, athletics, and student leadership give young people a place to practice maturity before the world asks for it.', '寄宿规律、公共典礼、艺术、运动和学生领导力，让年轻人在世界需要之前先练习成熟。') }}</p>
      </div>
    </section>

    <section class="page-section">
      <div class="page-heading reveal">
        <p class="eyebrow">{{ t('Daily Experience', '日常体验') }}</p>
        <h2>{{ t('A school culture students can feel.', '学生能感受到的校园文化。') }}</h2>
        <p>{{ t('Life at Xinhang is designed to make students visible: in the classroom, in the residence, on stage, on the field, and in the rituals that mark progress.', '在新航的生活让每个学生都能被看见：在课堂、在宿舍、在舞台、在操场、在标记成长的仪式中。') }}</p>
      </div>

      <div class="life-grid">
        <article class="life-card reveal">
          <div class="life-img art"></div>
          <div><h3>{{ t('Commencement & Leadership', '仪式与领导力') }}</h3><p>{{ t('Ceremony gives students a language for responsibility, honor, and belonging.', '仪式赋予学生关于责任、荣誉和归属的语言。') }}</p></div>
        </article>
        <article class="life-card reveal">
          <div class="life-img science"></div>
          <div><h3>{{ t('Arts & Assembly', '艺术与集会') }}</h3><p>{{ t('Students practice presence through performance, hosting, music, and public expression.', '学生通过表演、主持、音乐和公共表达练习如何有存在感。') }}</p></div>
        </article>
        <article class="life-card reveal">
          <div class="life-img sport"></div>
          <div><h3>{{ t('Boarding & Athletics', '寄宿与运动') }}</h3><p>{{ t('Routines, teams, and residential life build independence and steady habits.', '规律、团队和寄宿生活培养独立性和稳定的习惯。') }}</p></div>
        </article>
      </div>
    </section>

    <section class="image-split reverse" id="character">
      <figure class="visual-panel reveal">
        <img src="/assets/commencement-walk.jpg" alt="Graduates and faculty walking through campus" />
      </figure>
      <div class="text-stack reveal">
        <p class="eyebrow">{{ t('Traditions', '传统') }}</p>
        <h2>{{ t('Moments that students remember.', '学生会记住的时刻。') }}</h2>
        <p>{{ t('High-end private schools understand the power of ritual. At Xinhang, annual ceremonies and campus events help students see themselves as part of a larger story.', '好的学校懂得仪式的力量。在新航，年度典礼和校园活动让学生看到自己是更大故事的一部分。') }}</p>
        <dl class="detail-list">
          <div><dt>{{ t('Leadership', '领导力') }}</dt><dd>{{ t('Students take responsibility in assemblies, clubs, teams, and school events.', '学生在集会、社团、团队和校园活动中承担责任。') }}</dd></div>
          <div><dt>{{ t('Community', '社群') }}</dt><dd>{{ t('Shared milestones create a stronger sense of class identity and school pride.', '共同的里程碑增强年级归属和校园骄傲。') }}</dd></div>
          <div><dt>{{ t('Expression', '表达') }}</dt><dd>{{ t('Formal settings teach students how to speak, perform, and carry themselves.', '正式场合教会学生如何发言、如何表演、如何自持。') }}</dd></div>
        </dl>
      </div>
    </section>

    <section class="page-section student-signatures" id="ambassadors">
      <div class="page-heading reveal">
        <p class="eyebrow">{{ t('Student Voice', '学生之声') }}</p>
        <h2>{{ t('Leadership is practiced in public.', '领导力在公众面前练习。') }}</h2>
        <p>{{ t('From ambassadors to clubs, students learn to represent the school, organize peers, and turn interests into visible work.', '从学生大使到社团，学生学习代表学校、组织同伴、把兴趣变成可见的成果。') }}</p>
      </div>
      <div class="student-signature-grid">
        <article class="student-signature-card reveal"><span>01</span><h3>{{ t('Student Ambassadors', '学生大使') }}</h3><p>{{ t('Students welcome guests, guide campus visits, and learn the poise of speaking for a community.', '学生迎接来宾、带领参观，学习代表集体发言的从容。') }}</p></article>
        <article class="student-signature-card reveal"><span>02</span><h3>{{ t('Student Union', '学生会') }}</h3><p>{{ t('Student-led planning gives daily school life clearer rituals, events, and peer responsibility.', '学生主导的规划让日常校园生活有更清晰的仪式、活动和同伴责任。') }}</p></article>
        <article class="student-signature-card reveal"><span>03</span><h3>{{ t('Clubs & Societies', '社团') }}</h3><p>{{ t('Arts, sports, science, and service clubs help students build confidence beyond the classroom.', '艺术、体育、科学和公益社团帮助学生在课堂之外建立自信。') }}</p></article>
        <article class="student-signature-card reveal"><span>04</span><h3>{{ t('Challenge & Service', '挑战与服务') }}</h3><p>{{ t('Outdoor challenge, volunteer work, and project-based service give students a broader sense of purpose.', '户外挑战、志愿服务和项目式公益让学生拥有更广的目标感。') }}</p></article>
      </div>
    </section>

    <section class="page-section ivory">
      <div class="feature-grid">
        <article class="feature-card reveal"><span>01</span><h3>{{ t('Boarding', '寄宿') }}</h3><p>{{ t('Structured evenings, faculty support, and peer life help students practice independence with care.', '结构化晚间安排、教师支持和同伴生活帮助学生在关怀中练习独立。') }}</p></article>
        <article class="feature-card reveal" id="arts"><span>02</span><h3>{{ t('Clubs & Arts', '社团与艺术') }}</h3><p>{{ t('Creative and academic activities give students a wider way to discover interests and strengths.', '创意和学术活动为学生提供更宽的途径发现兴趣和特长。') }}</p></article>
        <article class="feature-card reveal" id="athletics"><span>03</span><h3>{{ t('Athletics', '运动') }}</h3><p>{{ t('Sports keep the day balanced and turn persistence, teamwork, and resilience into habit.', '运动让一天的节奏更均衡，把坚持、协作和韧性变成习惯。') }}</p></article>
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
