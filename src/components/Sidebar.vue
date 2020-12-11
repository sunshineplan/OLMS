<template>
  <nav class="nav flex-column navbar-light sidebar">
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
                selected: $router.currentRoute.value.path == '/' && user.super,
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
</template>

<script>
export default {
  name: "Sidebar",
  data() {
    return {
      user: this.$store.state.user,
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
  },
  methods: {
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
.sidebar {
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

@media (min-width: 1201px) {
  .sidebar {
    display: block !important;
  }
}

@media (max-width: 1200px) {
  .sidebar {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
}
</style>
