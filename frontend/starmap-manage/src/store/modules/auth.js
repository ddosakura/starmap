import * as types from '@/store/mutation-types'
import * as authAPI from '@/api/auth'
import {
  md5
} from '@/utils/func'
import {
  getUserInfo,
  cleanJWT
} from '@/utils/userinfo'
import {
  Notice,
  Message,
} from "iview"

export const moduleAuth = {
  state: {
    // auth-state for app
    // 0: loading
    // 1: logout
    // 2: login
    state: 0,
  },
  mutations: {
    [types.UPDATE_AUTH_STATE](state, payload) {
      state.state = payload
    },
  },
  actions: {
    async getAuthInfo({
      commit,
      // eslint-disable-next-line no-unused-vars
      state
    }) {
      const [d, keep] = getUserInfo()
      if (keep) {
        Message.info("欢迎回来，" + d.nickname + "！")
        commit(types.UPDATE_AUTH_STATE, 2)
        return
      }
      // TODO:
      // 1. 简化api写法
      // 2. 改为获取有权限的页面
      const res = await authAPI.getInfo()
      if (res.code === 0) {
        commit(types.UPDATE_AUTH_STATE, 2)
      } else if (res.msg == "jwt invalid") {
        commit(types.UPDATE_AUTH_STATE, 1)
      } else {
        Notice.error({
          title: '请求失败！',
          desc: res.msg,
        });
      }
    },
    async login({
      commit
    }, {
      user,
      pass
    }) {
      //if (!user || !pass) {
      //  alert("请填写完整")
      //  return
      //}
      const res = await authAPI.login(user, md5(pass))
      if (res.code === 0) {
        commit(types.UPDATE_AUTH_STATE, 2)
      } else {
        alert(res.msg);
      }
    },
    logout() {
      cleanJWT()
    }
  },
  getters: {
    // eslint-disable-next-line no-unused-vars
    authState: (state, getters) => {
      return state.state
    }
  }
}