import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Home', component: () => import('../views/HomeView.vue') },
  { path: '/about', name: 'About', component: () => import('../views/AboutView.vue') },
  { path: '/academics', name: 'Academics', component: () => import('../views/AcademicsView.vue') },
  { path: '/admission', name: 'Admission', component: () => import('../views/AdmissionView.vue') },
  { path: '/campus', name: 'Campus', component: () => import('../views/CampusView.vue') },
  { path: '/student-life', name: 'StudentLife', component: () => import('../views/StudentLifeView.vue') },
  { path: '/register', name: 'Register', component: () => import('../views/RegisterView.vue') },
  { path: '/login', name: 'Login', component: () => import('../views/LoginView.vue') },
  { path: '/apply', name: 'Apply', component: () => import('../views/ApplyView.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to) {
    if (to.hash) return { el: to.hash, behavior: 'smooth' }
    return { top: 0, behavior: 'smooth' }
  },
})

export default router
