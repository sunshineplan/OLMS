<template>
  <header style="padding-left: 20px">
    <a class="h3 title">{{ mode }}</a>
    <hr />
  </header>
  <div class="form">
    <div class="form-group">
      <label for="department">{{ $t("Department") }}</label>
      <input
        class="form-control"
        v-model.trim="department.name"
        id="department"
        required
      />
      <div class="invalid-feedback">{{ $t("RequiredField") }}</div>
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
  name: "Department",
  data() {
    return {
      mode:
        this.$route.params.mode == "add"
          ? this.$t("AddDepartment")
          : this.$t("EditDepartment"),
      department:
        this.$route.params.mode == "edit" ? this.$store.state.department : {},
      validated: false,
    };
  },
  mounted() {
    document.title = this.mode;
    window.addEventListener("keyup", this.cancel);
  },
  beforeUnmount() {
    window.removeEventListener("keyup", this.cancel);
  },
  methods: {
    async save() {
      if (valid()) {
        this.validated = false;
        let resp;
        if (this.$route.params.mode == "add")
          resp = await post("/department/add", { name: this.department.name });
        else resp = await post("/department/edit", this.department);
        await this.checkResp(resp, async () => {
          await this.checkJson(await resp.json(), async () =>
            this.goback(true)
          );
        });
      } else this.validated = true;
    },
    async del() {
      if (await this.confirm("Department")) {
        await this.checkResp(
          await post("/department/delete/" + this.department.id),
          async () => {
            await this.$store.dispatch("delDepartment", this.department.id);
            this.goback();
          }
        );
      }
    },
  },
};
</script>
