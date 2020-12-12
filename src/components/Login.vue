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
import { post } from "../misc.js";

const grecaptcha = window.grecaptcha;

export default {
  name: "Login",
  data() {
    return {
      user: this.$store.state.user,
      recaptcha: this.$store.state.recaptcha,
      username: "",
      password: "",
      rememberme: false,
      token: "",
    };
  },
  mounted() {
    document.title = this.$t("Login") + " - " + this.$t("OLMS");
    this.username = localStorage.getItem("username");
    if (this.recaptcha) {
      grecaptcha.ready(() => {
        this.execute();
        setInterval(() => {
          this.execute();
        }, 100000);
      });
    }
  },
  methods: {
    execute() {
      grecaptcha
        .execute(this.recaptcha, { action: "login" })
        .then((token) => (this.token = token));
    },
    async login() {
      if (!document.querySelector("#username").checkValidity())
        await this.prompt("Error", "EmptyUsername", "error");
      else if (!document.querySelector("#password").checkValidity())
        await this.prompt("Error", "EmptyPassword", "error");
      else {
        let data = {
          username: this.username,
          password: this.password,
          rememberme: this.rememberme,
        };
        if (this.recaptcha) data.recaptcha = this.token;
        const resp = await post("/login", data);
        this.checkResp(resp, () => {
          if (this.username != "root")
            localStorage.setItem("username", this.username);
          window.location = "/";
        });
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
