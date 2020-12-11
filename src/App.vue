<template>
  <component
    :is="script"
    src="https://www.recaptcha.net/recaptcha/api.js"
    v-if="recaptcha"
  />
  <nav class="navbar navbar-light topbar" v-if="user != null">
    <div class="d-flex" style="height: 100%">
      <a class="toggle" v-if="user && smallSize" @click="toggle">
        <i class="material-icons menu">menu</i>
      </a>
      <a class="brand full" href="/">
        {{ $t("OvertimeAndLeaveManagementSystem") }}
      </a>
      <a class="brand short" href="/">{{ $t("OLMS") }}</a>
    </div>
    <div class="navbar-nav flex-row" v-if="user">
      <a class="nav-link" v-text="user.realname"></a>
      <a class="nav-link link" @click="setting()">
        {{ $t("Setting") }}
      </a>
      <a class="nav-link link" href="/logout">{{ $t("Logout") }}</a>
    </div>
    <div class="navbar-nav flex-row" v-else>
      <a class="nav-link">{{ $t("Login") }}</a>
    </div>
  </nav>
  <Login v-if="user == false" />
  <div v-else-if="user">
    <transition name="slide">
      <Sidebar v-show="showSidebar || !smallSize" />
    </transition>
    <div
      class="content"
      style="padding-left: 200px"
      :style="{ opacity: loading ? 0.5 : 1 }"
      @mousedown="closeSidebar()"
    >
      <router-view />
    </div>
  </div>
  <div class="loading" v-show="loading">
    <div class="sk-wave sk-center">
      <div class="sk-wave-rect"></div>
      <div class="sk-wave-rect"></div>
      <div class="sk-wave-rect"></div>
      <div class="sk-wave-rect"></div>
      <div class="sk-wave-rect"></div>
    </div>
  </div>
</template>

<script>
import Cookies from "js-cookie";
import i18n, { loadLocaleMessages } from "./i18n";
import Login from "./components/Login.vue";
import Sidebar from "./components/Sidebar.vue";

export default {
  name: "App",
  components: { Login, Sidebar },
  data() {
    return {
      smallSize: window.innerWidth <= 1200,
    };
  },
  computed: {
    user() {
      return this.$store.state.user;
    },
    recaptcha() {
      return this.$store.state.recaptcha;
    },
    loading() {
      return this.$store.state.loading;
    },
    showSidebar() {
      return this.$store.state.showSidebar;
    },
  },
  async created() {
    const lang = Cookies.get("lang");
    if (lang) this.$i18n.locale = lang;
    await loadLocaleMessages(i18n, this.$i18n.locale);
    document.querySelector("html").setAttribute("lang", this.$i18n.locale);
    await this.$store.dispatch("info");
    if (!this.user) this.$router.push("/login");
    else {
      const personal = this.user.role ? false : true;
      this.$store.commit("personalRecord", personal);
    }
  },
  mounted() {
    window.addEventListener("resize", this.checkSize);
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.checkSize);
  },
  methods: {
    checkSize() {
      if (this.smallSize != window.innerWidth <= 1200)
        this.smallSize = !this.smallSize;
    },
    setting() {
      this.closeSidebar(() => this.$router.push("/setting"));
    },
    toggle() {
      this.$store.commit("toggleSidebar");
    },
  },
};
</script>

<style>
:root {
  --sk-color: #1a73e8;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
    "Helvetica Neue", Arial, "Noto Sans", "Microsoft YaHei New",
    "Microsoft Yahei", 微软雅黑, 宋体, SimSun, STXihei, 华文细黑, sans-serif;
}

a:hover {
  text-decoration: none;
}

.form,
.lang,
.subscribe {
  padding: 0 20px;
}

.h3 {
  cursor: default;
}

.input-group {
  padding: 5px 10px 5px 0px;
  max-width: 240px;
}

.form-control {
  width: 250px;
}

.content {
  position: fixed;
  top: 0;
  padding-top: 90px;
  height: 100%;
  width: 100%;
  overflow-y: auto;
}

.toolbar {
  padding-bottom: 10px;
}

.table-responsive {
  min-height: 300px;
  padding: 0 10px;
  cursor: default;
}

table,
table.record {
  table-layout: fixed;
}

td {
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

td:hover {
  white-space: normal;
}

.btn + .btn {
  margin-left: 4px;
}

.btn-info {
  margin-left: 8px !important;
}

.form-row {
  margin-right: 0px;
  margin-left: 0px;
}

.form-row .form-group {
  padding-right: 30px;
}

.delete {
  margin-top: 8px;
}

.swal {
  margin: 8px 6px;
}

#department {
  width: 135px;
}

#employee {
  width: 95px;
}

#year {
  width: 70px;
}

.sortable {
  cursor: pointer;
  background-position: right;
  background-repeat: no-repeat;
  padding-right: 30px !important;
}

.default {
  background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAQAAADYWf5HAAAAkElEQVQoz7X QMQ5AQBCF4dWQSJxC5wwax1Cq1e7BAdxD5SL+Tq/QCM1oNiJidwox0355mXnG/DrEtIQ6azioNZQxI0ykPhTQIwhCR+BmBYtlK7kLJYwWCcJA9M4qdrZrd8pPjZWPtOqdRQy320YSV17OatFC4euts6z39GYMKRPCTKY9UnPQ6P+GtMRfGtPnBCiqhAeJPmkqAAAAAElFTkSuQmCC");
}

.asc {
  background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAYAAAByUDbMAAAAZ0lEQVQ4y2NgGLKgquEuFxBPAGI2ahhWCsS/gDibUoO0gPgxEP8H4ttArEyuQYxAPBdqEAxPBImTY5gjEL9DM+wTENuQahAvEO9DMwiGdwAxOymGJQLxTyD+jgWDxCMZRsEoGAVoAADeemwtPcZI2wAAAABJRU5ErkJggg==);
}

.desc {
  background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAYAAAByUDbMAAAAZUlEQVQ4y2NgGAWjYBSggaqGu5FA/BOIv2PBIPFEUgxjB+IdQPwfC94HxLykus4GiD+hGfQOiB3J8SojEE9EM2wuSJzcsFMG4ttQgx4DsRalkZENxL+AuJQaMcsGxBOAmGvopk8AVz1sLZgg0bsAAAAASUVORK5CYII=);
}

@media (max-width: 1200px) {
  .content {
    padding-left: 0 !important;
  }

  table {
    table-layout: auto;
  }
}
</style>

<style scoped>
.topbar {
  position: fixed;
  top: 0px;
  z-index: 2;
  width: 100%;
  height: 70px;
  padding: 0 10px 0 0;
  background-color: #1a73e8;
  user-select: none;
}

.topbar .nav-link {
  padding-left: 8px;
  padding-right: 8px;
  color: white !important;
}

.topbar .link:hover {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 5px;
  cursor: pointer;
}

.toggle {
  padding: 20px;
  color: white !important;
}

.toggle:hover {
  color: #1a73e8 !important;
  background-color: rgb(232, 232, 232);
}

.material-icons.menu {
  font-size: 30px;
}

.brand {
  padding-left: 20px;
  margin: auto;
  font-size: 25px;
  letter-spacing: 0.3px;
  color: white;
}

.brand:hover {
  color: white;
  text-decoration: none;
}

.short {
  display: none;
}

.loading {
  position: fixed;
  z-index: 2;
  top: 70px;
  left: 250px;
  height: calc(100% - 70px);
  width: calc(100% - 250px);
  display: flex;
}

.slide-leave-active,
.slide-enter-active {
  transition: 0.5s;
}

.slide-enter-from,
.slide-leave-to {
  transform: translate(-100%, 0);
}

@media (max-width: 1200px) {
  .brand {
    padding-left: 10px;
  }

  .short {
    display: inline;
  }

  .full {
    display: none;
  }

  .loading {
    left: 0;
    width: 100%;
  }
}
</style>
