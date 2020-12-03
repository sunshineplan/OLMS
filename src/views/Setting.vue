<template>
  <header style="padding-left: 20px">
    <a class="h3 title">{{ $t("Setting") }}</a>
    <hr />
  </header>
  <div class="subscribe">
    <h4>{{ $t("Subscribe") }}</h4>
    <div class="form-group">
      <div class="form-group">
        <label for="email">{{ $t("Email") }}</label>
        <input
          type="email"
          class="form-control"
          v-model="email"
          id="email"
          placeholder="{{ $t('email') }})"
          :change="(subscribe = false)"
        />
      </div>
      <div class="form-group form-check">
        <input
          type="checkbox"
          class="form-check-input"
          v-model="subscribe"
          id="subscribe"
          :change="doSubscribe()"
        />
        <label class="form-check-label" for="subscribe">
          {{ $t("SubscribeNotification") }}
        </label>
      </div>
    </div>
  </div>
  <div class="lang">
    <h4>{{ $t("LanguageSetting") }}</h4>
    <div class="form-group">
      <label for="lang">{{ $t("Languages") }}</label>
      <select
        class="form-control"
        v-model.trim="lang"
        id="lang"
        :change="changeLanguage()"
      >
        <option value="en">English</option>
        <option value="zh">简体中文</option>
      </select>
    </div>
  </div>
  <div class="form" @keyup.enter="setting()">
    <h4>{{ $t("ChangePassword") }}</h4>
    <div class="form-group">
      <label for="password">{{ $t("CurrentPassword") }}</label>
      <input
        class="form-control"
        type="password"
        v-model.trim="password"
        id="password"
        maxlength="20"
        required
      />
      <div class="invalid-feedback">{{ $t("RequiredField") }}</div>
    </div>
    <div class="form-group">
      <label for="password1">{{ $t("NewPassword") }}</label>
      <input
        class="form-control"
        type="password"
        v-model.trim="password1"
        id="password1"
        maxlength="20"
        required
      />
      <div class="invalid-feedback">{{ $t("RequiredField") }}</div>
    </div>
    <div class="form-group">
      <label for="password2">{{ $t("ConfirmPassword") }}</label>
      <input
        class="form-control"
        type="password"
        v-model.trim="password2"
        id="password2"
        maxlength="20"
        required
      />
      <div class="invalid-feedback">{{ $t("RequiredField") }}</div>
      <small class="form-text text-muted">{{ $t("MaxPasswordLength") }}</small>
    </div>
    <button class="btn btn-primary" :click="setting()">
      {{ $t("SavePassword") }}
    </button>
    <button class="btn btn-primary" :click="goback()">{{ $t("Back") }}</button>
  </div>
</template>

<script>
import { BootstrapButtons, post, valid, validateEmail } from "../misc.js";

export default {
  name: "Setting",
  data() {
    return {
      email: "",
      subscribe: "",
      lang: document.documentElement.lang,
      password: "",
      password1: "",
      password2: "",
      validated: false,
    };
  },
  mounted() {
    document.title = "Setting";
  },
  methods: {
    changeLanguage() {
      document.cookie = `lang=${this.lang}; Path=/; max-age=31536000`;
      BootstrapButtons.fire(
        this.$t("Success"),
        this.$t("LanguageChanged"),
        "success"
      );
    },
    async doSubscribe() {
      let data;
      if (this.subscribe) {
        if (validateEmail(this.email))
          data = { subscribe: 1, email: this.email };
        else {
          await BootstrapButtons.fire(
            this.$t("Error"),
            this.$t("EmailNotValid"),
            "error"
          );
          this.subscribe = false;
        }
      } else data = { subscribe: 0 };
      const resp = await post("/subscribe", data);
      if ((await resp.json().status) == 1)
        await BootstrapButtons.fire(
          this.$t("Success"),
          this.$t("SubscribeChanged"),
          "success"
        );
      else {
        await BootstrapButtons.fire(
          this.$t("Error"),
          this.$t("EmailNotValid"),
          "error"
        );
        this.subscribe = false;
      }
    },
    async setting() {
      if (valid()) {
        this.validated = false;
        const resp = await post("/setting", {
          password: this.password,
          password1: this.password1,
          password2: this.password2,
        });
        if (!resp.ok)
          await BootstrapButtons.fire("Error", await resp.text(), "error");
        else {
          const json = await resp.json();
          if (json.status == 1) {
            await BootstrapButtons.fire(
              "Success",
              "Your password has changed. Please Re-login!",
              "success"
            );
            this.$store.commit("username", undefined);
            this.$router.push("/");
          } else {
            await BootstrapButtons.fire("Error", json.message, "error");
            if (json.error == 1) this.password = "";
            else {
              this.password1 = "";
              this.password2 = "";
            }
          }
        }
      } else this.validated = true;
    },
  },
};
</script>
