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
            @change="
              filter.userid = 0;
              year('department');
            "
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
            <option v-for="y in years" :key="y" :value="String(y)">
              {{ y }}
            </option>
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
            <label class="input-group-text" for="type">{{ $t("type") }}</label>
          </div>
          <select class="custom-select" v-model="filter.type" id="type">
            <option value="">{{ $t("All") }}</option>
            <option value="1">{{ $t("overtime") }}</option>
            <option value="0">{{ $t("leave") }}</option>
          </select>
        </div>
        <div class="input-group input-group-sm">
          <div class="input-group-prepend">
            <label class="input-group-text" for="status">
              {{ $t("status") }}
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
              {{ $t("describe") }}
            </label>
          </div>
          <input class="form-control" v-model="filter.describe" id="describe" />
        </div>
        <div class="input-group">
          <a
            class="btn btn-primary btn-sm"
            @click="
              sort = {};
              $store.commit('sort', {});
              $store.commit('page', 1);
              load('records');
            "
          >
            {{ $t("Filter") }}
          </a>
          <a class="btn btn-primary btn-sm" @click="reset()">
            {{ $t("Reset") }}
          </a>
          <a class="btn btn-info btn-sm" @click="download('records')">
            {{ $t("Export") }}
          </a>
        </div>
      </div>
    </div>
    <a class="btn btn-primary" @click="add()">{{ $t("New") }}</a>
    <p></p>
  </header>
  <Pagination :total="total">
    <div class="table-responsive">
      <table class="table table-hover table-sm record">
        <thead>
          <tr>
            <th
              v-for="(v, k) in field"
              :key="k"
              class="sortable"
              :class="[
                sort.sort == k
                  ? sort.order == 'desc'
                    ? 'desc'
                    : 'asc'
                  : 'default',
                { describe: k == 'describe' },
              ]"
              @click="sortBy(k)"
              :style="{ width: v.width }"
            >
              {{ $t(k) }}
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
            <td>{{ r.type ? $t("overtime") : $t("leave") }}</td>
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
                :class="{ disabled: r.status }"
                @click="edit(r)"
              >
                {{ $t("Edit") }}
              </a>
            </td>
            <td v-else-if="!user.super">
              <a
                class="btn btn-outline-primary btn-sm"
                :class="{ disabled: r.status }"
                @click="verify(r)"
              >
                {{ $t("Verify") }}
              </a>
            </td>
            <td v-else>
              <a class="btn btn-outline-primary btn-sm" @click="edit(r)">
                {{ $t("Edit") }}
              </a>
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
  name: "ShowRecords",
  components: {
    Pagination: defineAsyncComponent(() =>
      import(
        /* webpackChunkName: "show" */ "../components/Pagination.vue"
      )
    ),
  },
  data() {
    return {
      departments: this.$store.state.departments,
      years: [],
      records: [],
      total: 0,
      filter: this.$store.state.filter,
      sort: this.$store.state.sort,
    };
  },
  computed: {
    personal() {
      return this.$store.state.personalRecord;
    },
    field() {
      const field = {
        deptname: { width: "150px", personal: false },
        realname: { width: "100px", personal: false },
        date: { width: "150px", personal: true },
        type: { width: "80px", personal: true },
        duration: { width: "100px", personal: true },
        describe: { personal: true },
        created: { width: "100px", personal: true },
        status: { width: "100px", personal: true },
      };
      if (this.personal) {
        delete field.deptname;
        delete field.realname;
      }
      return field;
    },
    mode() {
      return this.personal
        ? this.$t("EmployeeRecords")
        : this.user.super
        ? this.$t("AllRecords")
        : this.$t("DepartmentRecords");
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
    async personal(to, from) {
      if (from != null) {
        document.title = this.mode + " - " + this.$t("OLMS");
        await this.year();
        await this.reset();
      }
    },
    async sort(sort) {
      if (Object.keys(sort).length) {
        this.$store.commit("page", 1);
        await this.load("records");
      }
    },
    async page() {
      await this.load("records");
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
    add() {
      this.$router.push("/record/add");
    },
    edit(record) {
      this.$store.commit("record", record);
      this.$router.push("/record/edit");
    },
    verify(record) {
      this.$store.commit("record", record);
      this.$router.push("/record/verify");
    },
    async reset() {
      this.filter = {
        deptid: 0,
        userid: 0,
        year: "",
        month: "",
        type: "",
        status: "",
        describe: "",
      };
      this.sort = {};
      this.$store.dispatch("reset", this.filter);
      await this.load("records");
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
