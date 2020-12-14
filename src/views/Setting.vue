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
          :placeholder="$t('Email')"
          @input="subscribe = false"
        />
      </div>
      <div class="form-group form-check">
        <input
          type="checkbox"
          class="form-check-input"
          v-model="subscribe"
          id="subscribe"
          @change="doSubscribe()"
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
        @change="changeLanguage()"
      >
        <option value="en">English</option>
        <option value="zh">简体中文</option>
      </select>
    </div>
  </div>
  <div class="form" @keyup.enter="changePassword()">
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
        autofocus
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
    <button class="btn btn-primary" @click="changePassword()">
      {{ $t("SavePassword") }}
    </button>
    <button class="btn btn-primary" @click="goback()">{{ $t("Back") }}</button>
  </div>
</template>

<script>
import Cookies from "js-cookie";
import { loadLocaleMessages } from "../i18n";
import { post, valid, validateEmail } from "../misc.js";

const grecaptcha = window.grecaptcha;

export default {
  name: "Setting",
  data() {
    return {
      recaptcha: this.$store.state.recaptcha,
      email: "",
      subscribe: false,
      lang: document.documentElement.lang,
      password: "",
      password1: "",
      password2: "",
      validated: false,
    };
  },
  async created() {
    await this.getSubscribe();
  },
  mounted() {
    document.title = this.$t("Setting") + " - " + this.$t("OLMS");
  },
  methods: {
    async changeLanguage() {
      Cookies.set("lang", this.lang, { expires: 365 });
      this.$i18n.locale = this.lang;
      await loadLocaleMessages(this.$i18n.locale);
      document.querySelector("html").setAttribute("lang", this.$i18n.locale);
      this.prompt("Success", "LanguageChanged", "success");
      this.$router.replace("/setting");
    },
    async getSubscribe() {
      const resp = await fetch("/subscribe");
      await this.checkResp(resp, async () => {
        const json = await resp.json();
        this.subscribe = json.subscribe;
        if (json.subscribe) this.email = json.email;
      });
    },
    async doSubscribe() {
      let data;
      if (this.subscribe) {
        if (validateEmail(this.email))
          data = { subscribe: true, email: this.email };
        else {
          this.subscribe = false;
          await this.prompt("Error", "EmailNotValid", "error");
          return;
        }
      } else data = { subscribe: false };
      const resp = await post("/subscribe", data);
      await this.checkResp(resp, async () => {
        const json = await resp.json();
        if (json.status == 1)
          await this.prompt("Success", "SubscribeChanged", "success");
        else {
          this.subscribe = false;
          await this.prompt("Error", "Error", "error");
        }
      });
    },
    async changePassword() {
      if (valid()) {
        this.validated = false;
        let data = {
          password: this.password,
          password1: this.password1,
          password2: this.password2,
        };
        if (this.recaptcha)
          data.recaptcha = await grecaptcha.execute(this.recaptcha, {
            action: "setting",
          });
        const resp = await post("/setting", data);
        await this.checkResp(resp, async () => {
          const json = await resp.json();
          if (json.status == 1) {
            await this.prompt("Success", "PasswordChanged", "success");
            this.$store.commit("user", null);
            this.$router.push("/login");
          } else {
            await this.prompt("Error", json.message, "error");
            if (json.error == 1) this.password = "";
            else {
              this.password1 = "";
              this.password2 = "";
            }
          }
        });
      } else this.validated = true;
    },
  },
};
</script>
