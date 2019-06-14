<template>
  <div class="login">
    <Card class="card">
      <p slot="title">Star-Map Manage</p>
      <Form ref="form" :model="formData" :rules="formRule">
        <FormItem prop="user">
          <i-input type="text" v-model="formData.user" placeholder="Username">
            <Icon type="ios-person-outline" slot="prepend"></Icon>
          </i-input>
        </FormItem>
        <FormItem prop="pass">
          <i-input type="password" v-model="formData.pass" placeholder="Password">
            <Icon type="ios-lock-outline" slot="prepend"></Icon>
          </i-input>
        </FormItem>
        <FormItem>
          <Button type="primary" @click="handleLogin('form')">Log In</Button>
        </FormItem>
      </Form>
    </Card>
  </div>
</template>

<script>
import { mapState, mapActions, mapGetters } from "vuex";

export default {
  data() {
    return {
      formData: {
        user: "",
        pass: ""
      },
      formRule: {
        user: [
          {
            required: true,
            message: "Please fill in the username",
            trigger: "blur"
          }
        ],
        pass: [
          {
            required: true,
            message: "Please fill in the password.",
            trigger: "blur"
          },
          {
            type: "string",
            min: 6,
            message: "The password length cannot be less than 6 bits",
            trigger: "blur"
          }
        ]
      }
    };
  },
  methods: {
    handleLogin(name) {
      this.$refs[name].validate(valid => {
        console.log(valid);
        if (valid) {
          this.login(this.formData);
        }
      });
    },
    ...mapActions({
      login: "login"
    })
  },
  created() {}
};
</script>

<style lang="less">
.card {
  width: 320px;
  margin: 0 auto;
  margin-top: 35%;
}
</style>
