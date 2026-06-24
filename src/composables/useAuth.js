import { ref, computed } from 'vue'

const token = ref(localStorage.getItem('token') || '')
const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

export function useAuth() {
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setAuth(t, u) {
    if (t !== undefined) {
      token.value = t
      localStorage.setItem('token', t)
    }
    if (u !== undefined) {
      user.value = u
      localStorage.setItem('user', JSON.stringify(u))
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function authHeader() {
    return token.value ? { Authorization: `Bearer ${token.value}` } : {}
  }

  return { token, user, isLoggedIn, isAdmin, setAuth, logout, authHeader }
}
