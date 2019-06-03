import Vue from 'vue'
import _ from "lodash";

// eslint-disable-next-line no-unused-vars
Plugin.install = function (Vue, options) {
  Vue._ = _;
  window._ = _;
  Object.defineProperties(Vue.prototype, {
    _: {
      get() {
        return _;
      }
    },
    $_: {
      get() {
        return _;
      }
    },
  });
};

Vue.use(Plugin)
