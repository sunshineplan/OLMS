<template>
  <header style="padding-left: 20px">
    <a class="h3 title">{{ mode }}</a>
    <hr />
  </header>
  <div class="form" @keyup.enter="save()">
    <div class="form-row" v-if="!personal">
      <div class="form-group">
        <label for="department">{{ $t("Department") }}</label>
        <select
          class="form-control"
          v-model.number="record.deptid"
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
        <select
          class="form-control"
          v-model.number="record.userid"
          id="employee"
          :disabled="!record.deptid"
          required
        >
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
        v-model="record.date"
        id="date"
        required
      />
    </div>
    <div class="form-row">
      <div class="form-group">
        <label for="type">{{ $t("Type") }}</label>
        <select class="form-control" v-model.number="record.type" id="type">
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
          v-model.number="record.duration"
          id="duration"
          required
        />
      </div>
    </div>
    <div class="form-group" v-if="user.super">
      <label for="status">{{ $t("Status") }}</label>
      <select class="form-control" v-model.number="record.status" id="status">
        <option value="0">{{ $t("Unverified") }}</option>
        <option value="1">{{ $t("Verified") }}</option>
        <option value="2">{{ $t("Rejected") }}</option>
      </select>
    </div>
    <div class="form-group">
      <label for="describe">{{ $t("Describe") }}</label>
      <textarea
        class="form-control"
        v-model="record.describe"
        id="describe"
        rows="3"
      />
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
import { post, valid } from "../misc.js";

export default {
  name: "Record",
  data() {
    return {
      user: this.$store.state.user,
      personal: this.$router.name == "departmentRecord" ? false : true,
      departments: this.$store.state.departments,
      mode:
        this.$route.params.mode == "add"
          ? this.$t("AddRecord")
          : this.$t("EditRecord"),
      record: this.$route.params.mode == "edit" ? this.$store.state.record : {},
      validated: false,
    };
  },
  computed: {
    employees() {
      return this.$store.state.employees.filter(
        (i) => i.deptid == this.record.deptid
      );
    },
  },
  mounted() {
    document.title = this.mode + " - " + this.$t("OLMS");
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
        if (this.user.super) data = this.record;
        const resp = await post(url, data);
        await this.checkResp(resp, async () => {
          await this.checkJson(await resp.json(), () => this.goback());
        });
      } else this.validated = true;
    },
    async del() {
      if (await this.confirm("Record")) {
        await this.checkResp(
          await post("/record/delete/" + this.record.id),
          () => this.goback()
        );
      }
    },
  },
};
</script>
