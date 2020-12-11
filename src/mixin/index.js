import Swal from 'sweetalert2'
import { BootstrapButtons, post } from '../misc.js'

export default {
  methods: {
    async confirm(type) {
      const confirm = await Swal.fire({
        title: this.$t('AreYouSure'),
        text: this.$t('Delete' + type),
        icon: 'warning',
        confirmButtonText: this.$t('Delete'),
        showCancelButton: true,
        focusCancel: true,
        customClass: {
          confirmButton: 'swal btn btn-danger',
          cancelButton: 'swal btn btn-primary'
        },
        buttonsStyling: false
      })
      return confirm.isConfirmed
    },
    closeSidebar(func) {
      if (func) {
        if (this.smallSize) {
          this.$store.commit('closeSidebar')
          setTimeout(() => func(), 500)
        } else (
          func()
        )
      } else if (this.smallSize) this.$store.commit('closeSidebar')
    },
    async checkResp(resp, success) {
      if (resp.ok) return success()
      else
        return await BootstrapButtons.fire(this.$t('Error'), this.$t(await resp.text()), 'error')
    },
    async checkJson(json, success) {
      if (json.status) return success()
      else
        return await BootstrapButtons.fire(this.$t('Error'), this.$t(json.message), 'error')
    },
    async year(mode) {
      this.$store.commit('startLoading')
      let resp
      if (this.personal) resp = await fetch('/year')
      else {
        let data
        if (mode == 'department') data = { deptid: this.filter.deptid }
        else if (mode == 'employee') data = { userid: this.filter.userid }
        else data = {}
        resp = await post('/year', data)
      }
      const json = await resp.json()
      if (json.year === 0) this.years = []
      else
        this.years = Array(new Date().getFullYear() - json.year + 1).fill().map((_, i) => json.year + i)
      this.$store.commit('stopLoading')
    },
    sortBy(field) {
      if (this.sort.sort == field && this.sort.order == 'desc')
        this.sort = { sort: field, order: 'asc' }
      else this.sort = { sort: field, order: 'desc' }
      this.$store.commit('sort', this.sort)
    },
    async load(mode) {
      this.$store.commit('startLoading')
      if (this.personal) this.filter.personal = true
      if (Object.keys(this.sort).length) {
        this.filter.sort = this.sort.sort
        this.filter.order = this.sort.order
      }
      this.filter.page = this.page
      this.$store.commit('filter', this.filter)
      const resp = await post(`/${mode}`, this.filter)
      const json = await resp.json()
      this[mode] = json.rows
      this.total = json.total
      this.$store.commit('stopLoading')
    },
    async download(mode) {
      this.$store.commit('startLoading')
      this.$store.commit('filter', this.filter)
      const resp = await post(`/${mode}/export`, this.filter)
      if (resp.status == 404)
        await BootstrapButtons.fire(this.$t('Info'), this.$t('NoResult'), 'info')
      else {
        const blob = new Blob(
          [new Uint8Array([0xef, 0xbb, 0xbf]), await resp.blob()],
          { type: 'text/csv;charset=utf-8' }
        )
        let link = document.createElement('a')
        link.href = window.URL.createObjectURL(blob)
        link.download = decodeURI(
          resp.headers
            .get('Content-Disposition')
            .split('filename=')[1]
            .replace(/"/g, '')
        )
        link.click()
      }
      this.$store.commit('stopLoading')
    },
    async goback(reload) {
      if (reload)
        await this.$store.dispatch('info')
      this.$router.go(-1)
    },
    cancel(event) { if (event.key == 'Escape') this.goback() }
  }
}
