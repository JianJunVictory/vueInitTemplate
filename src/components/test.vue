<template>
<div>
  测试页面
    {{$t('test')}}
    <div>A:{{message}}</div>
    <el-button type="warning" @click='Achange'>change A</el-button>
    <div>B:{{info}}</div>
    <el-button type="warning" @click='Bchange'>change B</el-button>
    <el-button type="warning" @click='logout'>LOGOUT</el-button>
    <div v-for="(item,index) in testData" :key="index">{{item}}</div>
</div>
</template>
<script>
import { mapState, mapActions } from 'vuex'
export default {
  name: 'Test',
  data () {
    return {
      title: 'Test',
      testData: []
    }
  },
  methods: {
    ...mapActions(['Bmodify', 'Amodify', 'doLogout', 'doTest']),
    Bchange () {
      this.$store.dispatch('Bmodify')
    },
    Achange () {
      this.$store.dispatch('Amodify')
    },
    logout () {
      let that = this
      this.$store.dispatch('doLogout').then(() => {
        that.$router.replace('/login')
      })
    }
  },
  computed: {
    ...mapState({
      message: state => state.a.message,
      info: state => state.b.info
    })
  },
  mounted () {
    let that = this
    this.$store.dispatch('doTest').then(respData => {
      that.testData = respData
    }).catch(err => {
      console.log(err)
    })
  }
}
</script>
