import Axios, { AxiosError, AxiosResponse } from "axios";
import { STORE } from "./index";

export const Server = Axios.create({
  baseURL: window.location.origin,
  timeout: 5000
});

export const API = Axios.create({
  baseURL: window.location.origin + "/api/",
  timeout: 5000
});

const errorHandler = async (error: AxiosError): Promise<AxiosResponse> => {
  if (undefined === error.response) {
    throw error;
  }
  if (
    error.response.status === 401 ||
    (error.response.status === 403 && !STORE.getters["auth/isAuthenticated"])
  ) {
    const authenticated = await STORE.dispatch("auth/requireLogin", true);
    if (!authenticated) {
      return Promise.reject(error);
    }
    return Server.request(error.config);
  }
  throw error;
};

API.interceptors.response.use(
  response => response,
  error => errorHandler(error)
);
