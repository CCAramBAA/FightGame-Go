import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'connection',
      component: () => import('@/views/ConnectionPage.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginRegisterPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/lobby',
      name: 'lobby',
      component: () => import('@/views/LobbyPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/tutorial',
      name: 'tutorial',
      component: () => import('@/views/TutorialPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/pvp-rooms',
      name: 'pvp-rooms',
      component: () => import('@/views/PVPRoomsPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/pvp-select',
      name: 'pvp-select',
      component: () => import('@/views/PVPSelectPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/pve-stages',
      name: 'pve-stages',
      component: () => import('@/views/PVEStagePage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/pve-select',
      name: 'pve-select',
      component: () => import('@/views/PVESelectPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/game',
      name: 'game',
      component: () => import('@/views/GamePage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/replay',
      name: 'replay',
      component: () => import('@/views/ReplayPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/shop',
      name: 'shop',
      component: () => import('@/views/ShopPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/herodex',
      name: 'herodex',
      component: () => import('@/views/HeroDexPage.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')

  // 需要登录的页面
  if (to.meta.requiresAuth && !token) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  // 已登录则跳过登录页
  if (to.meta.guest && token) {
    next({ name: 'lobby' })
    return
  }

  next()
})

export default router
