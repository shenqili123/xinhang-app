<template>
  <div ref="root">
    <section class="hero" aria-label="学校首屏介绍">
      <div class="hero-copy reveal">
        <h1><span>{{ t('Learn Deeply', '博学日新') }}</span><span>{{ t('Lead Wisely', '德行天下') }}</span></h1>
        <p class="hero-lede">{{ t('Academic rigor. Character formation. Global outlook.', '扎实学术 · 成熟人格 · 全球视野') }}</p>
        <div class="hero-actions">
          <router-link class="btn btn-primary" to="/apply">{{ t('Apply Now', '立即报名') }}</router-link>
          <router-link class="btn btn-secondary" to="/academics">{{ t('Explore Academics', '探索课程') }}</router-link>
        </div>
      </div>

      <figure class="hero-visual reveal">
        <img src="/assets/campus-overview.jpg" alt="Aerial view of Jinan Xinhang Experimental Foreign Language School campus" />
        <figcaption>
          <strong>{{ t('Campus at a Glance', '新航校园') }}</strong>
          <span>{{ t('A hillside campus built for serious learning and student growth.', '依山而建 · 面向世界') }}</span>
        </figcaption>
      </figure>

      <nav class="hero-portals" aria-label="Explore Xinhang">
        <router-link class="hero-portal reveal" to="/student-life">
          <img src="/assets/commencement-class.jpg" alt="Xinhang graduates celebrating together" />
          <span class="hero-portal-copy">
            <small>{{ t('Beyond the Classroom', '课堂之外') }}</small>
            <strong>{{ t('Student Life', '学生生活') }}</strong>
            <em>{{ t('Leadership, belonging and memorable traditions', '领导力、归属感与校园传统') }}</em>
          </span>
          <span class="hero-portal-arrow" aria-hidden="true">→</span>
        </router-link>
        <router-link class="hero-portal reveal" to="/about#faculty">
          <img src="/assets/hero-commencement.jpg" alt="A teacher accompanying Xinhang students at commencement" />
          <span class="hero-portal-copy">
            <small>{{ t('Known & Guided', '看见陪伴') }}</small>
            <strong>{{ t('Faculty & Mentorship', '教师团队') }}</strong>
            <em>{{ t('Teachers who challenge, support and stay alongside', '有要求、有支持、有长期陪伴') }}</em>
          </span>
          <span class="hero-portal-arrow" aria-hidden="true">→</span>
        </router-link>
        <router-link class="hero-portal reveal" to="/campus">
          <img src="/assets/campus-blossom.jpg" alt="Spring blossoms on the Xinhang campus" />
          <span class="hero-portal-copy">
            <small>{{ t('Life at Xinhang', '新航日常') }}</small>
            <strong>{{ t('Campus & Traditions', '校园文化') }}</strong>
            <em>{{ t('Spaces, ceremonies and the rhythm of school life', '校园空间、成长仪式与共同记忆') }}</em>
          </span>
          <span class="hero-portal-arrow" aria-hidden="true">→</span>
        </router-link>
      </nav>
    </section>

    <section class="pathfinder" id="explore" aria-label="Start here">
      <div class="pathfinder-copy reveal">
        <p class="eyebrow">{{ t('Start Here', '清晰入口') }}</p>
        <h2>{{ t('Find the right doorway into Xinhang.', '从这里进入新航。') }}</h2>
      </div>
      <div class="pathfinder-grid">
        <router-link class="pathfinder-card reveal" to="/apply"><span>01</span><strong>{{ t('Apply & Exam', '报名考试') }}</strong><em>{{ t('Online application, entrance permit, results and enrollment notice.', '在线填报、准考证生成、成绩查询与录取通知。') }}</em></router-link>
        <router-link class="pathfinder-card reveal" to="/academics"><span>02</span><strong>{{ t('Academic Pathways', '学段课程') }}</strong><em>{{ t('Primary, middle school, high school and international tracks.', '小学、初中、高中与国际方向的完整路径。') }}</em></router-link>
        <router-link class="pathfinder-card reveal" to="/student-life"><span>03</span><strong>{{ t('Student Life', '学生生活') }}</strong><em>{{ t('Leadership, clubs, boarding culture and school traditions.', '学生领导力、社团、寄宿文化与校园传统。') }}</em></router-link>
        <router-link class="pathfinder-card reveal" to="/campus"><span>04</span><strong>{{ t('Campus Tour', '校园探访') }}</strong><em>{{ t('Architecture, learning spaces, residence and athletics.', '建筑空间、学习场景、寄宿生活与运动设施。') }}</em></router-link>
      </div>
    </section>

    <section class="campus-feature" id="campus" aria-label="校园环境">
      <figure class="campus-feature-media reveal">
        <img src="/assets/campus-academic-building.jpg" alt="Xinhang academic building and landscaped campus path" />
      </figure>
      <div class="campus-feature-copy reveal">
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
      <div class="section-copy reveal">
        <p class="eyebrow">{{ t('About Xinhang', '关于新航') }}</p>
        <h2>{{ t('An international school with boarding-school polish.', '一所拥有寄宿学校气质的国际化学校。') }}</h2>
        <p>{{ t('Xinhang serves families who value serious academics, residential growth, and an international outlook. Classroom expectations, daily routines, college planning, and public expression are placed on one coherent path of development.', '新航面向期待扎实学术、寄宿成长与国际视野的家庭。学校把课堂要求、生活秩序、升学规划和公共表达放在同一条成长路径里，帮助学生形成更清楚的目标感。') }}</p>
        <p>{{ t('Education here is not only about results. It also values how students learn, work with others, and express themselves when the moment matters. A stable campus culture makes effort visible.', '这里的教育不是只追求结果，也重视学生如何学习、如何与人协作、如何在重要场合表达自己。稳定的校园文化，让每一次努力都能被看见。') }}</p>
      </div>
      <aside class="brand-panel reveal" aria-label="学校品牌细节">
        <img class="brand-panel-logo" src="/assets/school-name-gold.png" alt="山东新航实验国际学校" />
        <div class="brand-line"></div>
        <p>{{ t('Make Each Life Bloom Wonderfully', '让每一个生命精彩绽放') }}</p>
        <img class="brand-panel-motto" src="/assets/motto-gold-wide.png" alt="博学日新，德行天下" />
      </aside>
    </section>

    <section class="section academics-section" id="programs">
      <div class="section-heading reveal">
        <p class="eyebrow">{{ t('Academics', '学术课程') }}</p>
        <h2>{{ t('Built for students who are ready for more.', '为准备走得更远的学生而设。') }}</h2>
        <p>{{ t('From study habits to college pathways, the curriculum emphasizes strong fundamentals, English ability, academic expression, and future-ready choice.', '从学习习惯到升学路径，课程设计强调扎实基础、英文能力、学术表达和面向未来的选择力。') }}</p>
      </div>
      <div class="program-grid">
        <article class="program-card reveal"><span>01</span><h3>{{ t('Lower School', '小学部') }}</h3><p>{{ t('Lower School: reading, expression, inquiry, and self-management habits build a stable foundation for long-term learning.', '建立阅读、表达、探究和自我管理习惯，为长期学习打下稳定基础。') }}</p></article>
        <article class="program-card reveal"><span>02</span><h3>{{ t('Middle School', '初中部') }}</h3><p>{{ t('Middle School: national curriculum foundations, English ability, scientific thinking, and independent study methods.', '夯实国家课程，发展英语能力、科学思维和独立学习方法。') }}</p></article>
        <article class="program-card reveal"><span>03</span><h3>{{ t('Upper School', '高中部') }}</h3><p>{{ t('Upper School: academic planning, assessment readiness, and personalized support for multiple pathways.', '面向中高考和多元升学，提供规划、评估与个性化学业支持。') }}</p></article>
        <article class="program-card reveal accent"><span>04</span><h3>{{ t('Global Pathways', '国际方向') }}</h3><p>{{ t('International curriculum, cross-cultural confidence, and college guidance help students open wider choices.', '融合课程、跨文化能力和升学指导，帮助学生打开更宽的选择。') }}</p></article>
      </div>
    </section>

    <section class="statement" aria-label="学校教育理念">
      <div class="statement-inner reveal">
        <p class="eyebrow">{{ t('School Philosophy', '教育理念') }}</p>
        <h2>{{ t('Make Each Life Bloom Wonderfully.', '让每一个生命精彩绽放。') }}</h2>
        <p>{{ t("Every life should bloom wonderfully. The goal is not to shape students into one type, but to see each child's real possibility through academics, boarding, clubs, athletics, and college planning.", '不是把学生塑造成同一种样子，而是在学术、寄宿、社团、运动和升学规划里，看见每个孩子真正的可能。') }}</p>
      </div>
    </section>

    <section class="section life-section" id="life">
      <div class="section-heading reveal">
        <p class="eyebrow">{{ t('Student Life', '学生生活') }}</p>
        <h2>{{ t('Confidence grows outside the classroom, too.', '自信也在课堂之外生长。') }}</h2>
        <p>{{ t('Commencement, theater, athletics, and boarding life help students practice leadership, responsibility, and expression in real settings.', '毕业典礼、剧场活动、运动训练和寄宿生活，让学生在真实场景中练习领导力、责任感和表达能力。') }}</p>
      </div>
      <div class="life-grid">
        <article class="life-card reveal"><div class="life-img art"></div><div><h3>{{ t('Commencement & Leadership', '仪式与领导力') }}</h3><p>{{ t('Ceremony is not decoration. It teaches students responsibility, honor, and community.', '仪式感不是形式，而是让学生理解责任、荣誉和共同体的时刻。') }}</p></div></article>
        <article class="life-card reveal"><div class="life-img science"></div><div><h3>{{ t('Arts & Assembly', '艺术与集会') }}</h3><p>{{ t('Stage, music, visual work, and public expression help students become steady and confident communicators.', '舞台、音乐、视觉和公共表达训练，让学生形成稳定而自信的表达力。') }}</p></div></article>
        <article class="life-card reveal"><div class="life-img sport"></div><div><h3>{{ t('Boarding & Athletics', '寄宿与运动') }}</h3><p>{{ t('Rhythms, movement, teamwork, and self-management form a stable environment for growth.', '规律作息、运动习惯、同伴协作与生活自理，构成稳定的成长环境。') }}</p></div></article>
      </div>
    </section>

    <section class="admission-section" id="admission">
      <div class="admission-copy reveal">
        <p class="eyebrow">{{ t('Admission', '招生') }}</p>
        <h2>{{ t('See the campus. Meet the people. Feel the culture.', '看见校园，遇见师生，感受文化。') }}</h2>
        <p>{{ t("A strong school visit should reveal the school beyond classrooms: how students live, how teachers guide, and how the campus supports long-term growth.", '一次好的探校，应该让家庭看到课堂之外的学校气质：学生如何生活，教师如何陪伴，校园如何支持长期成长。') }}</p>
        <ol class="admission-route" aria-label="Admissions path">
          <li><span>01</span><strong>{{ t('Apply Online', '在线报名') }}</strong></li>
          <li><span>02</span><strong>{{ t('Receive Exam Notice', '接收考试通知') }}</strong></li>
          <li><span>03</span><strong>{{ t('Download Permit', '下载准考证') }}</strong></li>
          <li><span>04</span><strong>{{ t('Check Results', '查询成绩') }}</strong></li>
        </ol>
      </div>
      <div class="admission-actions reveal" id="visit">
        <a class="btn btn-primary" href="tel:0531-59562266">{{ t('Call Admission', '联系招生办') }}</a>
        <a class="btn btn-secondary" href="mailto:admissions@xinhang.school?subject=Schedule a Visit">{{ t('Schedule a Visit', '预约探校') }}</a>
        <router-link class="btn btn-light" to="/apply#application">{{ t('How to Apply', '申请流程') }}</router-link>
      </div>
    </section>

    <section class="steps" id="steps" aria-label="招生流程">
      <article class="reveal"><span>01</span><h3>{{ t('Inquire', '咨询') }}</h3><p>{{ t('Learn about divisions, boarding, curriculum, and enrollment plans.', '了解学段、寄宿、课程和入读计划。') }}</p></article>
      <article class="reveal"><span>02</span><h3>{{ t('Visit', '探校') }}</h3><p>{{ t('Join an open day or schedule a one-on-one campus visit.', '参加开放日或预约一对一校园参观。') }}</p></article>
      <article class="reveal"><span>03</span><h3>{{ t('Apply', '申请') }}</h3><p>{{ t('Complete the admission conversation, materials, and academic planning.', '完成入学沟通、材料提交与学业规划。') }}</p></article>
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
