<template>
  <footer class="footer">
    <div class="footer-brand">
      <img src="/assets/school-name-white.png" alt="山东新航实验国际学校" />
      <p>{{ t('Shandong Xinhang Experimental International School', '山东新航实验国际学校') }}</p>
    </div>
    <div class="footer-links">
      <router-link to="/admission">{{ t('Admission', '招生') }}</router-link>
      <router-link to="/about">{{ t('About', '关于新航') }}</router-link>
      <router-link to="/academics">{{ t('Academics', '学术课程') }}</router-link>
      <router-link to="/campus">{{ t('Campus', '校园') }}</router-link>
      <router-link to="/student-life">{{ t('Student Life', '学生生活') }}</router-link>
    </div>
  </footer>
</template>

<script setup>
import { useLanguage } from '../composables/useLanguage.js'
const { t } = useLanguage()
</script>
