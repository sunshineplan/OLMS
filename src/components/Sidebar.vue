<template>
  <nav class="nav flex-column navbar-light sidebar">
    <div class="panel">
      <div v-if="!user.super">
        <a class="navbar-brand">{{ $t("EmployeePanel") }}</a>
        <ul class="navbar-nav">
          <li>
            <router-link
              class="nav-link"
              :class="{ selected: $router.name == 'personalRecords' }"
              :to="{ name: 'personalRecords' }"
            >
              {{ $t("EmployeeRecords") }}
            </router-link>
          </li>
          <li>
            <router-link
              class="nav-link"
              :class="{ selected: $router.name == 'personalStatistics' }"
              :to="{ name: 'personalStatistics' }"
            >
              {{ $t("EmployeeStatistics") }}
            </router-link>
          </li>
        </ul>
      </div>
      <div v-if="user.role">
        <a class="navbar-brand">{{ $t("DepartmentPanel") }}</a>
        <ul class="navbar-nav">
          <li v-if="!user.super">
            <router-link
              class="nav-link"
              :class="{
                selected: $router.name == 'departmentRecords' && !user.super,
              }"
              :to="{ name: 'departmentRecords' }"
            >
              {{ $t("DepartmentRecords") }}
            </router-link>
          </li>
          <li>
            <router-link
              class="nav-link"
              :class="{ selected: $router.name == 'departmentStatistics' }"
              :to="{ name: 'departmentStatistics' }"
            >
              {{ $t("DepartmentStatistics") }}
            </router-link>
          </li>
        </ul>
      </div>
      <div v-if="user.role">
        <a class="navbar-brand">{{ $t("ControlPanel") }}</a>
        <ul class="navbar-nav">
          <li>
            <router-link
              class="nav-link"
              :class="{ selected: $router.name == 'employees' }"
              :to="{ name: 'employees' }"
            >
              {{ $t("ManageEmployee") }}
            </router-link>
          </li>
          <li v-if="user.super">
            <router-link
              class="nav-link"
              :class="{ selected: $router.name == 'departments' }"
              :to="{ name: 'departments' }"
            >
              {{ $t("ManageDepartment") }}
            </router-link>
          </li>
          <li v-if="user.super">
            <router-link
              class="nav-link"
              :class="{
                selected: $router.name == 'departmentRecords' && user.super,
              }"
              to="/records"
            >
              {{ $t("ManageRecord") }}
            </router-link>
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
    };
  },
  methods: {
    closeSidebar() {
      if (window.innerWidth <= 1200) this.$store.commit("closeSidebar");
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
    display: none;
  }
}
</style>
