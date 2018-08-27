<template>
<div>
   <div v-for="item in SkuList" :key="item.id" style="float:left">
     <img :src="item.iamgeUrl" width="50%">
   </div>
</div>
</template>
<script>
import { mapActions } from 'vuex'
export default {
  name: 'SkuList',
  data () {
    return {
      lists: [1, 2, 3, 4],
      SkuList: []
    }
  },
  methods: {
    ...mapActions(['doGetSkuList']),
    getData () {
      let that = this
      this.$store.dispatch('doGetSkuList').then(data => {
        console.log(data)
        that.SkuList = data
      }).catch(err => {
        console.log(err)
      })
    }
  },
  mounted () {
    this.getData()
  }
}
</script>
