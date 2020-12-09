import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    name: 'personalRecords',
    path: '/',
    component: () => import(/* webpackChunkName: 'records' */ '../views/ShowRecords.vue')
  },
  {
    name: 'personalStatistics',
    path: '/statistics',
    component: () => import(/* webpackChunkName: 'statistics' */ '../views/ShowStatistics.vue')
  },
  {
    name: 'departmentRecords',
    path: '/records',
    component: () => import(/* webpackChunkName: 'records' */ '../views/ShowRecords.vue')
  },
  {
    name: 'departmentStatistics',
    path: '/statistics',
    component: () => import(/* webpackChunkName: 'statistics' */ '../views/ShowStatistics.vue')
  },
  {
    name: 'departments',
    path: '/departments',
    component: () => import(/* webpackChunkName: 'departments' */ '../views/ShowDepartments.vue')
  },
  {
    name: 'employees',
    path: '/employees',
    component: () => import(/* webpackChunkName: 'employees' */ '../views/ShowEmployees.vue')
  },
  {
    path: '/setting',
    component: () => import(/* webpackChunkName: 'setting' */ '../views/Setting.vue')
  },
  {
    path: '/department/:mode',
    component: () => import(/* webpackChunkName: 'department' */ '../views/Department.vue')
  },
  {
    path: '/employee/:mode',
    component: () => import(/* webpackChunkName: 'employee' */ '../views/Employee.vue')
  },
  {
    name: 'personalRecord',
    path: '/record/:mode',
    component: () => import(/* webpackChunkName: 'record' */ '../views/Record.vue')
  },
  {
    name: 'departmentRecord',
    path: '/record/:mode',
    component: () => import(/* webpackChunkName: 'record' */ '../views/Record.vue')
  },
  {
    path: '/record/verify',
    component: () => import(/* webpackChunkName: 'verify' */ '../views/Verify.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
