import { createStore } from 'vuex'

export default createStore({
  state() {
    return {
      user: null,
      ready: false,
      showSidebar: false,
      loading: 0,
      departments: [],
      employees: [],
      recaptcha: null,
      personalRecord: null,
      personalStatistic: null,
      years: [],
      department: {},
      employee: {},
      record: {},
      filter: {},
      sort: {},
      page: 1
    }
  },
  mutations: {
    user(state, user) { state.user = user },
    ready(state) { state.ready = true },
    startLoading(state) { state.loading += 1 },
    stopLoading(state) { state.loading -= 1 },
    closeSidebar(state) { state.showSidebar = false },
    toggleSidebar(state) { state.showSidebar = !state.showSidebar },
    departments(state, departments) { state.departments = departments },
    employees(state, employees) { state.employees = employees },
    recaptcha(state, recaptcha) { state.recaptcha = recaptcha },
    personalRecord(state, personalRecord) { state.personalRecord = personalRecord },
    personalStatistic(state, personalStatistic) { state.personalStatistic = personalStatistic },
    years(state, years) { state.years = years },
    department(state, department) { state.department = department },
    employee(state, employee) { state.employee = employee },
    record(state, record) { state.record = record },
    filter(state, filter) { state.filter = filter },
    sort(state, sort) { state.sort = sort },
    page(state, page) { state.page = page }
  },
  actions: {
    async info({ commit }) {
      commit('startLoading')
      const resp = await fetch('/info')
      const json = await resp.json()
      if (json.recaptcha) {
        commit('recaptcha', json.recaptcha)
        delete json.recaptcha
      }
      if (Object.keys(json).length) {
        let user = json.user
        if (user.id == 0) user.super = true
        commit('user', user)
        commit('departments', json.departments)
        commit('employees', json.employees)
      } else {
        commit('user', false)
      }
      commit('stopLoading')
    },
    delDepartment({ commit, state }, id) {
      commit('departments', state.departments.filter(i => i.id != id))
    },
    delEmployee({ commit, state }, id) {
      commit('employees', state.employees.filter(i => i.id != id))
    },
    reset({ commit }, filter) {
      commit('filter', filter)
      commit('sort', {})
      commit('page', 1)
    }
  }
})
