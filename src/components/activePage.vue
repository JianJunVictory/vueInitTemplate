<template>
    <h3></h3>
</template>
<script>
import {mapActions} from 'vuex'
export default {
  name: 'ActivePage',
  data () {
    return {}
  },
  methods: {
    ...mapActions(['doActive']),
    active () {
      let token = this.$route.query.token
      let that = this
      this.$store.dispatch('doActive', token).then(() => {
        that.$router.replace('/')
      }).catch(err => {
        console.log(err)
      })
    }
  },
  mounted () {
    this.active()
  }
}
</script>
