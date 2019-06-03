<template>
  <div class="pg">
    <div>
      <Button type="primary" @click="act">Act</Button>
    </div>
    <br>
    <input v-model="msg">
    <span v-for="(item, i) in listB" :key="'k-' + i">{{i}}-{{item}}</span>
    <div v-if="msg" :title="msg">Looking {{msgpp}} {{testA}} {{testB}}</div>
  </div>
</template>

<script>
export default {
  name: "Playground",
  props: {
    message: String,
    obj: Object
  },
  data() {
    // console.log(axios)
    return {
      msg: "页面加载于 " + new Date().toLocaleString(),
      hb: "Hello, boy!",
      list: [
        {
          id: 0,
          name: "A"
        },
        {
          id: 1,
          name: "B"
        },
        {
          id: 2,
          name: "C"
        }
      ],
      listB: ["A", "B", "C"],
      testA: "A",
      testB: "B"
    };
  },
  computed: {
    msgpp() {
      return this.msg + "++";
    },
    test: {
      get() {
        return { a: this.testA, b: this.testB };
      },
      set(o) {
        this.testA = o.a;
        this.testB = o.b;
      }
    },
  },
  watch: {
    msg() {
      this.msgAlert()
    }
  },
  created() {
    this.act = _.debounce(this.reverseMessage, 100);
    this.msgAlert = _.debounce(() => alert("oh, my god!"), 1000)
  },
  methods: {
    reverseMessage() {
      this.test = { a: "---", b: "___" };
      this.msg = this.msg
        .split("")
        .reverse()
        .join("");
    }
  }
};
</script>

<style lang="less">
</style>
