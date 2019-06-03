module.exports = {
  publicPath: '/manage',
  devServer: {
    // host: 'localhost',
    // port: 3000,
    proxy: {
      '/': {
        target: 'http://127.0.0.1:8080/',
        ws: true,
        changeOrigin: true,
        //pathRewrite: { //重写路径
        //  "^/admin": ''
        //}
      },
    },
  },
}