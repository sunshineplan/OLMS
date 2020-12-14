<template>
  <component
    :is="script"
    src="https://www.recaptcha.net/recaptcha/api.js"
    v-if="recaptcha"
  />
  <nav class="navbar navbar-light topbar" v-if="user != null">
    <div class="d-flex" style="height: 100%">
      <a
        class="toggle"
        v-if="user && smallSize"
        @click="toggle"
        @mouseenter="hover = true"
        @mouseleave="hover = false"
      >
        <svg viewBox="0 0 70 70" width="40" height="30">
          <rect
            v-for="y in [10, 30, 50]"
            :key="y"
            :y="y"
            width="100%"
            height="10"
            :fill="hover ? '#1a73e8' : 'white'"
          />
        </svg>
      </a>
      <a class="brand" href="/" v-if="smallSize">{{ $t("OLMS") }}</a>
      <a class="brand" href="/" v-else>
        {{ $t("OvertimeAndLeaveManagementSystem") }}
      </a>
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
import { loadLocaleMessages } from "./i18n";
import Login from "./components/Login.vue";
import Sidebar from "./components/Sidebar.vue";

export default {
  name: "App",
  components: { Login, Sidebar },
  data() {
    return {
      hover: false,
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
    await loadLocaleMessages(this.$i18n.locale);
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

.toggle {
  padding: 20px;
  color: white !important;
}

.toggle:hover {
  background-color: rgb(232, 232, 232);
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

  .loading {
    left: 0;
    width: 100%;
  }
}
</style>
