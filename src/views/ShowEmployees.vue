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
              {{ $t("Dept") }}
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
            <label class="input-group-text" for="type">{{ $t("Type") }}</label>
          </div>
          <select class="custom-select" v-model="filter.type" id="type">
            <option value="">{{ $t("All") }}</option>
            <option value="0">{{ $t("GeneralEmployee") }}</option>
            <option value="1">{{ $t("Administrator") }}</option>
          </select>
        </div>
        <div class="input-group">
          <a class="btn btn-primary btn-sm" @click="doFilter()">
            {{ $t("Filter") }}
          </a>
          <a class="btn btn-primary btn-sm" @click="$store.dispatch('reset')">
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
            <th class="sortable" data-name="username">{{ $t("Username") }}</th>
            <th class="sortable" data-name="realname">{{ $t("Realname") }}</th>
            <th class="sortable" data-name="deptname">
              {{ $t("Department") }}
            </th>
            <th class="sortable" data-name="role" v-if="user.super">
              {{ $t("Role") }}
            </th>
            <th class="sortable" data-name="permission" v-if="user.super">
              {{ $t("Permission") }}
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
export default {
  name: "ShowEmployees",
  data() {
    return {
      user: this.$store.state.user,
      departments: this.$store.state.departments,
      employees: this.$store.state.employees,
      filter: this.$store.state.filter,
    };
  },
  computed: {
    total() {
      return this.employees.lenght;
    },
    page() {
      return this.state.page;
    },
  },
  created() {
    this.filter = {};
  },
  methods: {
    async doFilter() {
      this.$store.commit("filter", this.filter);
      if (!this.filter.deptid) this.employees = this.$store.state.employees;
      else
        this.employees = this.$store.state.employees.filter(
          (i) => i.deptid == this.filter.deptid
        );
    },
    add() {
      this.$router.push("/employee/add");
    },
    edit(employee) {
      this.$store.commit("employee", employee);
      this.$router.push("/employee/edit");
    },
  },
};
</script>
