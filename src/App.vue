<template>
  <component
    :is="script"
    src="https://www.recaptcha.net/recaptcha/api.js"
    v-if="recaptcha"
  ></component>
  <nav class="navbar navbar-light topbar">
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
      <router-link class="nav-link link" to="/setting">
        {{ $t("Setting") }}
      </router-link>
      <a class="nav-link link" href="/logout">{{ $t("Logout") }}</a>
    </div>
    <div class="navbar-nav flex-row" v-else>
      <a class="nav-link">{{ $t("Login") }}</a>
    </div>
  </nav>
  <Login v-if="!user" />
  <div v-else>
    <transition name="slide">
      <Sidebar v-show="showSidebar || !smallSize" />
    </transition>
    <div
      class="content"
      style="padding-left: 200px"
      :style="{ opacity: loading ? 0.5 : 1 }"
      @mousedown="closeSidebar"
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
import Login from "./components/Login.vue";
import Sidebar from "./components/Sidebar.vue";
export default {
  name: "App",
  components: { Login, Sidebar },
  data() {
    return {
      user: this.$store.state.user,
      recaptcha: this.$store.state.recaptcha,
      smallSize: window.innerWidth <= 900,
    };
  },
  computed: {
    loading() {
      return this.$store.state.loading;
    },
    showSidebar() {
      return this.$store.state.showSidebar;
    },
  },
  async created() {
    await this.$store.dispatch("info");
    if (!this.user) this.$router.push("/login");
    else this.$router.push({ name: "departmentRecords" });
  },
  mounted() {
    window.addEventListener("resize", this.checkSize900);
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.checkSize900);
  },
  methods: {
    checkSize900() {
      this.checkSize(900);
    },
    toggle() {
      this.$store.commit("toggleSidebar");
    },
    closeSidebar() {
      if (this.smallSize) this.$store.commit("closeSidebar");
    },
  },
};
</script>
