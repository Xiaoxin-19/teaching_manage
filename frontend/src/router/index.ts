import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'dashboard', component: () => import('../views/Dashboard/Dashboard.vue') },
  { path: '/students', name: 'students', component: () => import('../views/StudentManager/StudentManage.vue') },
  { path: '/teachers', name: 'teachers', component: () => import('../views/TeacherManager/TeacherManage.vue') },
  { path: '/records', name: 'records', component: () => import('../views/RecordManager/RecordManage.vue') },
  { path: '/settings', name: 'settings', component: () => import('../views/StudentManager/StudentManage.vue') },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})