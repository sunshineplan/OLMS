import { createStore } from 'vuex'
import { BootstrapButtons, post } from '../misc.js'

export default createStore({
  state() {
    return {
      user: {},
      sidebar: false,
      showSidebar: false,
      loading: 0,
      departments: [],
      employees: [],
      department: {},
      employee: {},
      record: {},
      filter: {},
      current: 1
    }
  },
  mutations: {
    user(state, user) { state.username = user },
    startLoading(state) { state.loading += 1 },
    stopLoading(state) { state.loading -= 1 },
    sidebar(state, status) { state.sidebar = status },
    closeSidebar(state) { state.showSidebar = false },
    toggleSidebar(state) { state.showSidebar = !state.showSidebar },
    departments(state, departments) { state.departments = departments },
    employees(state, employees) { state.employees = employees },
    department(state, department) { state.department = department },
    employee(state, employee) { state.employee = employee },
    record(state, record) { state.record = record },
    filter(state, filter) { state.filter = filter },
    current(state, current) { state.current = current }
  },
  actions: {
    async categories({ commit, state }) {
      commit('sidebar', false)
      commit('startLoading')
      const resp = await post('/category/get')
      commit('categories', await resp.json())
      if (state.category.count == undefined)
        commit('category', {
          id: -1,
          name: 'All Bookmarks',
          count: state.categories.reduce((total, i) => total + i.count, 0),
          start: 0
        })
      commit('stopLoading')
      commit('sidebar', true)
    },
    async bookmarks({ commit, state }, payload) {
      commit('startLoading')
      let params
      if (payload.more) {
        commit('more')
        params = { category: state.category.id, start: state.category.start }
      } else params = { category: payload.id }
      const resp = await post('/bookmark/get', params)
      if (!resp.ok)
        await BootstrapButtons.fire('Error', await resp.text(), 'error')
      else
        if (!payload.more) commit('bookmarks', await resp.json())
        else commit('bookmarks', state.bookmarks.concat(await resp.json()))
      commit('stopLoading')
    },
    reorder({ commit, state }, payload) {
      var bookmarks = state.bookmarks
      var bookmark = bookmarks[payload.old]
      bookmarks.splice(payload.old, 1)
      bookmarks.splice(payload.new, 0, bookmark)
      commit('bookmarks', bookmarks)
    },
    async addCategory({ dispatch, commit, state }, name) {
      await dispatch('categories')
      const category = state.categories.filter(i => i.name == name)
      if (category.length) {
        commit('category', category[0])
        commit('bookmarks', [])
      }
    },
    async editCategory({ dispatch, commit, state }, name) {
      await dispatch('categories')
      commit('category', {
        id: state.category.id,
        name,
        count: state.category.count,
        start: state.category.start
      })
      let bookmarks = state.bookmarks
      if (bookmarks)
        bookmarks.forEach(i => i.category = name)
      commit('bookmarks', bookmarks)
    },
    delBookmarks({ commit, state }, bookmark) {
      var categories = state.categories
      categories.forEach(i => { if (i.name == bookmark.category) i.count-- })
      commit('categories', categories)
      commit('bookmarks', state.bookmarks.filter(i => i.id != bookmark.id))
    }
  }
})
