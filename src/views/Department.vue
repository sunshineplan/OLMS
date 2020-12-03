<template>
  <header style="padding-left: 20px">
    <a class="h3 title">{{ mode }}</a>
    <hr />
  </header>
  <div class="form">
    <div class="form-group">
      <label for="dept">{{ $t("Department") }}</label>
      <input class="form-control" v-model.trim="name" id="dept" required />
      <div class="invalid-feedback">{{ $t("RequiredField") }}</div>
    </div>
    <button class="btn btn-primary" :click="save()">
      {{ mode }}
    </button>
    <button class="btn btn-primary" :click="goback()">
      {{ $t("Cancel") }}
    </button>
  </div>
  <div class="form" v-if="$route.params.mode == 'edit'">
    <button class="btn btn-danger delete" :click="del()">
      {{ $t("Delete") }}
    </button>
  </div>
</template>

<script>
import { BootstrapButtons, post, valid, confirm } from "../misc.js";

export default {
  name: "Department",
  data() {
    return {
      name: "",
      validated: false,
    };
  },
  computed: {
    department() {
      return this.$route.params.mode == "edit"
        ? this.$store.state.department
        : {};
    },
    mode() {
      return this.$route.params.mode == "add"
        ? this.$t("AddDepartment")
        : this.$t("EditDepartment");
    },
  },
  created() {
    this.name = this.department.name;
  },
  mounted() {
    document.title = this.mode + " Department";
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
          resp = await post("/department/add", { name: this.name });
        else
          resp = await post("/department/edit/" + this.department.id, {
            name: this.name,
          });
        if (!resp.ok)
          await BootstrapButtons.fire("Error", await resp.text(), "error");
        else {
          const json = await resp.json();
          if (json.status == 1) {
            if (this.$route.params.mode == "add")
              await this.$store.dispatch("addDepartment", this.name);
            else await this.$store.dispatch("editDepartment", this.name);
            this.goback();
          } else await BootstrapButtons.fire("Error", json.message, "error");
        }
      } else this.validated = true;
    },
    async del() {
      if (await confirm("department")) {
        const resp = await post("/department/delete/" + this.department.id);
        if (!resp.ok)
          await BootstrapButtons.fire("Error", await resp.text(), "error");
        else this.goback();
      }
    },
  },
};
</script>
