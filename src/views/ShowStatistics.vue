<template>
  <header style="padding-left: 20px">
    <div style="height: 50px">
      <a class="h3 title" v-if="user.role">{{ $t("DepartmentStatistics") }}</a>
      <a class="h3 title" v-else>{{ $t("EmployeeStatistics") }}</a>
    </div>
    <div class="toolbar">
      <div class="form-inline" v-if="user.role">
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="department">
              {{ $t("Dept") }}
            </label>
          </div>
          <select class="custom-select" v-model="department" id="department">
            <option value="">{{ $t("All") }}</option>
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
          <select class="custom-select" v-model="employee" id="employee">
            <option value="">{{ $t("All") }}</option>
            <option v-for="e in employees" :key="e.id" :value="e.id">
              {{ e.realname }}
            </option>
          </select>
        </div>
      </div>
      <div class="form-inline">
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="period">
              {{ $t("Period") }}
            </label>
          </div>
          <select class="custom-select" v-model="period" id="period">
            <option value="month">{{ $t("Month") }}</option>
            <option value="year">{{ $t("Year") }}</option>
          </select>
        </div>
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="year">{{ $t("Year") }}</label>
          </div>
          <select class="custom-select" v-model="year" id="year">
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
            v-model="month"
            id="month"
            v-show="period == 'month'"
            :disabled="year == ''"
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
        <div class="input-group">
          <a class="btn btn-primary btn-sm" :click="filter()">
            {{ $t("Filter") }}
          </a>
          <a class="btn btn-primary btn-sm" :click="reset()">
            {{ $t("Reset") }}
          </a>
          <a class="btn btn-info btn-sm" :click="download()">
            {{ $t("Export") }}
          </a>
        </div>
      </div>
    </div>
  </header>
  <Pagination :total="total">
    <div class="table-responsive">
      <table class="table table-hover table-sm">
        <thead>
          <tr>
            <th class="sortable" data-name="period">{{ $t("Period") }}</th>
            <th class="sortable" data-name="dept_name" v-if="user.role">
              {{ $t("Department") }}
            </th>
            <th class="sortable" data-name="realname" v-if="user.role">
              {{ $t("Realname") }}
            </th>
            <th class="sortable" data-name="overtime">{{ $t("Overtime") }}</th>
            <th class="sortable" data-name="leave">{{ $t("Leave") }}</th>
            <th class="sortable" data-name="summary">{{ $t("Summary") }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in statistics" :key="s.period + s.deptname + s.realname">
            <td>{{ s.period }}</td>
            <td v-if="user.role">{{ s.deptname }}</td>
            <td v-if="user.role">{{ s.realname }}</td>
            <td>{{ s.overtime }} {{ $t("Hours") }}</td>
            <td>{{ s.leave }} {{ $t("Hours") }}</td>
            <td>{{ s.summary }} {{ $t("Hours") }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </Pagination>
</template>

<script>
export default {
  name: "ShowStatistics",
  data() {
    return {
      statistics: [],
      total: 0,
      department: "",
      employee: "",
      period: "month",
      year: "",
      month: "",
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
  },
  async created() {
    await load();
  },
  methods: {
    async load() {
      this.$store.commit("loading");
      const resp = await fetch("/statistics");
      json = await resp.json();
      this.statistics = json.rows;
      this.total = json.total;
      this.$store.commit("loading");
    },
    async filter() {},
    async reset() {
      this.department = "";
      this.employee = "";
      this.period = "month";
      this.year = "";
      this.month = "";
      await load();
    },
    download() {},
  },
};
</script>
