<template>
  <header style="padding-left: 20px">
    <a class="h3 title">{{ mode }}</a>
    <hr />
  </header>
  <div class="form">
    <div class="form-row" v-if="user.admin || user.super">
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
      <div class="form-group">
        <label for="employee">{{ $t("Employee") }}</label>
        <select class="form-control" v-model="employee" id="employee" required>
          <option disabled>-- {{ $t("SelectEmployee") }} --</option>
          <option v-for="e in employees" :key="e.id" :value="e.id">
            {{ e.realname }}
          </option>
        </select>
      </div>
    </div>
    <div class="form-group">
      <label for="date">{{ $t("Date") }}</label>
      <input
        class="form-control"
        type="datetime-local"
        v-model="date"
        id="date"
        required
      />
    </div>
    <div class="form-row">
      <div class="form-group">
        <label for="type">{{ $t("Type") }}</label>
        <select class="form-control" v-model="type" id="type">
          <option value="1">{{ $t("Overtime") }}</option>
          <option value="0">{{ $t("Leave") }}</option>
        </select>
      </div>
      <div class="form-group">
        <label for="duration">{{ $t("Duration") }}</label>
        <input
          class="form-control"
          type="number"
          min="1"
          v-model="duration"
          id="duration"
          required
        />
      </div>
    </div>
    <div class="form-group" v-if="user.super">
      <label for="status">{{ $t("Status") }}</label>
      <select class="form-control" v-model="status" id="status">
        <option value="0">{{ $t("Unverified") }}</option>
        <option value="1">{{ $t("Verified") }}</option>
        <option value="2">{{ $t("Rejected") }}</option>
      </select>
    </div>
    <div class="form-group">
      <label for="describe">{{ $t("Describe") }}</label>
      <textarea
        class="form-control"
        v-model="describe"
        id="describe"
        rows="3"
      />
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
  name: "Record",
  data() {
    return {
      department: "",
      employee: "",
      date: "",
      type: "",
      duration: "",
      status: "",
      describe: "",
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
    employees() {
      return this.$store.state.employees;
    },
    record() {
      return this.$route.params.mode == "edit" ? this.$store.state.record : {};
    },
    mode() {
      return this.$route.params.mode == "add"
        ? this.$t("AddRecord")
        : this.$t("EditRecord");
    },
  },
  created() {
    this.department = this.record.department;
    this.employee = this.record.employee;
    this.date = this.record.date;
    this.type = this.record.type;
    this.duration = this.record.duration;
    this.status = this.record.status;
    this.describe = this.record.describe;
  },
  mounted() {
    document.title = this.mode + " Record";
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
        if (this.$route.params.mode == "add") url = "/record/add";
        else url = "/record/edit/" + this.record.id;
        let data = {
          department: this.record.department,
          employee: this.record.employee,
          date: this.record.date,
          type: this.record.type,
          duration: this.record.duration,
          describe: this.record.describe,
        };
        if (this.user.super) data.status = this.record.status;
        const resp = await post(url, data);
        if (!resp.ok)
          await BootstrapButtons.fire("Error", await resp.text(), "error");
        else {
          const json = await resp.json();
          if (json.status == 1) {
            if (this.$route.params.mode == "add")
              await this.$store.dispatch("addRecord", this.date);
            else await this.$store.dispatch("editRecord", this.date);
            this.goback();
          } else await BootstrapButtons.fire("Error", json.message, "error");
        }
      } else this.validated = true;
    },
    async del() {
      if (await confirm("employee")) {
        const resp = await post("/record/delete/" + this.record.id);
        if (!resp.ok)
          await BootstrapButtons.fire("Error", await resp.text(), "error");
        else this.goback();
      }
    },
  },
};
</script>
