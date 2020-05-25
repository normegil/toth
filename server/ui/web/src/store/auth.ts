import { Module } from "vuex";
import { RootState } from "./model";
import { API, Server } from "./http";
import { AxiosResponse } from "axios";
import sleep from "../tools/sleep";
import User from "../model/User";
import ResourceRights from "../model/ResourceRights";
import Role from "../model/Role";

interface AuthState {
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
    authentifiedUser: undefined,
    rights: undefined
  },
  getters: {
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
    },
    loadRights: async (ctx): Promise<void> => {
      const response: AxiosResponse<ResourceRights> = await API.get(
        "/rights/current"
      );
      ctx.commit("setRights", response.data);
    },
    signIn: async (ctx, login: LoginInformations): Promise<void> => {
      const response: AxiosResponse<User> = await API.get("/auth/log-in", {
        auth: login
      });
      return ctx.dispatch("loadAuthentication");
    },
    signOut: async (ctx): Promise<void> => {
      await API.get("/auth/log-out", {
        headers: {
          "X-Authentication-Action": "sign-out"
        }
      });
      ctx.commit("setAuthentified", undefined);
    }
  }
};
