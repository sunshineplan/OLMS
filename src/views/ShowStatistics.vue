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
            <label class="input-group-text" for="period">
              {{ $t("Period") }}
            </label>
          </div>
          <select
            class="custom-select"
            v-model="filter.period"
            id="period"
            @change="
              filter.year = '';
              filter.month = '';
            "
          >
            <option value="month">{{ $t("Month") }}</option>
            <option value="year">{{ $t("Year") }}</option>
          </select>
        </div>
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
            v-show="filter.period == 'month'"
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
        <div class="input-group">
          <a
            class="btn btn-primary btn-sm"
            @click="
              this.sort = {};
              this.$store.commit('sort', {});
              this.$store.commit('page', 1);
              load('statistics');
            "
          >
            {{ $t("Filter") }}
          </a>
          <a class="btn btn-primary btn-sm" @click="reset('statistics')">
            {{ $t("Reset") }}
          </a>
          <a class="btn btn-info btn-sm" @click="download('statistics')">
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
            <th
              class="sortable"
              :class="
                sort.sort == 'period'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('period')"
            >
              {{ $t("Period") }}
            </th>
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
              v-if="!personal"
            >
              {{ $t("Realname") }}
            </th>
            <th
              class="sortable"
              :class="
                sort.sort == 'overtime'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('overtime')"
            >
              {{ $t("Overtime") }}
            </th>
            <th
              class="sortable"
              :class="
                sort.sort == 'leave'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('leave')"
            >
              {{ $t("Leave") }}
            </th>
            <th
              class="sortable"
              :class="
                sort.sort == 'summary'
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy('summary')"
            >
              {{ $t("Summary") }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in statistics" :key="s.period + s.deptname + s.realname">
            <td>{{ s.period }}</td>
            <td v-if="!personal">{{ s.deptname }}</td>
            <td v-if="!personal">{{ s.realname }}</td>
            <td>
              {{ s.overtime }}
              {{
                s.overtime == 0 || s.overtime == 1 ? $t("Hour") : $t("Hours")
              }}
            </td>
            <td>
              {{ s.leave }}
              {{ s.leave == 0 || s.leave == 1 ? $t("Hour") : $t("Hours") }}
            </td>
            <td>
              {{ s.summary }}
              {{
                s.summary == 0 || Math.abs(s.summary) == 1
                  ? $t("Hour")
                  : $t("Hours")
              }}
            </td>
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
      user: this.$store.state.user,
      personal: this.$router.name == "departmentStatistics" ? false : true,
      departments: this.$store.state.departments,
      years: [],
      statistics: [],
      total: 0,
      filter: this.$store.state.filter,
      sort: this.$store.state.sort,
    };
  },
  computed: {
    mode() {
      return this.personal
        ? this.$t("EmployeeStatistics")
        : this.$t("DepartmentStatistics");
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
      this.$store.commit("page", 1);
      await this.load("records");
    },
    async page() {
      await this.load("statistics");
    },
  },
  async created() {
    await this.year();
    await this.reset("statistics");
  },
  mounted() {
    document.title = this.mode + " - " + this.$t("OLMS");
  },
};
</script>
