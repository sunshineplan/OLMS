import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import(/* webpackChunkName: 'show' */ '../views/ShowRecords.vue')
  },
  {
    path: '/statistics',
    component: () => import(/* webpackChunkName: 'show' */ '../views/ShowStatistics.vue')
  },
  {
    path: '/departments',
    component: () => import(/* webpackChunkName: 'manage' */ '../views/ShowDepartments.vue')
  },
  {
    path: '/employees',
    component: () => import(/* webpackChunkName: 'manage' */ '../views/ShowEmployees.vue')
  },
  {
    path: '/setting',
    component: () => import(/* webpackChunkName: 'config' */ '../views/Setting.vue')
  },
  {
    path: '/department/:mode',
    component: () => import(/* webpackChunkName: 'manage' */ '../views/Department.vue')
  },
  {
    path: '/employee/:mode',
    component: () => import(/* webpackChunkName: 'manage' */ '../views/Employee.vue')
  },
  {
    path: '/record/:mode',
    component: () => import(/* webpackChunkName: 'show' */ '../views/Record.vue')
  },
  {
    path: '/record/verify',
    component: () => import(/* webpackChunkName: 'manage' */ '../views/Verify.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
