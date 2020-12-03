<template>
  <header style="padding-left: 20px">
    <a class="h3 title">{{ mode }}</a>
    <hr />
  </header>
  <div class="form">
    <div class="form-row">
      <div class="form-group">
        <label for="username">{{ $t("Username") }}</label>
        <input
          class="form-control"
          v-model.trim="username"
          id="username"
          required
        />
        <div class="invalid-feedback">{{ $t("RequiredField") }}</div>
      </div>
      <div class="form-group">
        <label for="realname">{{ $t("Realname") }}</label>
        <input class="form-control" v-model.trim="realname" id="realname" />
      </div>
    </div>
    <div class="form-group">
      <label for="department">{{ $t("Department") }}</label>
      <select
        class="form-control"
        v-model="department"
        id="department"
        required
      >
        <option disabled>-- {{ $t("SelectDepartment") }} --</option>
        <option v-for="d in departments" :key="d.id" :value="d.id">
          {{ d.name }}
        </option>
      </select>
    </div>
    <div class="form-group" v-if="$route.params.mode == 'edit' && user.super">
      <label for="password">{{ $t("Password") }}</label>
      <input
        class="form-control"
        type="password"
        v-model.trim="password"
        id="password"
        maxlength="20"
      />
      <small class="form-text text-muted">{{ $t("LeaveBlankPassword") }}</small>
    </div>
    <div class="form-group" v-if="user.super">
      <label for="role">{{ $t("Role") }}</label>
      <select class="form-control" v-model="role" id="role">
        <option value="0">{{ $t("GeneralEmployee") }}</option>
        <option value="1">{{ $t("Administrator") }}</option>
      </select>
    </div>
    <div
      class="form-group"
      id="permission-selector"
      v-if="user.super"
      v-show="role == 1"
    >
      <label for="permission">
        {{ $t("Permission") }} ({{ $t("MultipleChoice") }})
      </label>
      <select
        multiple
        class="form-control"
        v-model="permission"
        id="permission"
      >
        <option v-for="d in departments" :key="d.id" :value="d.id"></option>
      </select>
    </div>
    <button class="btn btn-primary" :click="save()">
      {{ mode }}
    </button>
    <button class="btn btn-primary" :click="goback()">
      {{ $t("Cancel") }}
    </button>
  </div>
  <div class="form" v-if="$route.params.mode == 'edit'">
    <button class="btn btn-danger delete" :click="del()">
      {{ $t("Delete") }}
    </button>
  </div>
</template>

<script>
import { BootstrapButtons, post, valid, confirm } from "../misc.js";

export default {
  name: "Employee",
  data() {
    return {
      username: "",
      realname: "",
      department: "",
      password: "",
      role: false,
      permission: [],
      validated: false,
    };
  },
  computed: {
    user() {
      return this.$store.state.user;
    },
    departments() {
      return this.$store.state.departments;
    },
    employee() {
      return this.$route.params.mode == "edit"
        ? this.$store.state.employee
        : {};
    },
    mode() {
      return this.$route.params.mode == "add"
        ? this.$t("AddEmployee")
        : this.$t("EditEmployee");
    },
  },
  created() {
    this.username = this.employee.username;
    this.realname = this.employee.realname;
    this.department = this.employee.department;
    this.role = this.employee.role;
    this.permission = this.employee.permission;
  },
  mounted() {
    document.title = this.mode + " Employee";
    window.addEventListener("keyup", this.cancel);
  },
  beforeUnmount() {
    window.removeEventListener("keyup", this.cancel);
  },
  methods: {
    async save() {
      if (valid()) {
        this.validated = false;
        let url;
        if (this.$route.params.mode == "add") url = "/employee/add";
        else url = "/employee/edit/" + this.record.id;
        let data = {
          username: this.username,
          realname: this.realname,
          department: this.department,
        };
        if (this.user.super) {
          data.password = this.password;
          data.role = this.role;
          data.permission = this.permission;
        }
        const resp = await post(url, data);
        if (!resp.ok)
          await BootstrapButtons.fire("Error", await resp.text(), "error");
        else {
          const json = await resp.json();
          if (json.status == 1) {
            if (this.$route.params.mode == "add")
              await this.$store.dispatch("addEmployee", this.realname);
            else await this.$store.dispatch("editEmployee", this.realname);
            this.goback();
          } else await BootstrapButtons.fire("Error", json.message, "error");
        }
      } else this.validated = true;
    },
    async del() {
      if (await confirm("employee")) {
        const resp = await post("/employee/delete/" + this.employee.id);
        if (!resp.ok)
          await BootstrapButtons.fire("Error", await resp.text(), "error");
        else this.goback();
      }
    },
  },
};
</script>
