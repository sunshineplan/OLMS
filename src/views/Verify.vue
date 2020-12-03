<template>
  <header style="padding-left: 20px">
    <a class="h3 title">{{ $t("VerifyRecord") }}</a>
    <hr />
  </header>
  <div class="form">
    <div class="form-row">
      <div class="form-group">
        <label for="department">{{ $t("Department") }}</label>
        <input
          class="form-control"
          v-text="department"
          id="department"
          readonly
        />
      </div>
      <div class="form-group">
        <label for="employee">{{ $t("Employee") }}</label>
        <input class="form-control" v-text="employee" id="employee" readonly />
      </div>
    </div>
    <div class="form-group">
      <label for="date">{{ $t("Date") }}</label>
      <input class="form-control" v-text="date" id="date" readonly />
    </div>
    <div class="form-group">
      <label for="type">{{ $t("Type") }}</label>
      <input class="form-control" v-text="type" id="type" readonly />
    </div>
    <div class="form-group">
      <label for="duration">{{ $t("Duration") }}</label>
      <input class="form-control" v-text="duration" id="duration" readonly />
    </div>
    <div class="form-group">
      <label for="describe">{{ $t("Describe") }}</label>
      <textarea
        class="form-control"
        v-text="describe"
        id="describe"
        rows="3"
        readonly
      />
    </div>
    <button class="btn btn-success" :click="verify(true)">
      {{ $t("Accept") }}
    </button>
    <button class="btn btn-danger" :click="verify(false)">
      {{ $t("Reject") }}
    </button>
    <button class="btn btn-primary" :click="goback()">
      {{ $t("Cancel") }}
    </button>
  </div>
</template>

<script>
import { BootstrapButtons, post, confirm } from "../misc.js";

export default {
  name: "Verify",
  data() {
    return {
      department: "",
      employee: "",
      date: "",
      type: "",
      duration: "",
      describe: "",
    };
  },
  computed: {
    record() {
      return this.$store.state.record;
    },
  },
  created() {
    this.department = this.record.department;
    this.employee = this.record.employee;
    this.date = this.record.date;
    this.type = this.record.type;
    this.duration = this.record.duration;
    this.describe = this.record.describe;
  },
  mounted() {
    document.title = "Verify Record";
    window.addEventListener("keyup", this.cancel);
  },
  beforeUnmount() {
    window.removeEventListener("keyup", this.cancel);
  },
  methods: {
    async verify(status) {
      const resp = await post("/record/verify/" + this.record.id, { status });
      if (!resp.ok)
        await BootstrapButtons.fire("Error", await resp.text(), "error");
      else {
        const json = await resp.json();
        if (json.status == 1) this.goback();
        else await BootstrapButtons.fire("Error", json.message, "error");
      }
    },
  },
};
</script>
