/**
 * Vue → Static HTML Export Script
 * 将当前 Vue 前端导出为静态 HTML，供同事在纯 HTML 环境中编辑
 * 用法: node export-html.js
 * 输出: ./html-export/ 文件夹
 */
const fs = require('fs')
const path = require('path')

const SRC = path.join(__dirname, 'src')
const OUT = path.join(__dirname, 'html-export')

// Pages to export (content pages only, skip functional pages like Login/Register)
const pages = [
  { vue: 'HomeView.vue', html: 'index.html', title: '首页' },
  { vue: 'AboutView.vue', html: 'about.html', title: '关于新航' },
  { vue: 'AcademicsView.vue', html: 'academics.html', title: '学术课程' },
  { vue: 'CampusView.vue', html: 'campus.html', title: '校园' },
  { vue: 'StudentLifeView.vue', html: 'student-life.html', title: '学生生活' },
  { vue: 'ApplyView.vue', html: 'apply.html', title: '招生报名' },
]

const routeToFile = {
  '/': 'index.html',
  '/about': 'about.html',
  '/academics': 'academics.html',
  '/campus': 'campus.html',
  '/student-life': 'student-life.html',
  '/apply': 'apply.html',
  '/news': 'news.html',
  '/login': 'login.html',
  '/register': 'register.html',
  '/profile': 'profile.html',
}

function extractTemplate(vueContent) {
  const match = vueContent.match(/<template>([\s\S]*?)<\/template>/)
  return match ? match[1].trim() : ''
}

// Replace {{ t('English text', '中文文本') }} → English text / 中文文本
function resolveTranslations(html) {
  // t('en', 'zh') → "en" with Chinese as title attribute
  return html.replace(/\{\{\s*t\(\s*'([^']*)'\s*,\s*'([^']*)'\s*\)\s*\}\}/g, (_, en, zh) => {
    return `<span data-en="${en}" data-zh="${zh}">${en}</span>`
  })
}

// Convert <router-link to="/path"> → <a href="file.html">
function convertRouterLinks(html) {
  // <router-link ... to="/path#anchor" ...>text</router-link>
  return html.replace(/<router-link\b[^>]*\bto="([^"]*)"[^>]*>([\s\S]*?)<\/router-link>/g, (_, to, content) => {
    let href = to
    const hashIdx = to.indexOf('#')
    const route = hashIdx >= 0 ? to.substring(0, hashIdx) : to
    const hash = hashIdx >= 0 ? to.substring(hashIdx) : ''

    if (routeToFile[route]) {
      href = routeToFile[route] + hash
    } else {
      href = route.replace(/^\//, '') + '.html' + hash
    }
    // Preserve class if it has btn
    const classMatch = _.match(/class="([^"]*)"/)
    const cls = classMatch ? ` class="${classMatch[1].replace(/[^a-zA-Z0-9\s-]/g, '')}"` : ''
    return `<a href="${href}"${cls}>${content}</a>`
  })
}

// Remove Vue directives and dynamic bindings
function cleanVueDirectives(html) {
  // Remove v-if, v-else, v-else-if, v-for, v-show, v-model, v-bind, @events
  html = html.replace(/\s+v-if="[^"]*"/g, '')
  html = html.replace(/\s+v-else-if="[^"]*"/g, '')
  html = html.replace(/\s+v-else/g, '')
  html = html.replace(/\s+v-for="[^"]*"/g, '')
  html = html.replace(/\s+v-show="[^"]*"/g, '')
  html = html.replace(/\s+v-model="[^"]*"/g, '')
  html = html.replace(/\s+v-html="[^"]*"/g, '')
  html = html.replace(/\s+@[\w.-]+="[^"]*"/g, '')
  html = html.replace(/\s+@[\w.-]+/g, '')
  // Remove :class, :style, :aria-*, :disabled, :src bindings
  html = html.replace(/\s+:class="[^"]*"/g, '')
  html = html.replace(/\s+:style="[^"]*"/g, '')
  html = html.replace(/\s+:aria-[\w-]+="[^"]*"/g, '')
  html = html.replace(/\s+:disabled="[^"]*"/g, '')
  html = html.replace(/\s+:src="[^"]*"/g, '')
  html = html.replace(/\s+:alt="[^"]*"/g, '')
  html = html.replace(/\s+:placeholder="[^"]*"/g, '')
  html = html.replace(/\s+active-class="[^"]*"/g, '')
  html = html.replace(/\s+exact-active-class="[^"]*"/g, '')
  // Remove <template> wrapper tags (Vue fragments)
  html = html.replace(/<template>/g, '')
  html = html.replace(/<\/template>/g, '')
  // Remove remaining {{ ... }} expressions
  html = html.replace(/\{\{[^}]*\}\}/g, '')
  // Remove ref="..."
  html = html.replace(/\s+ref="[^"]*"/g, '')
  return html
}

function buildHeader() {
  return `  <header class="site-header">
    <a class="brand-block" href="index.html" aria-label="山东新航实验国际学校首页">
      <img src="assets/school-name-white.png" alt="山东新航实验国际学校" />
    </a>

    <div class="utility-bar">
      <nav class="utility-nav" aria-label="辅助导航">
        <a href="news.html"><span data-en="News" data-zh="新闻动态">News</span></a>
        <a href="login.html"><span data-en="Sign In" data-zh="登录">Sign In</span></a>
        <a href="register.html"><span data-en="Register" data-zh="注册">Register</span></a>
      </nav>
    </div>

    <div class="primary-bar">
      <button class="lang-toggle" type="button" aria-label="切换语言">中文</button>

      <button class="nav-toggle" aria-label="打开导航" aria-expanded="false">☰</button>
      <nav class="main-nav" aria-label="主导航">
        <div class="nav-item">
          <a href="about.html">About / 关于新航</a>
          <div class="nav-submenu">
            <a href="about.html">School Overview / 学校概况</a>
            <a href="about.html#character">Introduction / 学校简介</a>
            <a href="campus.html">Campus Walk / 漫步校园</a>
            <a href="about.html#exchange">International Exchange / 对外交流</a>
          </div>
        </div>
        <div class="nav-item">
          <a href="academics.html">Academics / 学术课程</a>
          <div class="nav-submenu">
            <a href="academics.html">Teaching & Research / 教学教研</a>
            <a href="academics.html#curriculum">School-based Curriculum / 校本课程</a>
            <a href="academics.html#global">International Academics / 国际部学术</a>
            <a href="academics.html#pathways">Curriculum System / 课程体系</a>
          </div>
        </div>
        <div class="nav-item">
          <a href="campus.html">Campus / 校园</a>
          <div class="nav-submenu">
            <a href="news.html">Campus News / 校园动态</a>
            <a href="campus.html#culture">Civilized Campus / 文明校园风采</a>
            <a href="campus.html">Library / 图书馆</a>
            <a href="campus.html#environment">Campus Environment / 校园环境</a>
          </div>
        </div>
        <div class="nav-item">
          <a href="student-life.html">Student Life / 学生生活</a>
          <div class="nav-submenu">
            <a href="student-life.html">Student Activities / 学生活动</a>
            <a href="student-life.html#character">Character Education / 德育天地</a>
            <a href="student-life.html#ambassadors">Student Ambassadors / 学生大使</a>
            <a href="student-life.html#arts">Arts / 艺术</a>
            <a href="student-life.html#athletics">Athletics / 运动</a>
          </div>
        </div>
        <a class="nav-apply" href="apply.html">Apply Now / 立即报名</a>
      </nav>
    </div>
  </header>`
}

function buildFooter() {
  return `  <footer class="footer">
    <div class="footer-brand">
      <img src="assets/school-name-white.png" alt="山东新航实验国际学校" />
      <p>Shandong Xinhang Experimental International School / 山东新航实验国际学校</p>
    </div>
    <div class="footer-links">
      <a href="apply.html">Apply Now / 立即报名</a>
      <a href="about.html">About / 关于新航</a>
      <a href="academics.html">Academics / 学术课程</a>
      <a href="campus.html">Campus / 校园</a>
      <a href="student-life.html">Student Life / 学生生活</a>
      <a href="news.html">News / 新闻动态</a>
      <a href="login.html">Sign In / 登录</a>
    </div>
  </footer>`
}

function buildPage(title, bodyContent) {
  return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>${title} - 山东新航实验国际学校</title>
  <link rel="stylesheet" href="styles.css" />
  <script src="script.js" defer><\/script>
</head>
<body>
${buildHeader()}

${bodyContent}

${buildFooter()}
</body>
</html>
`
}

// --- Main ---
if (fs.existsSync(OUT)) fs.rmSync(OUT, { recursive: true })
fs.mkdirSync(OUT, { recursive: true })

// Export pages
for (const page of pages) {
  const vuePath = path.join(SRC, 'views', page.vue)
  if (!fs.existsSync(vuePath)) {
    console.log(`⚠ Skipped ${page.vue} (not found)`)
    continue
  }
  const raw = fs.readFileSync(vuePath, 'utf-8')
  let template = extractTemplate(raw)
  template = resolveTranslations(template)
  template = convertRouterLinks(template)
  template = cleanVueDirectives(template)
  // Remove the outer <div ref="root"> wrapper
  template = template.replace(/^\s*<div>\s*/i, '').replace(/\s*<\/div>\s*$/i, '')

  const html = buildPage(page.title, template)
  fs.writeFileSync(path.join(OUT, page.html), html, 'utf-8')
  console.log(`✓ ${page.html}`)
}

// Copy CSS
const cssPath = path.join(SRC, 'assets', 'styles.css')
if (fs.existsSync(cssPath)) {
  let css = fs.readFileSync(cssPath, 'utf-8')
  // Convert /assets/ paths to assets/ (relative)
  css = css.replace(/url\(["']?\/assets\//g, 'url("assets/')
  css = css.replace(/url\(["']?assets\//g, 'url("assets/')
  fs.writeFileSync(path.join(OUT, 'styles.css'), css, 'utf-8')
  console.log('✓ styles.css')
}

// Generate script.js with language toggle and navigation
const scriptContent = `// Language toggle & navigation
document.addEventListener('DOMContentLoaded', () => {
  let lang = localStorage.getItem('lang') || 'en';

  function applyLang() {
    document.querySelectorAll('[data-en]').forEach(el => {
      el.textContent = lang === 'zh' ? el.dataset.zh : el.dataset.en;
    });
    const btn = document.querySelector('.lang-toggle');
    if (btn) btn.textContent = lang === 'zh' ? 'EN' : '中文';
  }

  const langBtn = document.querySelector('.lang-toggle');
  if (langBtn) {
    langBtn.addEventListener('click', () => {
      lang = lang === 'zh' ? 'en' : 'zh';
      localStorage.setItem('lang', lang);
      applyLang();
    });
  }

  applyLang();

  // Mobile menu toggle
  const toggle = document.querySelector('.nav-toggle');
  const nav = document.querySelector('.main-nav');
  if (toggle && nav) {
    toggle.addEventListener('click', () => {
      const open = nav.classList.toggle('open');
      toggle.setAttribute('aria-expanded', open);
    });
  }

  // Scroll header
  const header = document.querySelector('.site-header');
  if (header) {
    const onScroll = () => header.classList.toggle('scrolled', window.scrollY > 24);
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
  }

  // Reveal animation
  const reveals = document.querySelectorAll('.reveal');
  if (reveals.length) {
    const io = new IntersectionObserver(entries => {
      entries.forEach(e => { if (e.isIntersecting) { e.target.classList.add('visible'); io.unobserve(e.target); } });
    }, { threshold: 0.15 });
    reveals.forEach(el => io.observe(el));
  }
});
`
fs.writeFileSync(path.join(OUT, 'script.js'), scriptContent, 'utf-8')
console.log('✓ script.js')

// Copy assets folder
const assetsDir = path.join(__dirname, 'public', 'assets')
const outAssets = path.join(OUT, 'assets')
if (fs.existsSync(assetsDir)) {
  fs.cpSync(assetsDir, outAssets, { recursive: true })
  console.log('✓ assets/ (images)')
}

// Create README
const readme = `# 新航前端 HTML 版本（自动导出）

## 说明
这是从 Vue 前端自动导出的静态 HTML 版本，用于设计修改。

## 文件结构
- index.html — 首页
- about.html — 关于新航
- academics.html — 学术课程
- campus.html — 校园
- student-life.html — 学生生活
- apply.html — 招生报名
- styles.css — 全部样式
- script.js — 语言切换 + 导航 + 动画
- assets/ — 图片资源

## 使用方式
1. 直接用浏览器打开任意 .html 文件即可预览
2. 修改 HTML 内容后保存，刷新浏览器查看效果
3. 点击右上角"中文"按钮可切换中英文

## 双语文本格式
所有双语文本都使用这种结构:
\`\`\`html
<span data-en="English text" data-zh="中文文本">English text</span>
\`\`\`
- 默认显示 data-en 的内容
- 切换语言后显示 data-zh 的内容
- 修改时请同时更新 data-en 和 data-zh 属性

## 注意事项
- 这只是展示层，不包含后端逻辑（登录、报名提交等）
- 报名表仅保留 HTML 结构，表单提交功能需要后端支持
- 修改完成后请将改动后的 HTML 文件发回，我们会同步到 Vue 版本

## 导出时间
${new Date().toLocaleString('zh-CN')}
`
fs.writeFileSync(path.join(OUT, 'README.md'), readme, 'utf-8')
console.log('✓ README.md')

console.log(`\n✅ 导出完成! 输出目录: ${OUT}`)
