<template>
  <div class="toolbar">
    <div class="form-inline" v-if="!personal || mode == 'employees'">
      <div class="input-group input-group-sm">
        <div class="input-group-prepend">
          <label class="input-group-text" for="department">
            {{ $t("Department") }}
          </label>
        </div>
        <select
          class="custom-select"
          :value="filter.deptid"
          id="department"
          @change="
            $emit('update', 'deptid', Number($event.target.value));
            $emit('update', 'userid', 0);
            $emit('year', 'department');
          "
        >
          <option value="0">{{ $t("All") }}</option>
          <option v-for="d in departments" :key="d.id" :value="d.id">
            {{ d.name }}
          </option>
        </select>
      </div>
      <div
        class="input-group input-group-sm"
        v-if="user.super && mode == 'employees'"
      >
        <div class="input-group-prepend">
          <label class="input-group-text" for="role">{{ $t("role") }}</label>
        </div>
        <select
          class="custom-select"
          :value="filter.role"
          id="role"
          @change="$emit('update', 'role', $event.target.value)"
        >
          <option value="">{{ $t("All") }}</option>
          <option value="0">{{ $t("GeneralEmployee") }}</option>
          <option value="1">{{ $t("Administrator") }}</option>
        </select>
      </div>
      <div class="input-group" id="employees" v-if="mode == 'employees'" />
      <div class="input-group input-group-sm" v-else>
        <div class="input-group-prepend">
          <label class="input-group-text" for="employee">
            {{ $t("Name") }}
          </label>
        </div>
        <select
          class="custom-select"
          :value="filter.userid"
          id="employee"
          :disabled="!filter.deptid"
          @change="
            $emit('update', 'userid', Number($event.target.value));
            $emit('year', 'employee');
          "
        >
          <option value="0">{{ $t("All") }}</option>
          <option v-for="e in employees" :key="e.id" :value="e.id">
            {{ e.realname }}
          </option>
        </select>
      </div>
    </div>
    <div class="form-inline" v-if="mode != 'employees'">
      <div class="input-group input-group-sm" v-if="mode == 'statistics'">
        <div class="input-group-prepend">
          <label class="input-group-text" for="period">
            {{ $t("period") }}
          </label>
        </div>
        <select
          class="custom-select"
          :value="filter.period"
          id="period"
          @change="
            $emit('update', 'period', $event.target.value);
            $emit('update', 'year', '');
            $emit('update', 'month', '');
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
        <select
          class="custom-select"
          :value="filter.year"
          id="year"
          @change="$emit('update', 'year', $event.target.value)"
        >
          <option value="">{{ $t("All") }}</option>
          <option v-for="y in years" :key="y" :value="String(y)">
            {{ y }}
          </option>
        </select>
      </div>
      <div
        class="input-group input-group-sm"
        v-show="!filter.period || filter.period == 'month'"
      >
        <div class="input-group-prepend">
          <label class="input-group-text" for="month">
            {{ $t("Month") }}
          </label>
        </div>
        <select
          class="custom-select"
          :value="filter.month"
          id="month"
          :disabled="filter.year == ''"
          @change="$emit('update', 'month', $event.target.value)"
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
      <div class="input-group input-group-sm" v-if="mode == 'records'">
        <div class="input-group-prepend">
          <label class="input-group-text" for="type">{{ $t("type") }}</label>
        </div>
        <select
          class="custom-select"
          :value="filter.type"
          id="type"
          @change="$emit('update', 'type', $event.target.value)"
        >
          <option value="">{{ $t("All") }}</option>
          <option value="1">{{ $t("overtime") }}</option>
          <option value="0">{{ $t("leave") }}</option>
        </select>
      </div>
      <div class="input-group input-group-sm" v-if="mode == 'records'">
        <div class="input-group-prepend">
          <label class="input-group-text" for="status">
            {{ $t("status") }}
          </label>
        </div>
        <select
          class="custom-select"
          :value="filter.status"
          id="status"
          @change="$emit('update', 'status', $event.target.value)"
        >
          <option value="">{{ $t("All") }}</option>
          <option value="0">{{ $t("Unverified") }}</option>
          <option value="1">{{ $t("Verified") }}</option>
          <option value="2">{{ $t("Rejected") }}</option>
        </select>
      </div>
      <div class="input-group" id="statistics" v-if="mode == 'statistics'" />
    </div>
    <div class="form-inline" v-if="mode == 'records'">
      <div class="input-group input-group-sm">
        <div class="input-group-prepend">
          <label class="input-group-text" for="describe">
            {{ $t("describe") }}
          </label>
        </div>
        <input
          class="form-control"
          :value="filter.describe"
          id="describe"
          @input="$emit('update', 'describe', $event.target.value)"
        />
      </div>
      <div class="input-group" id="records" />
    </div>
  </div>
  <teleport :to="`#${mode}`" v-if="isMounted">
    <a class="btn btn-primary btn-sm" @click="$emit('filter')">
      {{ $t("Filter") }}
    </a>
    <a class="btn btn-primary btn-sm" @click="$emit('reset')">
      {{ $t("Reset") }}
    </a>
    <a
      class="btn btn-info btn-sm"
      @click="download(mode)"
      v-if="mode != 'employees'"
    >
      {{ $t("Export") }}
    </a>
  </teleport>
</template>

<script>
export default {
  name: "Toolbar",
  emits: ["update", "year", "filter", "reset"],
  props: {
    mode: String,
    personal: Boolean,
    filter: Object,
    departments: Array,
    employees: Array,
    years: Array,
  },
  data() {
    return {
      isMounted: false,
    };
  },
  mounted() {
    this.isMounted = true;
  },
};
</script>

<style scoped>
.toolbar {
  padding-bottom: 10px;
}

.input-group {
  padding: 5px 10px 5px 0px;
  max-width: 240px;
}

#department {
  width: 135px;
}

#employee {
  width: 95px;
}

#year {
  width: 70px;
}
</style>
