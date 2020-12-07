import { BootstrapButtons, post } from '../misc.js'

export default {
  methods: {
    checkSize(size) {
      if (this.smallSize != window.innerWidth <= size)
        this.smallSize = !this.smallSize
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
    async load(mode) {
      this.$store.commit('loading')
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
        await BootstrapButtons.fire(this.$t('Info'), this.$t(NoResult), 'info')
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
