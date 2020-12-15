<template>
  <a
    class="toggle"
    v-if="user && smallSize"
    @click="toggle()"
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
  <transition name="slide">
    <nav
      class="nav flex-column navbar-light"
      v-show="showSidebar || !smallSize"
    >
      <div class="panel">
        <div v-if="!user.super">
          <a class="navbar-brand">{{ $t("EmployeePanel") }}</a>
          <ul class="navbar-nav">
            <li>
              <a
                class="nav-link"
                :class="{
                  selected:
                    $router.currentRoute.value.path == '/' && personalRecord,
                }"
                @click="goto('/', true)"
              >
                {{ $t("EmployeeRecords") }}
              </a>
            </li>
            <li>
              <a
                class="nav-link"
                :class="{
                  selected:
                    $router.currentRoute.value.path == '/statistics' &&
                    personalStatistic,
                }"
                @click="goto('/statistics', true)"
              >
                {{ $t("EmployeeStatistics") }}
              </a>
            </li>
          </ul>
        </div>
        <div v-if="user.role">
          <a class="navbar-brand">{{ $t("DepartmentPanel") }}</a>
          <ul class="navbar-nav">
            <li v-if="!user.super">
              <a
                class="nav-link"
                :class="{
                  selected:
                    $router.currentRoute.value.path == '/' &&
                    !personalRecord &&
                    !user.super,
                }"
                @click="goto('/', false)"
              >
                {{ $t("DepartmentRecords") }}
              </a>
            </li>
            <li>
              <a
                class="nav-link"
                :class="{
                  selected:
                    $router.currentRoute.value.path == '/statistics' &&
                    !personalStatistic,
                }"
                @click="goto('/statistics', false)"
              >
                {{ $t("DepartmentStatistics") }}
              </a>
            </li>
          </ul>
        </div>
        <div v-if="user.role">
          <a class="navbar-brand">{{ $t("ControlPanel") }}</a>
          <ul class="navbar-nav">
            <li>
              <a
                class="nav-link"
                :class="{
                  selected: $router.currentRoute.value.path == '/employees',
                }"
                @click="goto('/employees')"
              >
                {{ $t("ManageEmployee") }}
              </a>
            </li>
            <li v-if="user.super">
              <a
                class="nav-link"
                :class="{
                  selected: $router.currentRoute.value.path == '/departments',
                }"
                @click="goto('/departments')"
              >
                {{ $t("ManageDepartment") }}
              </a>
            </li>
            <li v-if="user.super">
              <a
                class="nav-link"
                :class="{
                  selected:
                    $router.currentRoute.value.path == '/' && user.super,
                }"
                @click="goto('/')"
              >
                {{ $t("ManageRecord") }}
              </a>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  </transition>
</template>

<script>
export default {
  name: "Sidebar",
  data() {
    return {
      hover: false,
      smallSize: window.innerWidth <= 1200,
    };
  },
  computed: {
    personalRecord() {
      return this.$store.state.personalRecord;
    },
    personalStatistic() {
      return this.$store.state.personalStatistic;
    },
    showSidebar() {
      return this.$store.state.showSidebar;
    },
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
        this.smallSize = window.innerWidth <= 1200;
    },
    toggle() {
      this.$store.commit("toggleSidebar");
    },
    goto(router, personal) {
      this.closeSidebar(() => {
        if (personal != undefined)
          if (router == "/") this.$store.commit("personalRecord", personal);
          else this.$store.commit("personalStatistic", personal);
        this.$router.push(router);
      });
    },
  },
};
</script>

<style scoped>
.toggle {
  position: fixed;
  z-index: 100;
  top: 0;
  padding: 20px;
  color: white !important;
}

.toggle:hover {
  background-color: rgb(232, 232, 232);
}

nav {
  position: fixed;
  top: 0;
  z-index: 1;
  height: 100%;
  width: 200px;
  padding-top: 70px;
  user-select: none;
}

.panel {
  height: 100%;
  width: 100%;
  padding-top: 10px;
  overflow-x: hidden;
  border-right: 1px solid #e9ecef;
  background-color: white;
}

.panel .navbar-brand {
  text-indent: 10px;
}

.panel .navbar-nav {
  text-indent: 20px;
}

.panel .nav-link {
  display: block;
  cursor: pointer;
  margin: 0;
  border-left: 5px solid transparent;
  color: rgba(0, 0, 0, 0.7) !important;
}

.panel .nav-link:hover {
  background-color: rgb(232, 232, 232) !important;
}

.selected {
  border-left: 5px solid #1a73e8 !important;
  color: #1a73e8 !important;
}

.panel .nav-link.selected {
  background-color: rgba(161, 194, 250, 0.16);
  color: #3367d6 !important;
}

.slide-leave-active,
.slide-enter-active {
  transition: 0.3s;
}

.slide-enter-from,
.slide-leave-to {
  transform: translate(-100%, 0);
}

@media (max-width: 1200px) {
  nav {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
}
</style>
