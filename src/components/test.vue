<template>
<div>
    {{$t('test')}}
    <div>A:{{message}}</div>
    <el-button type="warning" @click='Achange'>change A</el-button>
    <div>B:{{info}}</div>
    <el-button type="warning" @click='Bchange'>change B</el-button>
</div>
</template>
<script>
import { mapState, mapActions } from 'vuex'
export default {
  name: 'Test',
  data () {
    return {
      title: 'Test'
    }
  },
  methods: {
    ...mapActions(['Bmodify', 'Amodify']),
    Bchange () {
      this.$store.dispatch('Bmodify')
    },
    Achange () {
      this.$store.dispatch('Amodify')
    }
  },
  computed: {
    ...mapState({
      message: state => state.a.message,
      info: state => state.b.info
    })
  }
}
</script>
