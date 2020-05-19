import Vuex from "vuex";
import Vue from "vue";
import { AUTH } from "./auth";
import { SEARCH } from "./search";

Vue.use(Vuex);

export const STORE = new Vuex.Store({
  modules: {
    auth: AUTH,
    search: SEARCH
  }
});
