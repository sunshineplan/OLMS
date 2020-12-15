<template>
  <header style="padding-left: 20px">
    <div style="height: 50px">
      <a class="h3 title">{{ mode }}</a>
    </div>
    <Toolbar
      mode="statistics"
      :personal="personal"
      :filter="filter"
      :departments="departments"
      :employees="employees"
      :years="years"
      @update="updateFilter"
      @year="year"
      @filter="doFilter"
      @reset="reset"
    />
  </header>
  <Pagination :total="total">
    <div class="table-responsive">
      <table class="table table-hover table-sm">
        <thead>
          <tr>
            <th
              v-for="i in field"
              :key="i"
              class="sortable"
              :class="
                sort.sort == i
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default'
              "
              @click="sortBy(i)"
            >
              {{ $t(i) }}
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
import { defineAsyncComponent } from "vue";

export default {
  name: "ShowStatistics",
  components: {
    Toolbar: defineAsyncComponent(() =>
      import(/* webpackChunkName: "show" */ "../components/Toolbar.vue")
    ),
    Pagination: defineAsyncComponent(() =>
      import(/* webpackChunkName: "show" */ "../components/Pagination.vue")
    ),
  },
  data() {
    return {
      departments: this.$store.state.departments,
      years: [],
      statistics: [],
      total: 0,
      filter: this.$store.state.filter,
      sort: this.$store.state.sort,
    };
  },
  computed: {
    personal() {
      return this.$store.state.personalStatistic;
    },
    field() {
      const field = new Set([
        "period",
        "deptname",
        "realname",
        "overtime",
        "leave",
        "summary",
      ]);
      if (this.personal) {
        field.delete("deptname");
        field.delete("realname");
      }
      return field;
    },
    mode() {
      return this.personal
        ? this.$t("EmployeeStatistics")
        : this.$t("DepartmentStatistics");
    },
    employees() {
      if (this.personal == false)
        return this.$store.state.employees.filter(
          (i) => i.deptid == this.filter.deptid
        );
      return [];
    },
    page() {
      return this.$store.state.page;
    },
  },
  watch: {
    async personal() {
      document.title = this.mode + " - " + this.$t("OLMS");
      await this.year();
      await this.reset();
    },
    async sort(sort) {
      if (Object.keys(sort).length) {
        this.$store.commit("page", 1);
        await this.load("statistics");
      }
    },
    async page() {
      await this.load("statistics");
    },
  },
  async created() {
    await this.year();
    await this.reset();
  },
  mounted() {
    document.title = this.mode + " - " + this.$t("OLMS");
  },
  methods: {
    async doFilter() {
      this.sort = {};
      this.$store.commit("sort", {});
      this.$store.commit("page", 1);
      await this.load("statistics");
    },
    async reset() {
      this.filter = {
        deptid: 0,
        userid: 0,
        period: "month",
        year: "",
        month: "",
      };
      this.sort = {};
      this.$store.dispatch("reset", this.filter);
      await this.load("statistics");
    },
  },
};
</script>
