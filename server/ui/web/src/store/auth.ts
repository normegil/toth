import { Module } from "vuex";
import { RootState } from "./model";
import { API, Server } from "./http";
import { AxiosResponse } from "axios";
import User from "../model/User";
import ResourceRights from "../model/ResourceRights";

interface AuthState {
  authentifiedUser: User | undefined;
  rights: ResourceRights[] | undefined;
}

interface LoginInformations {
  username: string;
  password: string;
}

interface PasswordChange {
  id: string;
  currentPassword: string;
  newPassword: string;
}

interface ProfileChange {
  id: string;
  name: string;
  mail: string;
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
    updateProfile: async (ctx, info: ProfileChange): Promise<void> => {
      let response = await API.put("/users", {
        "id": info.id,
        "name":info.name,
        "mail":info.mail,
      }, {
        headers: {
          "Content-Type": "application/json"
        }
      });
      return ctx.dispatch("loadAuthentication");
    },
    updatePassword: async (ctx, info: PasswordChange): Promise<void> => {
      let response = await API.put("/users/" + info.id + "/password", {
        "current_password":info.currentPassword,
        "new_password":info.newPassword,
      }, {
        headers: {
          "Content-Type": "application/json"
        }
      });
      return ctx.dispatch("loadAuthentication");
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
