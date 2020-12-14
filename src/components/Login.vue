<template>
  <div class="content">
    <header>
      <h3
        class="d-flex justify-content-center align-items-center"
        style="height: 100%"
      >
        {{ $t("Login") }}
      </h3>
    </header>
    <div class="login" @keyup.enter="login()">
      <div class="form-group">
        <label for="username">{{ $t("Username") }}</label>
        <input
          autofocus
          class="form-control"
          v-model.trim="username"
          id="username"
          maxlength="20"
          :placeholder="$t('Username')"
          required
        />
      </div>
      <div class="form-group">
        <label for="password">{{ $t("Password") }}</label>
        <input
          class="form-control"
          type="password"
          v-model.trim="password"
          id="password"
          maxlength="20"
          :placeholder="$t('Password')"
          required
        />
      </div>
      <div class="form-group form-check">
        <input
          type="checkbox"
          class="form-check-input"
          v-model="rememberme"
          id="rememberme"
        />
        <label class="form-check-label" for="rememberme">
          {{ $t("RememberMe") }}
        </label>
      </div>
      <hr />
      <button class="btn btn-primary login" @click="login()">
        {{ $t("Login") }}
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: "Login",
  data() {
    return {
      user: this.$store.state.user,
      username: "",
      password: "",
      rememberme: false,
    };
  },
  mounted() {
    document.title = this.$t("Login") + " - " + this.$t("OLMS");
    this.username = localStorage.getItem("username");
  },
  methods: {
    async login() {
      if (!document.querySelector("#username").checkValidity())
        await this.prompt("Error", "EmptyUsername", "error");
      else if (!document.querySelector("#password").checkValidity())
        await this.prompt("Error", "EmptyPassword", "error");
      else {
        const data = {
          username: this.username,
          password: this.password,
          rememberme: this.rememberme,
        };
        const login = async () => {
          const resp = await this.post("/login", data, "login");
          this.checkResp(resp, () => {
            if (this.username != "root")
              localStorage.setItem("username", this.username);
            window.location = "/";
          });
        };
        if (!this.recaptcha) login();
        else window.grecaptcha.ready(login);
      }
    },
  },
};
</script>

<style scoped>
.login {
  width: 250px;
  margin: 0 auto 20px;
}
</style>
