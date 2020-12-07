export default {
  methods: {
    checkSize(size) {
      if (this.smallSize != window.innerWidth <= size)
        this.smallSize = !this.smallSize
    },
    async checkResp(resp, success) {
      if (resp.ok) return success()
      else return await BootstrapButtons.fire(this.$t('Error'), this.$t(await resp.text()), 'error')
    },
    async checkJson(json, success) {
      if (json.status) return success()
      else return await BootstrapButtons.fire(this.$t('Error'), this.$t(json.message), 'error')
    },
    async goback(reload) {
      if (reload)
        await this.$store.dispatch('info')
      this.$router.go(-1)
    },
    cancel(event) { if (event.key == 'Escape') this.goback() }
  }
}
