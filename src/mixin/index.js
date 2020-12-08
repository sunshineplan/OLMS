import Swal from 'sweetalert2'
import { BootstrapButtons, post } from '../misc.js'

export default {
  methods: {
    async confirm(type) {
      return await Swal.fire({
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
      }).isConfirmed
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
      this.$store.commit('loading')
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
      this.$store.commit('loading')
    },
    async load(mode) {
      this.$store.commit('loading')
      if (this.personal) this.filter.personal = true
      this.filter.page = this.page
      this.$store.commit('filter', this.filter)
      const resp = await post(`/${mode}`, this.filter)
      const json = await resp.json()
      this.records = json.rows
      this.total = json.total
      this.$store.commit('loading')
    },
    async download(mode) {
      this.$store.commit('loading')
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
            .replace(/'/g, '')
        )
        link.click()
      }
      this.$store.commit('loading')
    },
    async goback(reload) {
      if (reload)
        await this.$store.dispatch('info')
      this.$router.go(-1)
    },
    cancel(event) { if (event.key == 'Escape') this.goback() },
    async reset(mode) {
      this.filter = {}
      await this.load(mode)
    }
  }
}
