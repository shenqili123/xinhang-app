<template>
  <footer class="footer">
    <div class="footer-brand">
      <img src="/assets/school-name-white.png" alt="济南新航实验外国语学校" />
    </div>
    <div class="footer-links">
      <router-link to="/apply">{{ t('Apply Now', '立即报名') }}</router-link>
      <router-link to="/about">{{ t('About', '关于新航') }}</router-link>
      <router-link to="/academics">{{ t('Academics', '学术课程') }}</router-link>
      <router-link to="/campus">{{ t('Campus', '校园') }}</router-link>
      <router-link to="/student-life">{{ t('Student Life', '学生生活') }}</router-link>
      <router-link to="/news">{{ t('News', '新闻动态') }}</router-link>
      <router-link to="/login">{{ t('Sign In', '登录') }}</router-link>
    </div>
  </footer>
</template>

<script setup>
import { useLanguage } from '../composables/useLanguage.js'
const { t } = useLanguage()
</script>
