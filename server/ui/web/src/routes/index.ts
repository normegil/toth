import VueRouter from "vue-router";
import Settings from "../pages/Settings.vue";

export const ROUTER = new VueRouter({
  routes: [
    {
      path: "/settings",
      component: Settings
    },
  ]
});
