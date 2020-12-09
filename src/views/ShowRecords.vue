<template>
  <header style="padding-left: 20px">
    <div style="height: 50px">
      <a class="h3 title">{{ mode }}</a>
    </div>
    <div class="toolbar">
      <div class="form-inline" v-if="!personal">
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
            @change="year('department')"
          >
            <option value="0">{{ $t("All") }}</option>
            <option v-for="d in departments" :key="d.id" :value="d.id">
              {{ d.name }}
            </option>
          </select>
        </div>
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="employee">
              {{ $t("Name") }}
            </label>
          </div>
          <select
            class="custom-select"
            v-model.number="filter.userid"
            id="employee"
            :disabled="!filter.deptid"
            @change="year('employee')"
          >
            <option value="0">{{ $t("All") }}</option>
            <option v-for="e in employees" :key="e.id" :value="e.id">
              {{ e.realname }}
            </option>
          </select>
        </div>
      </div>
      <div class="form-inline">
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="year">{{ $t("Year") }}</label>
          </div>
          <select class="custom-select" v-model="filter.year" id="year">
            <option value="">{{ $t("All") }}</option>
            <option v-for="y in years" :key="y" :value="y">{{ y }}</option>
          </select>
        </div>
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="month">
              {{ $t("Month") }}
            </label>
          </div>
          <select
            class="custom-select"
            v-model="filter.month"
            id="month"
            :disabled="filter.year == ''"
          >
            <option value="">{{ $t("All") }}</option>
            <option value="01">1</option>
            <option value="02">2</option>
            <option value="03">3</option>
            <option value="04">4</option>
            <option value="05">5</option>
            <option value="06">6</option>
            <option value="07">7</option>
            <option value="08">8</option>
            <option value="09">9</option>
            <option value="10">10</option>
            <option value="11">11</option>
            <option value="12">12</option>
          </select>
        </div>
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="type">{{ $t("Type") }}</label>
          </div>
          <select class="custom-select" v-model="filter.type" id="type">
            <option value="">{{ $t("All") }}</option>
            <option value="1">{{ $t("Overtime") }}</option>
            <option value="0">{{ $t("Leave") }}</option>
          </select>
        </div>
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="status">
              {{ $t("Status") }}
            </label>
          </div>
          <select class="custom-select" v-model="filter.status" id="status">
            <option value="">{{ $t("All") }}</option>
            <option value="0">{{ $t("Unverified") }}</option>
            <option value="1">{{ $t("Verified") }}</option>
            <option value="2">{{ $t("Rejected") }}</option>
          </select>
        </div>
      </div>
      <div class="form-inline">
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="describe">
              {{ $t("Describe") }}
            </label>
          </div>
          <input class="form-control" v-model="filter.describe" id="describe" />
        </div>
        <div class="input-group">
          <a
            class="btn btn-primary btn-sm"
            @click="
              this.sort = {};
              this.$store.commit('sort', {});
              this.$store.commit('page', 1);
              load('records');
            "
          >
            {{ $t("Filter") }}
          </a>
          <a class="btn btn-primary btn-sm" @click="reset('records')">
            {{ $t("Reset") }}
          </a>
          <a class="btn btn-info btn-sm" @click="download('records')">
            {{ $t("Export") }}
          </a>
        </div>
      </div>
    </div>
    <a class="btn btn-primary" @click="add(personal)">{{ $t("New") }}</a>
    <p></p>
  </header>
  <Pagination :total="total">
    <div class="table-responsive">
      <table class="table table-hover table-sm record">
        <thead>
          <tr>
            <th
              class="sortable"
              :class="
                sort.sort == 'deptname'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('deptname')"
              style="width: 150px"
              v-if="!personal"
            >
              {{ $t("Department") }}
            </th>
            <th
              class="sortable"
              :class="
                sort.sort == 'realname'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('realname')"
              style="width: 100px"
              v-if="!personal"
            >
              {{ $t("Realname") }}
            </th>
            <th
              class="sortable"
              :class="
                sort.sort == 'date'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('date')"
              style="width: 150px"
            >
              {{ $t("Date") }}
            </th>
            <th
              class="sortable"
              :class="
                sort.sort == 'type'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('type')"
              style="width: 80px"
            >
              {{ $t("Type") }}
            </th>
            <th
              class="sortable"
              @click="sortBy('duration')"
              style="width: 100px"
            >
              {{ $t("Duration") }}
            </th>
            <th
              class="describe sortable"
              :class="
                sort.sort == 'describe'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('describe')"
            >
              {{ $t("Describe") }}
            </th>
            <th
              class="sortable"
              :class="
                sort.sort == 'created'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('created')"
              style="width: 100px"
            >
              {{ $t("Created") }}
            </th>
            <th
              class="sortable"
              :class="
                sort.sort == 'status'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('status')"
              style="width: 100px"
            >
              {{ $t("Status") }}
            </th>
            <th style="width: 100px">{{ $t("Operation") }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in records" :key="r.id">
            <td v-if="!personal">{{ r.deptname }}</td>
            <td v-if="!personal">{{ r.realname }}</td>
            <td>
              {{
                r.date.replace(":00Z", "").replace(/-/g, "/").replace("T", " ")
              }}
            </td>
            <td>{{ r.type ? $t("Overtime") : $t("Leave") }}</td>
            <td>
              {{ r.duration }} {{ r.duration == 1 ? $t("Hour") : $t("Hours") }}
            </td>
            <td class="describe">{{ r.describe }}</td>
            <td>{{ r.created.split("T")[0] }}</td>
            <td>
              <a class="text-success" v-if="r.status == 1">
                {{ $t("Verified") }}
              </a>
              <a class="text-danger" v-else-if="r.status == 2">
                {{ $t("Rejected") }}
              </a>
              <a class="text-muted" v-else>{{ $t("Unverified") }}</a>
            </td>
            <td v-if="personal">
              <a
                class="btn btn-outline-primary btn-sm"
                :class="{ disabled: !r.status }"
                @click="edit(r, personal)"
              >
                {{ t("Edit") }}
              </a>
            </td>
            <td v-else-if="!user.super">
              <a
                class="btn btn-outline-primary btn-sm"
                :class="{ disabled: r.status }"
                @click="verify(r)"
              >
                {{ t("Verify") }}
              </a>
            </td>
            <td v-else>
              <a class="btn btn-outline-primary btn-sm" @click="edit(r)">
                {{ t("Edit") }}
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
  name: "ShowRecords",
  data() {
    return {
      user: this.$store.state.user,
      personal: this.$router.name == "departmentRecords" ? false : true,
      departments: this.$store.state.departments,
      years: [],
      records: [],
      total: 0,
      filter: this.$store.state.filter,
      sort: this.$store.state.sort,
    };
  },
  computed: {
    mode() {
      return this.personal
        ? this.$t("EmployeeRecords")
        : this.user.super
        ? this.$t("AllRecords")
        : this.$t("DepartmentRecords");
    },
    employees() {
      return this.$store.state.employees.filter(
        (i) => i.deptid == this.record.deptid
      );
    },
    page() {
      return this.$store.state.page;
    },
  },
  watch: {
    async sort() {
      await this.load("records");
    },
    async page() {
      await this.load("records");
    },
  },
  async created() {
    await this.year();
    await this.reset("records");
  },
  mounted() {
    document.title = this.mode + " - " + this.$t("OLMS");
  },
  methods: {
    add(personal) {
      if (personal)
        this.$router.push({ name: "personalRecord", params: { mode: "add" } });
      this.$router.push({ name: "departmentRecord", params: { mode: "add" } });
    },
    edit(record, personal) {
      this.$store.commit("record", record);
      if (personal)
        this.$router.push({ name: "personalRecord", params: { mode: "edit" } });
      this.$router.push({ name: "departmentRecord", params: { mode: "edit" } });
    },
    verify(record) {
      this.$store.commit("record", record);
      this.$router.push("/record/verify");
    },
  },
};
</script>

<style scoped>
@media (max-width: 1200px) {
  .describe {
    width: 150px;
    white-space: normal;
  }
}
</style>
