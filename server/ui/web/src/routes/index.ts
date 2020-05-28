import VueRouter from "vue-router";
import Settings from "../pages/Settings.vue";
import Classes from "../pages/Classes.vue";

export const ROUTER = new VueRouter({
  routes: [
    {
      path: "/classes",
      component: Classes
    },{
      path: "/settings",
      component: Settings
    },
  ]
});
