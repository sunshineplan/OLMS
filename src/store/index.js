import { createStore } from 'vuex'
import { post } from '../misc.js'

export default createStore({
  state() {
    return {
      user: null,
      showSidebar: false,
      loading: 0,
      departments: [],
      employees: [],
      recaptcha: null,
      years: [],
      department: {},
      employee: {},
      record: {},
      filter: {},
      page: 1
    }
  },
  mutations: {
    user(state, user) { state.username = user },
    startLoading(state) { state.loading += 1 },
    stopLoading(state) { state.loading -= 1 },
    closeSidebar(state) { state.showSidebar = false },
    toggleSidebar(state) { state.showSidebar = !state.showSidebar },
    departments(state, departments) { state.departments = departments },
    employees(state, employees) { state.employees = employees },
    recaptcha(state, recaptcha) { state.recaptcha = recaptcha },
    years(state, years) { state.years = years },
    department(state, department) { state.department = department },
    employee(state, employee) { state.employee = employee },
    record(state, record) { state.record = record },
    filter(state, filter) { state.filter = filter },
    page(state, page) { state.page = page }
  },
  actions: {
    async info({ commit }) {
      commit('startLoading')
      const resp = await fetch('/info')
      const json = await resp.json()
      let user = json.user
      if (user.id == 0) user.super = true
      commit('user', user)
      commit('departments', json.departments)
      commit('employees', json.employees)
      commit('recaptcha', json.recaptcha)
      commit('stopLoading')
    },
    delDepartment({ commit, state }, id) {
      commit('departments', state.departments.filter(i => i.id != id))
    },
    delEmployee({ commit, state }, id) {
      commit('employees', state.employees.filter(i => i.id != id))
    },
    reset({ commit }) {
      commit('filter', {})
      commit('page', 1)
    }
  }
})
