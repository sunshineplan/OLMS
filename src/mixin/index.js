import Swal from 'sweetalert2'

export default {
  computed: {
    user() { return this.$store.state.user },
    recaptcha() { return this.$store.state.recaptcha }
  },
  methods: {
    async post(url, data, challenge) {
      const json = {}
      if (data) Object.keys(data).forEach(key => data[key] !== '' && (json[key] = data[key]))
      if (this.recaptcha && challenge)
        json.recaptcha = await window.grecaptcha.execute(this.recaptcha, { action: challenge });
      return fetch(url, {
        method: 'post',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(json)
      })
    },
    prompt(title, content, icon) {
      const prompt = Swal.mixin({
        customClass: { confirmButton: 'swal btn btn-primary' },
        confirmButtonText: this.$t('OK'),
        buttonsStyling: false
      })
      return prompt.fire(this.$t(title), this.$t(content), icon)
    },
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
        if (window.innerWidth > 1200 || !this.$store.state.showSidebar)
          func()
        else {
          this.$store.commit('closeSidebar')
          setTimeout(func, 300)
        }
      } else this.$store.commit('closeSidebar')
    },
    async checkResp(resp, success) {
      if (resp.ok) return success()
      else
        return await this.prompt('Error', await resp.text(), 'error')
    },
    async checkJson(json, success) {
      if (json.status) return success()
      else
        return await this.prompt('Error', json.message, 'error')
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
        resp = await this.post('/year', data)
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
      const resp = await this.post(`/${mode}`, this.filter, mode)
      const json = await resp.json()
      this[mode] = json.rows
      this.total = json.total
      this.$store.commit('stopLoading')
    },
    async download(mode) {
      this.$store.commit('startLoading')
      this.$store.commit('filter', this.filter)
      const resp = await this.post(`/${mode}/export`, this.filter, mode)
      if (resp.status == 404)
        await this.prompt('Info', 'NoResult', 'info')
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
