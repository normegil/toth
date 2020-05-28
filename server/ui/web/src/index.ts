import Vue from "vue";
import VueI18n from "vue-i18n";
import VueRouter from "vue-router";
import App from "./App.vue";
import { ROUTER } from "./routes";
import { STORE } from "./store";
import { TRANSLATOR_OPTIONS } from "./assets/translations";
import "./assets/scss/index.scss";
import "line-awesome/dist/line-awesome/css/line-awesome.min.css";

Vue.use(VueI18n);
Vue.use(VueRouter);

const i18n = new VueI18n(TRANSLATOR_OPTIONS);

STORE.dispatch("auth/loadAuthentication")
    .then(() => {
      new Vue({
        i18n: i18n,
        router: ROUTER,
        store: STORE,
        // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
        render: h => h(App)
      }).$mount("#app");
    })
    .catch((err: Error) => {
      console.error(err)
    })
