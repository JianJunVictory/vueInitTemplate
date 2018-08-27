<template>
<div id="slide">
  <el-carousel indicator-position="outside"  arrow="always">
    <el-carousel-item v-for="item in imgs" :key="item">
      <div><img :src="item" height="600px"/></div>
    </el-carousel-item>
  </el-carousel>
</div>
</template>
<script>
export default {
  name: 'Slide',
  data () {
    return {
      imgs: ['http://cdn.shopify.com/s/files/1/2500/8302/files/saturdays_banner_1440x.jpg?v=1535315434', 'http://cdn.shopify.com/s/files/1/2500/8302/files/matin_banner_1440x.jpg?v=1535326842', 'http://cdn.shopify.com/s/files/1/2500/8302/files/stone_island_carasel_1440x.jpg?v=1534745085']
    }
  }
}
</script>
<style>
#slide {
    width:74%;
    margin:0 auto;
}
.el-carousel__item div {
    font-size: 18px;
    line-height: 300px;
    margin: 0;
  }
  .el-carousel__item:nth-child(2n) {
    background-color: #99a9bf;
  }
  .el-carousel__item:nth-child(2n+1) {
    background-color: #d3dce6;
  }
.el-carousel__container {
    position: relative;
    height: 600px;
}
</style>
