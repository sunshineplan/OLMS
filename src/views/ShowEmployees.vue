<template>
  <header style="padding-left: 20px">
    <div style="height: 50px">
      <a class="h3 title">{{ $t("EmployeesList") }}</a>
    </div>
    <div class="toolbar">
      <div class="form-inline">
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="department">
              {{ $t("Department") }}
            </label>
          </div>
          <select
            class="custom-select"
            v-model.number="filter.deptid"
            id="department"
          >
            <option value="0">{{ $t("All") }}</option>
            <option v-for="d in departments" :key="d.id" :value="d.id">
              {{ d.name }}
            </option>
          </select>
        </div>
        <div class="input-group input-group-sm" v-if="user.super">
          <div class="input-group-prepend">
            <label class="input-group-text" for="role">{{ $t("role") }}</label>
          </div>
          <select class="custom-select" v-model="filter.role" id="role">
            <option value="">{{ $t("All") }}</option>
            <option value="0">{{ $t("GeneralEmployee") }}</option>
            <option value="1">{{ $t("Administrator") }}</option>
          </select>
        </div>
        <div class="input-group">
          <a
            class="btn btn-primary btn-sm"
            @click="
              sort = {};
              $store.commit('sort', {});
              $store.commit('page', 1);
              doFilter();
            "
          >
            {{ $t("Filter") }}
          </a>
          <a
            class="btn btn-primary btn-sm"
            @click="
              reset();
              employees = $store.state.employees;
            "
          >
            {{ $t("Reset") }}
          </a>
        </div>
      </div>
    </div>
    <a class="btn btn-primary" @click="add()">{{ $t("Add") }}</a>
    <p></p>
  </header>
  <Pagination :total="total">
    <div class="table-responsive">
      <table class="table table-hover table-sm">
        <thead>
          <tr>
            <th
              v-for="i in field"
              :key="i"
              class="sortable"
              :class="
                sort.sort == i
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="doSort(i)"
            >
              {{ $t(i) }}
            </th>
            <th v-if="user.super">{{ $t("Operation") }}</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="e in employees.slice((page - 1) * 10, page * 10)"
            :key="e.id"
          >
            <td>{{ e.username }}</td>
            <td>{{ e.realname }}</td>
            <td>{{ e.deptname }}</td>
            <td v-if="user.super">
              {{ e.role ? $t("Administrator") : $t("GeneralEmployee") }}
            </td>
            <td v-if="user.super">{{ e.permission }}</td>
            <td v-if="user.super">
              <a class="btn btn-outline-primary btn-sm" @click="edit(e)">
                {{ $t("Edit") }}
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </Pagination>
</template>

<script>
import { defineAsyncComponent } from "vue";

export default {
  name: "ShowEmployees",
  components: {
    Pagination: defineAsyncComponent(() =>
      import(
        /* webpackChunkName: "show" */ "../components/Pagination.vue"
      )
    ),
  },
  data() {
    return {
      departments: this.$store.state.departments,
      employees: this.$store.state.employees,
      filter: this.$store.state.filter,
      sort: this.$store.state.sort,
    };
  },
  computed: {
    field() {
      const field = new Set([
        "username",
        "realname",
        "deptname",
        "role",
        "permission",
      ]);
      if (!this.user.super) {
        field.delete("role");
        field.delete("permission");
      }
      return field;
    },
    total() {
      return this.employees.length;
    },
    page() {
      return this.$store.state.page;
    },
  },
  created() {
    this.reset();
  },
  mounted() {
    document.title = this.$t("EmployeesList") + " - " + this.$t("OLMS");
  },
  methods: {
    doFilter() {
      this.$store.commit("filter", this.filter);
      if (!this.filter.deptid && !this.filter.role)
        this.employees = this.$store.state.employees;
      else if (this.filter.deptid && !this.filter.role)
        this.employees = this.$store.state.employees.filter(
          (i) => i.deptid == this.filter.deptid
        );
      else if (!this.filter.deptid && this.filter.role)
        this.employees = this.$store.state.employees.filter(
          (i) => Number(i.role) == this.filter.role
        );
      else
        this.employees = this.$store.state.employees.filter(
          (i) =>
            i.deptid == this.filter.deptid && Number(i.role) == this.filter.role
        );
    },
    doSort(field) {
      if (this.sort.sort == field && this.sort.order == "desc")
        this.sort.order = "asc";
      else this.sort = { sort: field, order: "desc" };
      this.$store.commit("sort", this.sort);
      this.employees.sort((a, b) => {
        if (a[field] == b[field]) return 0;
        if (this.sort.order == "desc")
          if (a[field] > b[field]) return 1;
          else return -1;
        else {
          if (a[field] > b[field]) return -1;
          else return 1;
        }
      });
    },
    add() {
      this.$router.push("/employee/add");
    },
    edit(employee) {
      this.$store.commit("employee", employee);
      this.$router.push("/employee/edit");
    },
    reset() {
      this.filter = { deptid: 0, role: "" };
      this.sort = {};
      this.$store.dispatch("reset", this.filter);
    },
  },
};
</script>
