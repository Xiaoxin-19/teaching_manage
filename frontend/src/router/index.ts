import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'dashboard', component: () => import('../views/StudentManage.vue') },
  { path: '/students', name: 'students', component: () => import('../views/StudentManage.vue') },
  { path: '/teachers', name: 'teachers', component: () => import('../views/StudentManage.vue') },
  { path: '/records', name: 'records', component: () => import('../views/StudentManage.vue') },
  { path: '/settings', name: 'settings', component: () => import('../views/StudentManage.vue') },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
