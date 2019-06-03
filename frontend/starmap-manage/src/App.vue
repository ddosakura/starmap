<template>
  <div id="app">
    <template v-if="authState == 0">
      <!-- Loading -->
      <LoadingPage/>
    </template>
    <template v-else-if="authState == 1">
      <!-- 登录(后台无法注册) -->
      <LoginPage />
    </template>
    <template v-else>
      <!-- 后台页面 -->
      <div id="nav">
        <router-link to="/">Home</router-link>|
        <router-link to="/about">About</router-link>
      </div>
      <router-view/>
    </template>
  </div>
</template>

<script>
import LoadingPage from "@/views/LoadingPage.vue";
import LoginPage from "@/views/LoginPage.vue";

export default {
  components: {
    LoadingPage,
    LoginPage
  },
  data() {
    //this.$axios
    //  .get("/auth/user/login?username=regtest&password=123")
    //  .then(res => console.log(res.data))
    //  .then(() => {});
    return {
      // authState for app
      // 0: loading
      // 1: logout
      // 2: login
      authState: 0
    };
  },
  created() {
    this.$axios
      .get("/auth/user/info")
      .then(res => {
        this.authState = 2
      }).catch((e) => {
        console.log(e)
        this.authState = 1
      });
  }
};
</script>

<style lang="less">
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#nav {
  padding: 30px;
  a {
    font-weight: bold;
    color: #2c3e50;
    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
