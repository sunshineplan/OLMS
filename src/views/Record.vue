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
          <option value="0" disabled>-- {{ $t("SelectDepartment") }} --</option>
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
          <option value="0" disabled>-- {{ $t("SelectEmployee") }} --</option>
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
        <select class="form-control" v-model="record.type" id="type">
          <option :value="true">{{ $t("Overtime") }}</option>
          <option :value="false">{{ $t("Leave") }}</option>
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
      personal: this.$store.state.personalRecord,
      departments: this.$store.state.departments,
      mode:
        this.$route.params.mode == "add"
          ? this.$t("AddRecord")
          : this.$t("EditRecord"),
      record:
        this.$route.params.mode == "edit"
          ? this.$store.state.record
          : { deptid: 0, userid: 0, type: true, status: 0 },
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
  created() {
    if (this.record.date)
      this.record.date = this.record.date.replace(":00Z", "");
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
        const data = {
          date: this.record.date + ":00Z",
          type: this.record.type,
          duration: this.record.duration,
          describe: this.record.describe,
        };
        if (!this.personal) {
          data.deptid = this.record.deptid;
          data.userid = this.record.userid;
        }
        if (this.user.super) data.status = this.record.status;
        let url;
        if (this.$route.params.mode == "add") url = "/record/add";
        else {
          url = "/record/edit";
          data.id = this.record.id;
        }
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

<style scoped>
.form-control {
  width: 250px !important;
}
</style>
