import axios from "axios";
import Qs from "qs"
import {
  Notice
} from "iview"
import { getToken, freshJWT } from '@/utils/userinfo'

// Full config:  https://github.com/axios/axios#request-config
// axios.defaults.baseURL = process.env.baseURL || process.env.apiUrl || '';
// axios.defaults.headers.common['Authorization'] = AUTH_TOKEN;
// axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';

let config = {
  transformRequest: [function (data) {
    return Qs.stringify(data)
  }],
  // baseURL: process.env.baseURL || process.env.apiUrl || ""
  // timeout: 60 * 1000, // Timeout
  // withCredentials: true, // Check cross-site Access-Control
};

export const _axios = axios.create(config);
_axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
_axios.defaults.headers.put['Content-Type'] = 'application/x-www-form-urlencoded'

_axios.interceptors.request.use(
  function (config) {
    // Do something before request is sent
    const token = getToken()
    if (token) {
      config.headers["token"] = token
    }
    return config;
  },
  function (error) {
    // Do something with request error
    Notice.error({
      title: '请求异常！',
      desc: error,
    });
    return Promise.reject(error);
  }
);

// Add a response interceptor
_axios.interceptors.response.use(
  function (response) {
    // Do something with response data
    const token = response.headers["set-token"]
    if (token) {
      freshJWT(token)
    }
    return response;
  },
  function (error) {
    // Do something with response error
    Notice.error({
      title: '服务器异常！',
      desc: error,
    });
    return Promise.reject(error);
  }
);