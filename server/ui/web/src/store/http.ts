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
