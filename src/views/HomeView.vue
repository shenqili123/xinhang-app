<template>
  <section class="hero" aria-label="学校首屏介绍">
    <div class="hero-copy reveal visible">
      <p class="school-label">{{ t('Shandong Xinhang Experimental International School', '山东新航实验国际学校') }}</p>
      <h1 v-html="t('<span>Think</span><span>Lead</span><span>Serve</span><span>Bloom</span>', '<span>博学</span><span>日新</span><span>德行</span><span>天下</span>')"></h1>
      <p class="hero-lede">{{ t('Challenge builds confidence. Mentorship, boarding culture, and global pathways prepare students to lead on campus and beyond.', '挑战塑造自信。导师陪伴、寄宿文化与全球路径，让学生在校园内外都能学会领导。') }}</p>
      <div class="hero-actions">
        <router-link class="btn btn-primary" to="/admission#visit">{{ t('Schedule a Visit', '预约探校') }}</router-link>
        <router-link class="btn btn-secondary" to="/academics">{{ t('Explore Academics', '探索课程') }}</router-link>
      </div>
    </div>
    <figure class="hero-visual reveal visible">
      <img src="/images/hero-commencement.jpg" alt="Xinhang commencement" />
      <figcaption>
        <strong>{{ t('Class of 2026', '2026届毕业生') }}</strong>
        <span>{{ t('Ceremony, mentorship, and aspiration in motion.', '仪式、陪伴与向上的愿望正在发生。') }}</span>
      </figcaption>
    </figure>
  </section>

  <section class="proof-bar" aria-label="学校亮点">
    <div><strong>K-12</strong><span>{{ t('continuous pathway', '一贯制成长路径') }}</span></div>
    <div><strong>Boarding</strong><span>{{ t('structured student life', '有秩序的寄宿生活') }}</span></div>
    <div><strong>Global</strong><span>{{ t('bilingual perspective', '双语与全球视野') }}</span></div>
    <div><strong>Jinan</strong><span>{{ t('rooted in Shandong', '根植济南与山东') }}</span></div>
  </section>

  <section class="campus-feature" id="campus" aria-label="校园环境">
    <figure class="campus-feature-media reveal visible">
      <img src="/images/campus-overview.jpg" alt="Xinhang campus aerial view" />
    </figure>
    <div class="campus-feature-copy reveal visible">
      <p class="eyebrow">{{ t('Campus', '校园') }}</p>
      <h2>{{ t('A place with room for ambition.', '一座让志向有空间生长的校园。') }}</h2>
      <p>{{ t('Built into the hillside, the campus brings academics, athletics, theater, and residential life into one daily rhythm. Open grounds and clear order give students a strong sense of direction, pace, and belonging.', '依山而建的校园，把学习、运动、剧场和生活空间放进同一条日常动线。开阔的场地与清晰的秩序，让学生每天都能感到目标、节奏和归属。') }}</p>
      <dl class="campus-facts">
        <div><dt>Academic</dt><dd>{{ t('Focused learning, projects, and college preparation', '专注学习、项目探究与升学准备') }}</dd></div>
        <div><dt>Residential</dt><dd>{{ t('Boarding routines and daily growth support', '寄宿生活与日常成长支持') }}</dd></div>
        <div><dt>Arts &amp; Athletics</dt><dd>{{ t('Theater, athletics, and school ceremonies', '剧场、运动场与校园仪式') }}</dd></div>
      </dl>
    </div>
  </section>

  <section class="section about-section" id="about">
    <div class="section-copy reveal visible">
      <p class="eyebrow">{{ t('About Xinhang', '关于新航') }}</p>
      <h2>{{ t('An international school with boarding-school polish.', '一所拥有寄宿学校气质的国际化学校。') }}</h2>
      <p>{{ t('Xinhang serves families who value serious academics, residential growth, and an international outlook. Classroom expectations, daily routines, college planning, and public expression are placed on one coherent path of development.', '新航面向期待扎实学术、寄宿成长与国际视野的家庭。学校把课堂要求、生活秩序、升学规划和公共表达放在同一条成长路径里，帮助学生形成更清楚的目标感。') }}</p>
      <p>{{ t("Education here is not only about results. It also values how students learn, work with others, and express themselves when the moment matters. A stable campus culture makes effort visible.", '这里的教育不是只追求结果，也重视学生如何学习、如何与人协作、如何在重要场合表达自己。稳定的校园文化，让每一次努力都能被看见。') }}</p>
    </div>
    <aside class="brand-panel reveal visible" aria-label="学校品牌">
      <img class="brand-panel-logo" src="/images/school-name-gold.png" alt="山东新航实验国际学校" />
      <div class="brand-line"></div>
      <p>{{ t('Make Each Life Bloom Wonderfully', '让每一个生命精彩绽放') }}</p>
      <img class="brand-panel-motto" src="/images/motto-gold-wide.png" alt="博学日新，德行天下" />
    </aside>
  </section>

  <section class="section academics-section" id="programs">
    <div class="section-heading reveal visible">
      <p class="eyebrow">{{ t('Academics', '学术课程') }}</p>
      <h2>{{ t('Built for students who are ready for more.', '为准备走得更远的学生而设。') }}</h2>
      <p>{{ t('From study habits to college pathways, the curriculum emphasizes strong fundamentals, English ability, academic expression, and future-ready choice.', '从学习习惯到升学路径，课程设计强调扎实基础、英文能力、学术表达和面向未来的选择力。') }}</p>
    </div>
    <div class="program-grid">
      <article class="program-card reveal visible" v-for="(p, i) in programs" :key="i" :class="{ accent: p.accent }">
        <span>{{ String(i + 1).padStart(2, '0') }}</span>
        <h3>{{ t(p.en.title, p.zh.title) }}</h3>
        <p>{{ t(p.en.desc, p.zh.desc) }}</p>
      </article>
    </div>
  </section>

  <section class="statement" aria-label="教育理念">
    <div class="statement-inner reveal visible">
      <p class="eyebrow">{{ t('School Philosophy', '教育理念') }}</p>
      <h2>{{ t('Make Each Life Bloom Wonderfully.', '让每一个生命精彩绽放。') }}</h2>
      <p>{{ t("Every life should bloom wonderfully. The goal is not to shape students into one type, but to see each child's real possibility through academics, boarding, clubs, athletics, and college planning.", '不是把学生塑造成同一种样子，而是在学术、寄宿、社团、运动和升学规划里，看见每个孩子真正的可能。') }}</p>
    </div>
  </section>

  <section class="section life-section" id="life">
    <div class="section-heading reveal visible">
      <p class="eyebrow">{{ t('Student Life', '学生生活') }}</p>
      <h2>{{ t('Confidence grows outside the classroom, too.', '自信也在课堂之外生长。') }}</h2>
      <p>{{ t('Commencement, theater, athletics, and boarding life help students practice leadership, responsibility, and expression in real settings.', '毕业典礼、剧场活动、运动训练和寄宿生活，让学生在真实场景中练习领导力、责任感和表达能力。') }}</p>
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

  <section class="admission-section" id="admission">
    <div class="admission-copy reveal visible">
      <p class="eyebrow">{{ t('Admission', '招生') }}</p>
      <h2>{{ t('See the campus. Meet the people. Feel the culture.', '看见校园，遇见师生，感受文化。') }}</h2>
      <p>{{ t('A strong school visit should reveal the school beyond classrooms: how students live, how teachers guide, and how the campus supports long-term growth.', '一次好的探校，应该让家庭看到课堂之外的学校气质：学生如何生活，教师如何陪伴，校园如何支持长期成长。') }}</p>
    </div>
    <div class="admission-actions reveal visible" id="visit">
      <a class="btn btn-primary" href="tel:0531-00000000">{{ t('Call Admission', '联系招生办') }}</a>
      <router-link class="btn btn-secondary" to="/apply">{{ t('Apply Now', '立即报名') }}</router-link>
      <router-link class="btn btn-light" to="/admission#steps">{{ t('How to Apply', '申请流程') }}</router-link>
    </div>
  </section>

  <section class="steps" id="steps" aria-label="招生流程">
    <article class="reveal visible"><span>01</span><h3>{{ t('Inquire', '咨询') }}</h3><p>{{ t('Learn about divisions, boarding, curriculum, and enrollment plans.', '了解学段、寄宿、课程和入读计划。') }}</p></article>
    <article class="reveal visible"><span>02</span><h3>{{ t('Visit', '探校') }}</h3><p>{{ t('Join an open day or schedule a one-on-one campus visit.', '参加开放日或预约一对一校园参观。') }}</p></article>
    <article class="reveal visible"><span>03</span><h3>{{ t('Apply', '申请') }}</h3><p>{{ t('Complete the admission conversation, materials, and academic planning.', '完成入学沟通、材料提交与学业规划。') }}</p></article>
  </section>
</template>

<script setup>
import { useLanguage } from '../composables/useLanguage'
const { t } = useLanguage()

const programs = [
  { en: { title: 'Lower School', desc: 'Lower School: reading, expression, inquiry, and self-management habits build a stable foundation for long-term learning.' }, zh: { title: '小学部', desc: '建立阅读、表达、探究和自我管理习惯，为长期学习打下稳定基础。' } },
  { en: { title: 'Middle School', desc: 'Middle School: national curriculum foundations, English ability, scientific thinking, and independent study methods.' }, zh: { title: '初中部', desc: '夯实国家课程，发展英语能力、科学思维和独立学习方法。' } },
  { en: { title: 'Upper School', desc: 'Upper School: academic planning, assessment readiness, and personalized support for multiple pathways.' }, zh: { title: '高中部', desc: '面向中高考和多元升学，提供规划、评估与个性化学业支持。' } },
  { accent: true, en: { title: 'Global Pathways', desc: 'International curriculum, cross-cultural confidence, and college guidance help students open wider choices.' }, zh: { title: '国际方向', desc: '融合课程、跨文化能力和升学指导，帮助学生打开更宽的选择。' } },
]

const lifeCards = [
  { cls: 'art', en: { title: 'Commencement & Leadership', desc: 'Ceremony is not decoration. It teaches students responsibility, honor, and community.' }, zh: { title: '仪式与领导力', desc: '仪式感不是形式，而是让学生理解责任、荣誉和共同体的时刻。' } },
  { cls: 'science', en: { title: 'Arts & Assembly', desc: 'Stage, music, visual work, and public expression help students become steady and confident communicators.' }, zh: { title: '艺术与集会', desc: '舞台、音乐、视觉和公共表达训练，让学生形成稳定而自信的表达力。' } },
  { cls: 'sport', en: { title: 'Boarding & Athletics', desc: 'Rhythms, movement, teamwork, and self-management form a stable environment for growth.' }, zh: { title: '寄宿与运动', desc: '规律作息、运动习惯、同伴协作与生活自理，构成稳定的成长环境。' } },
]
</script>
