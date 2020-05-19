import { Module } from "vuex";
import { RootState } from "./model";
import SearchResult, { PrintableResult } from "../model/SearchResult";
import Namable from "../model/Namable";
import { VueRouter } from "vue-router/types/router";
import { API } from "./http";
import { AxiosError, AxiosResponse } from "axios";

interface SearchState {
  searched: string;
  searchResults: SearchResult<Namable>[];
}

export const SEARCH: Module<SearchState, RootState> = {
  namespaced: true,
  state: {
    searched: "",
    searchResults: []
  },
  getters: {},
  mutations: {
    setSearched: (state, searched: string): void => {
      state.searched = searched;
    },
    setSearchResults: (
      state,
      searchResults: SearchResult<PrintableResult>[]
    ): void => {
      state.searchResults = searchResults;
    }
  },
  actions: {
    search: (ctx, router: VueRouter): void => {
      console.log("Searching: " + ctx.state.searched);
      if (ctx.state.searched === "") {
        console.error("Search query is empty");
        return;
      }
      API.put("/searches", {
        search: ctx.state.searched
      })
        .then((r: AxiosResponse<SearchResponse>) => {
          ctx.commit("setSearchResults", r.data.results);
          const pathToReach = "/search";
          if (router.currentRoute.path !== pathToReach) {
            return router.push(pathToReach);
          }
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    }
  }
};

interface SearchResponse {
  results: SearchResult<PrintableResult>[];
}
