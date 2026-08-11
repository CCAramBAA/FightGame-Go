import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/dashboard/Dashboard.vue'),
    meta: { title: '仪表盘' },
  },
  {
    path: '/users',
    name: 'UserList',
    component: () => import('@/views/user/UserList.vue'),
    meta: { title: '用户管理' },
  },
  {
    path: '/characters',
    name: 'CharacterManage',
    component: () => import('@/views/character/CharacterManage.vue'),
    meta: { title: '角色管理' },
  },
  {
    path: '/skills',
    name: 'SkillManage',
    component: () => import('@/views/skill/SkillManage.vue'),
    meta: { title: '技能管理' },
  },
  {
    path: '/skins',
    name: 'SkinManage',
    component: () => import('@/views/skin/SkinManage.vue'),
    meta: { title: '皮肤管理' },
  },
  {
    path: '/special-rules',
    name: 'SpecialRuleManage',
    component: () => import('@/views/rule/SpecialRuleManage.vue'),
    meta: { title: '英雄特殊规则' },
  },
  {
    path: '/stages',
    name: 'StageManage',
    component: () => import('@/views/stage/StageManage.vue'),
    meta: { title: '关卡管理' },
  },
  {
    path: '/shop',
    name: 'ShopManage',
    component: () => import('@/views/shop/ShopManage.vue'),
    meta: { title: '商城管理' },
  },
  {
    path: '/game',
    name: 'GameManage',
    component: () => import('@/views/game/GameManage.vue'),
    meta: { title: '游戏管理' },
  },
  {
    path: '/logs',
    name: 'LogView',
    component: () => import('@/views/log/LogView.vue'),
    meta: { title: '系统日志' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
