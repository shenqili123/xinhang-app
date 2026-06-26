// Language toggle & navigation
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
