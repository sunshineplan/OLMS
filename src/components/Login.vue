<template>
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
        placeholder='{{$t("Username")}}'
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
        placeholder='{{$t("Password")}}'
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
    <input
      type="hidden"
      name="g-recaptcha-response"
      v-model="token"
      v-if="recaptcha"
    />
    <hr />
    <button class="btn btn-primary login" @click="login()">
      {{ $t("Login") }}
    </button>
  </div>
</template>

<script>
import { BootstrapButtons, post } from "../misc.js";

export default {
  name: "Login",
  data() {
    return {
      username: "",
      password: "",
      rememberme: false,
      token: "",
    };
  },
  mounted() {
    document.title = "Log In";
    this.username = localStorage.getItem("username");
    if (this.$store.state.recaptcha) {
      function execute() {
        grecaptcha
          .execute(this.$store.state.recaptcha, { action: "login" })
          .then((token) => (this.token = token));
      }
      grecaptcha.ready(() => {
        execute();
        setInterval(() => {
          execute();
        }, 100000);
      });
    }
  },
  methods: {
    async login() {
      if (!document.querySelector("#username").checkValidity())
        await BootstrapButtons.fire(
          "Error",
          "Username cannot be empty.",
          "error"
        );
      else if (!document.querySelector("#password").checkValidity())
        await BootstrapButtons.fire(
          "Error",
          "Password cannot be empty.",
          "error"
        );
      else {
        const resp = await post("/login", {
          username: this.username,
          password: this.password,
          rememberme: this.rememberme,
        });
        if (!resp.ok)
          await BootstrapButtons.fire("Error", await resp.text(), "error");
        else {
          if (this.username != "root")
            localStorage.setItem("username", this.username);
          this.$store.commit("username", this.username);
          this.$router.push("/");
        }
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
