<template>
  <div>
    <el-container>
      <el-header><Header></Header></el-header>
      <el-main>
      <Slide></Slide>
      <SkuList></SkuList>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import Slide from './slide'
import Header from './header'
import SkuList from './skulist'
export default {
  name: 'HelloWorld',
  components: {Slide, Header, SkuList},
  data () {
    return {
      msg: 'Welcome to Your Vue.js App'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.el-menu-item {
  padding: 0 5px !important;
}
.el-submenu__title {
  padding:none !important;
}
.el-header {
  padding: 0 0px;
}
.el-main {
  margin-top: 11px;
}
</style>
