<template>
    <div class="loading">
        Loading <span v-for="i in point" :key="i">.</span>
    </div>
</template>

<script>
export default {
    data() {
        return {
            point: 1,
        }
    },
    created() {
        this.nextTip()
    },
    methods: {
        nextTip() {
            this.point %= 6
            this.point++
            setTimeout(this.nextTip, 300)
        }
    }
}
</script>

<style lang="less">
.loading {
    margin-top: 40%;
    text-align: center;
    font-size: 200%;
}
</style>
