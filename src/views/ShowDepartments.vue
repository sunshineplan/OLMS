<template>
  <header style="padding-left: 20px">
    <div style="height: 50px">
      <a class="h3 title">{{ $t("DepartmentsList") }}</a>
    </div>
    <a class="btn btn-primary" :click="add()">{{ $t("Add") }}</a>
    <p></p>
  </header>

  <div class="table-responsive">
    <table class="table table-hover">
      <thead>
        <tr>
          <th>ID</th>
          <th>{{ $t("Department") }}</th>
          <th>{{ $t("Operation") }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="d in departments" :key="d.id">
          <td>{{ d.id }}</td>
          <td>{{ d.name }}</td>
          <td>
            <a class="btn btn-outline-primary btn-sm" :click="edit(d)">
              {{ $t("Edit") }}
            </a>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: "ShowDepartments",
  data() {
    return {
      departments: [],
    };
  },
  async created() {
    await load();
  },
  methods: {
    async load() {
      this.$store.commit("loading");
      const resp = await fetch("/departments");
      this.departments = await resp.json();
      this.$store.commit("loading");
    },
    add() {
      this.$store.commit("department", {});
      this.$router.push("/department/add");
    },
    edit(department) {
      this.$store.commit("department", department);
      this.$router.push("/department/edit");
    },
  },
};
</script>
