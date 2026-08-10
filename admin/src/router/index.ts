import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/dashboard/Dashboard.vue'),
    },
    {
      path: '/user',
      name: 'user',
      component: () => import('@/views/user/UserList.vue'),
    },
    {
      path: '/game',
      name: 'game',
      component: () => import('@/views/game/GameManage.vue'),
    },
    {
      path: '/log',
      name: 'log',
      component: () => import('@/views/log/LogView.vue'),
    },
  ],
})

export default router
