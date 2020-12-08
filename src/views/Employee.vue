<template>
  <header style="padding-left: 20px">
    <a class="h3 title">{{ mode }}</a>
    <hr />
  </header>
  <div class="form" @keyup.enter="save()">
    <div class="form-row">
      <div class="form-group">
        <label for="username">{{ $t("Username") }}</label>
        <input
          class="form-control"
          v-model.trim="employee.username"
          id="username"
          required
          autofocus
        />
        <div class="invalid-feedback">{{ $t("RequiredField") }}</div>
      </div>
      <div class="form-group">
        <label for="realname">{{ $t("Realname") }}</label>
        <input
          class="form-control"
          v-model.trim="employee.realname"
          id="realname"
        />
      </div>
    </div>
    <div class="form-group">
      <label for="department">{{ $t("Department") }}</label>
      <select
        class="form-control"
        v-model.number="employee.department"
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
        v-model.trim="employee.password"
        id="password"
        maxlength="20"
      />
      <small class="form-text text-muted">{{ $t("LeaveBlankPassword") }}</small>
    </div>
    <div class="form-group" v-if="user.super">
      <label for="role">{{ $t("Role") }}</label>
      <select
        class="form-control"
        v-model.number="employee.role"
        id="role"
        @change="if (!employee.role) permission = [];"
      >
        <option value="0">{{ $t("GeneralEmployee") }}</option>
        <option value="1">{{ $t("Administrator") }}</option>
      </select>
    </div>
    <div class="form-group" v-if="user.super" v-show="employee.role">
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
    <button class="btn btn-primary" @click="save()">
      {{ mode }}
    </button>
    <button class="btn btn-primary" @click="goback()">
      {{ $t("Cancel") }}
    </button>
  </div>
  <div class="form" v-if="$route.params.mode == 'edit'">
    <button class="btn btn-danger delete" @click="del()">
      {{ $t("Delete") }}
    </button>
  </div>
</template>

<script>
import { BootstrapButtons, post, valid } from "../misc.js";

export default {
  name: "Employee",
  data() {
    return {
      user: this.$store.state.user,
      departments: this.$store.state.departments,
      mode:
        this.$route.params.mode == "add"
          ? this.$t("AddEmployee")
          : this.$t("EditEmployee"),
      employee:
        this.$route.params.mode == "edit" ? this.$store.state.employee : {},
      permission: [],
      validated: false,
    };
  },
  created() {
    this.permission = this.employee.permission.split(",");
  },
  mounted() {
    document.title = this.mode;
    window.addEventListener("keyup", this.cancel);
  },
  beforeUnmount() {
    window.removeEventListener("keyup", this.cancel);
  },
  methods: {
    async save() {
      if (valid()) {
        this.validated = false;
        let url, data;
        if (this.$route.params.mode == "add") url = "/employee/add";
        else {
          url = "/employee/edit";
          data.id = this.employee.id;
        }
        data.username = this.employee.username;
        data.realname = this.employee.realname;
        data.deptid = this.employee.department;
        if (this.user.super) {
          data.password = this.employee.password;
          data.role = this.employee.role;
          data.permission = this.permission;
          if (data.role && !data.permission.length) {
            await BootstrapButtons.fire(
              this.$t("Error"),
              this.$t("EmptyPermission"),
              "error"
            );
            return;
          }
        }
        const resp = await post(url, data);
        await this.check(resp, async () => {
          await this.checkJson(await resp.json(), async () =>
            this.goback(true)
          );
        });
      } else this.validated = true;
    },
    async del() {
      if (await this.confirm("Employee")) {
        await this.checkResp(
          await post("/employee/delete/" + this.employee.id),
          async () => {
            await this.$store.dispatch("delEmployee", this.employee.id);
            this.goback();
          }
        );
      }
    },
  },
};
</script>
