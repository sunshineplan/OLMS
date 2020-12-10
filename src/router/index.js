import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import(/* webpackChunkName: 'record' */ '../views/ShowRecords.vue')
  },
  {
    path: '/statistics',
    component: () => import(/* webpackChunkName: 'statistic' */ '../views/ShowStatistics.vue')
  },
  {
    path: '/departments',
    component: () => import(/* webpackChunkName: 'department' */ '../views/ShowDepartments.vue')
  },
  {
    path: '/employees',
    component: () => import(/* webpackChunkName: 'employee' */ '../views/ShowEmployees.vue')
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
    path: '/record/:mode',
    component: () => import(/* webpackChunkName: 'record' */ '../views/Record.vue')
  },
  {
    path: '/record/verify',
    component: () => import(/* webpackChunkName: 'record' */ '../views/Verify.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
