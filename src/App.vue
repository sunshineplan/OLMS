<template>
  <nav class="navbar navbar-light topbar" v-if="user != null">
    <div class="d-flex" style="height: 100%">
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
    <Sidebar />
    <div
      class="content"
      style="padding-left: 200px"
      :style="{ opacity: loading ? 0.5 : 1 }"
      @mousedown="closeSidebar()"
    >
      <router-view v-if="!recaptcha || (recaptcha && ready)" />
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
import { defineAsyncComponent } from "vue";
import { loadLocaleMessages } from "./i18n";

export default {
  name: "App",
  components: {
    Login: defineAsyncComponent(() =>
      import(/* webpackChunkName: "config" */ "./components/Login.vue")
    ),
    Sidebar: defineAsyncComponent(() =>
      import(/* webpackChunkName: "show" */ "./components/Sidebar.vue")
    ),
  },
  computed: {
    loading() {
      return this.$store.state.loading;
    },
    ready() {
      return this.$store.state.ready;
    },
  },
  async created() {
    const lang = Cookies.get("lang") || navigator.language;
    await loadLocaleMessages(lang);
    document.querySelector("html").setAttribute("lang", this.$i18n.locale);
    await this.$store.dispatch("info");
    if (this.recaptcha) {
      const script = document.createElement("script");
      script.setAttribute(
        "src",
        "https://www.recaptcha.net/recaptcha/api.js?render=" + this.recaptcha
      );
      script.async = true;
      document.head.appendChild(script);
      setTimeout(() => {
        const grecaptcha = window.grecaptcha;
        if (grecaptcha)
          grecaptcha.ready(() => {
            this.$store.commit("ready");
          });
        else this.prompt("Error", "reCAPTCHALoadingFailed", "error");
      }, 500);
    }
    if (!this.user) this.$router.push("/login");
    else {
      const personal = this.user.role ? false : true;
      this.$store.commit("personalRecord", personal);
    }
  },
  methods: {
    setting() {
      this.closeSidebar(() => this.$router.push("/setting"));
    },
  },
};
</script>

<style>
@import "./style.css";
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

.short {
  display: none;
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

.loading {
  position: fixed;
  z-index: 2;
  top: 70px;
  left: 250px;
  height: calc(100% - 70px);
  width: calc(100% - 250px);
  display: flex;
}

@media (max-width: 1200px) {
  .short {
    display: inline;
  }

  .full {
    display: none;
  }

  .brand {
    padding-left: 90px;
  }

  .loading {
    left: 0;
    width: 100%;
  }
}
</style>
