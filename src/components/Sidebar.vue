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
                selected: $router.currentRoute.value.name == 'personalRecords',
              }"
              @click="
                $store.commit('personal', true);
                goto({ name: 'personalRecords' });
              "
            >
              {{ $t("EmployeeRecords") }}
            </a>
          </li>
          <li>
            <a
              class="nav-link"
              :class="{
                selected:
                  $router.currentRoute.value.name == 'personalStatistics',
              }"
              @click="
                $store.commit('personal', true);
                goto({ name: 'personalStatistics' });
              "
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
                  $router.currentRoute.value.name == 'departmentRecords' &&
                  !user.super,
              }"
              @click="
                $store.commit('personal', false);
                goto({ name: 'departmentRecords' });
              "
            >
              {{ $t("DepartmentRecords") }}
            </a>
          </li>
          <li>
            <a
              class="nav-link"
              :class="{
                selected:
                  $router.currentRoute.value.name == 'departmentStatistics',
              }"
              @click="
                $store.commit('personal', false);
                goto({ name: 'departmentStatistics' });
              "
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
                selected: $router.currentRoute.value.name == 'employees',
              }"
              @click="goto({ name: 'employees' })"
            >
              {{ $t("ManageEmployee") }}
            </a>
          </li>
          <li v-if="user.super">
            <a
              class="nav-link"
              :class="{
                selected: $router.currentRoute.value.name == 'departments',
              }"
              @click="goto({ name: 'departments' })"
            >
              {{ $t("ManageDepartment") }}
            </a>
          </li>
          <li v-if="user.super">
            <a
              class="nav-link"
              :class="{
                selected:
                  $router.currentRoute.value.name == 'departmentRecords' &&
                  user.super,
              }"
              @click="
                $store.commit('personal', false);
                goto({ name: 'departmentRecords' });
              "
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
    };
  },
  methods: {
    goto(router) {
      if (window.innerWidth <= 1200) this.$store.commit("closeSidebar");
      this.$router.push(router);
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
