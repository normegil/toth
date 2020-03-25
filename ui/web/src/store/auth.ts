import { Module } from "vuex";
import { RootState } from "./model";
import { API, Server } from "./http";
import { AxiosResponse } from "axios";
import sleep from "../tools/sleep";
import User from "../model/User";
import ResourceRights from "../model/ResourceRights";
import Role from "../model/Role";

interface AuthState {
  showLoginModal: boolean;
  authentifiedUser: User | undefined;
  rights: ResourceRights[] | undefined;
}

interface LoginInformations {
  username: string;
  password: string;
}

export const AUTH: Module<AuthState, RootState> = {
  namespaced: true,
  state: {
    showLoginModal: false,
    authentifiedUser: new User(
      "0",
      "Marie-Odile Barvaux",
      new Role("0", "teacher")
    ),
    rights: undefined
  },
  getters: {
    isAuthenticated: (state): boolean => {
      return !(
        undefined === state.authentifiedUser ||
        state.authentifiedUser.name === "anonymous"
      );
    },
    hasAccess: state => {
      return (resource: string, action: string): boolean => {
        if (state.rights === undefined) {
          return false;
        }

        for (const right of state.rights) {
          if (right.name === resource) {
            for (const rightAction of right.allowedActions) {
              if (rightAction === action) {
                return true;
              }
            }
          }
        }
        return false;
      };
    }
  },
  mutations: {
    setShowLoginModal: (state, show: boolean): void => {
      state.showLoginModal = show;
    },
    setRights: (state, rights: ResourceRights[]): void => {
      state.rights = rights;
    },
    setAuthentified: (state, authentifiedUser: User | undefined): void => {
      let username: string | undefined;
      if (authentifiedUser === undefined) {
        username = undefined;
      } else {
        username = authentifiedUser.name;
      }
      console.log("Change user to: " + username);
      state.authentifiedUser = authentifiedUser;
    }
  },
  actions: {
    loadAuthentication: async (ctx): Promise<void> => {
      const response: AxiosResponse<User> = await API.get("/users/current");
      ctx.commit("setAuthentified", response.data);
      return ctx.dispatch("loadRights");
    },
    loadRights: async (ctx): Promise<void> => {
      const response: AxiosResponse<ResourceRights> = await API.get(
        "/rights/current"
      );
      ctx.commit("setRights", response.data);
    },
    signIn: async (ctx, login: LoginInformations): Promise<void> => {
      const response: AxiosResponse<User> = await Server.get("/auth/sign-in", {
        auth: login
      });
      console.log("Sign in as: " + response.data.name);
      ctx.commit("setAuthentified", response.data);
      ctx.commit("setShowLoginModal", false);
      return ctx.dispatch("loadRights");
    },
    signOut: async (ctx): Promise<void> => {
      await Server.get("/auth/sign-out", {
        headers: {
          "X-Authentication-Action": "sign-out"
        }
      });
      ctx.commit("setAuthentified", undefined);
      return ctx.dispatch("loadRights");
    },
    requireLogin: async (ctx): Promise<boolean> => {
      ctx.commit("setShowLoginModal", true);
      while (!ctx.getters.isAuthenticated && ctx.state.showLoginModal) {
        await sleep(200);
      }
      return new Promise((resolve): void => {
        resolve(ctx.getters.isAuthenticated);
      });
    }
  }
};
