import { createRouter, createWebHistory } from 'vue-router'


const routes = [
  {
    path: '/',
    component: () => import(/* webpackChunkName: 'records' */ '../views/ShowRecords.vue')
  },
  {
    path: '/departments',
    component: () => import(/* webpackChunkName: 'departments' */ '../views/ShowDepartments.vue')
  },
  {
    path: '/employees',
    component: () => import(/* webpackChunkName: 'employees' */ '../views/ShowEmployees.vue')
  },
  {
    path: '/statistics',
    component: () => import(/* webpackChunkName: 'statistics' */ '../views/ShowStatistics.vue')
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
    component: () => import(/* webpackChunkName: 'verify' */ '../views/Verify.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
